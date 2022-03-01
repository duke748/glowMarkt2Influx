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
)

// Gets all the virtual Entity Id's and stores them in Variables
func getVirtualEntities() error {
	checkIfTokenExpired(tokenStruct.Exp)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", baseUrl+"virtualentity/", nil)
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

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		// Create a var of ve's details of type virtualEntities struct
		var ves virtualEntities

		// Unmarshall into the ves object
		json.Unmarshal([]byte(bodyString), &ves)

		// Loop of the resources slice and return the Ids
		for k := range ves[0].Resources {
			// fmt.Println("name: ", ves[0].Resources[k].Name)
			switch ves[0].Resources[k].Name {
			case "gas consumption":
				gasConsumptionId = ves[0].Resources[k].ResourceID
			case "gas cost":
				gasCostId = ves[0].Resources[k].ResourceID
			case "electricity consumption":
				electricityConsumptionId = ves[0].Resources[k].ResourceID
			case "electricity cost":
				electricityCostId = ves[0].Resources[k].ResourceID

			}
		}

		fmt.Println("gasConsumptionId: ", gasConsumptionId)
		fmt.Println("gasCostId: ", gasCostId)
		fmt.Println("electricityConsumptionId: ", electricityConsumptionId)
		fmt.Println("electricityCostId: ", electricityCostId)
	}

	defer response.Body.Close()

	return nil
}

func sendCatchup(entityId string) {
	checkIfTokenExpired(tokenStruct.Exp)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", baseUrl+"resource/"+entityId+"/catchup", nil)
	if err != nil {
		fmt.Println("Error")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("applicationId", applicationId)
	req.Header.Add("token", token)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error")
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Sent Refresh")
	}

	defer response.Body.Close()
	time.Sleep(30 * time.Second)
}

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
