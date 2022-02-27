package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// set global variables
var token string                    // Bearer Token used to authenticate with the api
var tokenExpiry int64               // Unix Time stamp date that the token expires
var tokenStruct authorisation       // Struct from the Auth API
var min int                         // Minutes
var strMin string                   // Stringified version of minutes
var gasConsumptionId string         // Id returned from the Virtual Entities API for gas Consumption
var gasCostId string                // Id returned from the Virtual Entities API for gas cost
var electricityConsumptionId string // Id returned from the Virtual Entities API for electricity Consumption
var electricityCostId string        // Id returned from the Virtual Entities API for electricity Cost
var influxDbToken string            // token to communicate with InfluxDb
var influxDbUrl string              // URL for InfluxDB
var influxDbOrg string              // InfluxDb organisation
var influxDbBucket string           // InfluxDb bucket
var glowUsername string             // Username for the Glow App
var glowPassword string             // Password for the Glow App
var catchupRequired bool            // Do we need to send a catchup to the API

// Constants
const baseUrl = "https://api.glowmarkt.com/api/v0-1/"
const authUrl = baseUrl + "auth"
const resourceUrl = baseUrl + "resource/"
const applicationId = "b0f1b774-a586-4f72-9edd-27ead8aa7a8d"

func pullTodaysReadings() {
	// Grab the current time to nearest 30 minutes of data that will be available
	timeTo := getTimeToNearest30()

	// Grab todays date range
	dateStart := time.Now().Format("2006-01-02") + "T00:00:00"
	dateEnd := time.Now().Format("2006-01-02") + timeTo

	// Yesterdays date range
	yesterdayDateStart := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + "T00:00:00"
	yesterdayDateEnd := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + "T23:59:59"

	fmt.Println("dateStart: ", dateStart)
	fmt.Println("dateEnd: ", dateEnd)
	fmt.Println("yesterdayDateStart: ", yesterdayDateStart)
	fmt.Println("yesterdayDateEnd: ", yesterdayDateEnd)

	getMeterReadings(electricityConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Electricity", catchupRequired)
	getMeterReadings(gasConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Gas", catchupRequired)

}

func main() {
	// Populate all the required
	setRequiredEnvironmentVariables()

	// Grab token if required
	if token == "" {
		tokenStruct, err := getToken()
		if err != nil {
			log.Fatal("ERROR: ", err)
		}
		token = tokenStruct.Token
		tokenExpiry = tokenStruct.Exp

		fmt.Println("token: ", token)
		fmt.Println("tokenExpiry: ", tokenExpiry)
	}

	// populate the virtual entities
	getVirtualEntities()

	// Setup a func to run every 3 minutes
	go func() {
		c := time.Tick(3 * time.Minute)
		for range c {
			// Note this purposfully runs the function
			// in the same goroutine so we make sure there is
			// only ever one. If it might take a long time and
			// it's safe to have several running just add "go" here.
			pullTodaysReadings()
		}
	}()

	// Setup Gin Routes
	r := gin.Default()

	r.POST("/catchup", postCatchup)
	r.GET("/catchup", getCatchup)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// Grab the current time to nearest 30 minutes of data that will be available
	timeTo := getTimeToNearest30()

	// Grab todays date range
	dateStart := time.Now().Format("2006-01-02") + "T00:00:00"
	dateEnd := time.Now().Format("2006-01-02") + timeTo

	// Yesterdays date range
	yesterdayDateStart := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + "T00:00:00"
	yesterdayDateEnd := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + "T23:59:59"

	fmt.Println("dateStart: ", dateStart)
	fmt.Println("dateEnd: ", dateEnd)
	fmt.Println("yesterdayDateStart: ", yesterdayDateStart)
	fmt.Println("yesterdayDateEnd: ", yesterdayDateEnd)

	// Yesterdays data
	// getMeterReadings(electricityConsumptionId, "PT30M&function=sum&from="+yesterdayDateStart+"&to="+yesterdayDateEnd, "Electricity")
	// getMeterReadings(gasConsumptionId, "PT30M&function=sum&from="+yesterdayDateStart+"&to="+yesterdayDateEnd, "Gas")

	// getLast30days()

	catchupRequired = true
	// Todays data
	getMeterReadings(electricityConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Electricity", catchupRequired)
	getMeterReadings(gasConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Gas", catchupRequired)

}

func (r *vEConsumptionDataSlice) UnmarshalJSON(b []byte) error {

	var vals []float64
	if err := json.Unmarshal(b, &vals); err != nil {
		return err
	}

	if len(vals) != 2 {
		return fmt.Errorf("Expected two values in '%s' but got %s", string(b), string(len(vals)))
	}

	r.Timestamp = time.Unix(int64(vals[0]), 0)
	r.Kwh = vals[1]

	return nil
}

func getMeterReadings(resourceId string, period string, endpoint string, catchup bool) error {
	checkIfTokenExpired(tokenStruct.Exp)

	if catchup == true {
		sendCatchup(resourceId)
	}

	fmt.Println("Executing GET to ", resourceUrl+resourceId+"/readings?period="+period)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", resourceUrl+resourceId+"/readings?period="+period, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("applicationId", applicationId)
	req.Header.Add("token", token)

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}

	fmt.Println("Response Code: ", response.StatusCode)
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		// Create a var of veConsumptionDetails details of type veConsumptionData struct
		var veConsumptionDetails veConsumptionData

		// Unmarshall into the veConsumptionDetails object
		json.Unmarshal([]byte(bodyString), &veConsumptionDetails)

		fmt.Println(veConsumptionDetails.Status)

		totalKw := 0.0

		// Create a client
		// You can generate an API Token from the "API Tokens Tab" in the UI
		influxClient := influxdb2.NewClient(influxDbUrl, influxDbToken)

		// always close client at the end
		defer influxClient.Close()

		for k := range veConsumptionDetails.Data {

			// Convert to unixNanotime for influxDB
			unixNanoTime := veConsumptionDetails.Data[k].Timestamp.UnixNano()

			fmt.Printf("Time %s Kw Useage %g \n", veConsumptionDetails.Data[k].Timestamp, veConsumptionDetails.Data[k].Kwh)
			totalKw += veConsumptionDetails.Data[k].Kwh

			// get non-blocking write client
			writeAPI := influxClient.WriteAPI("Greenlands", "glowAPI")

			// write line protocol
			lineProtocol := "Reading,meter=" + endpoint + " kwh=" + fmt.Sprint(veConsumptionDetails.Data[k].Kwh) + " " + fmt.Sprint(unixNanoTime)
			fmt.Println("lineProtocol: ", lineProtocol)

			writeAPI.WriteRecord(lineProtocol)
			// Flush writes
			writeAPI.Flush()

		}
		fmt.Printf("Total KW used was  %g \n", totalKw)
	}

	defer response.Body.Close()

	return nil
}

// Get Token
func getToken() (authorisation, error) {
	// Create authorisationContent obj of type authorisation struct
	var authorisationContent authorisation

	// Construct JSON to send to glowmarkt api auth
	var jsonData = []byte(`{
		"username": "` + glowUsername + `",
		"password": "` + glowPassword + `",
		"applicationId" : "b0f1b774-a586-4f72-9edd-27ead8aa7a8d"
	}`)

	request, error := http.NewRequest("POST", authUrl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if error != nil {
		return authorisationContent, fmt.Errorf("ERROR creating Post to %s -> %s", authUrl, error)
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		return authorisationContent, fmt.Errorf("ERROR creating http client -> %s", error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	respBody, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode == http.StatusOK {
		// Unmarshall into the authorisationContent object
		json.Unmarshal([]byte(respBody), &authorisationContent)
	}

	// Return the struct
	return authorisationContent, nil
}
