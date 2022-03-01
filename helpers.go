package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// func getLast30days() {

// 	for i := -1; i > -31; i-- {
// 		fmt.Println(i)

// 		// Yesterdays date range
// 		yesterdayDateStart := time.Now().AddDate(0, 0, i).Format("2006-01-02") + "T00:00:00"
// 		yesterdayDateEnd := time.Now().AddDate(0, 0, i).Format("2006-01-02") + "T23:59:59"

// 		fmt.Println("yesterdayDateStart: ", yesterdayDateStart)
// 		fmt.Println("yesterdayDateEnd: ", yesterdayDateEnd)

// 		getMeterReadings(electricityConsumptionId, "PT30M&function=sum&from="+yesterdayDateStart+"&to="+yesterdayDateEnd, "Electricity")
// 		getMeterReadings(gasConsumptionId, "PT30M&function=sum&from="+yesterdayDateStart+"&to="+yesterdayDateEnd, "Gas")

// 	}

// }

func getLastXdays(days int) {

	for i := -1; i > -days; i-- {
		fmt.Println(i)

		// date range
		dateStart := time.Now().AddDate(0, 0, i).Format("2006-01-02") + "T00:00:00"
		dateEnd := time.Now().AddDate(0, 0, i).Format("2006-01-02") + "T23:59:59"

		getMeterReadings(electricityConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Electricity", false)
		getMeterReadings(gasConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Gas", false)

	}

}

// Returns time to the nearest 30 minutes of Glowmarkt data
func getTimeToNearest30() string {

	// Get current HH:MM & ignore seconds
	hour, min, _ := time.Now().Clock()

	switch true {
	case min <= 29:
		strMin = "00"
		break
	case min >= 30:
		strMin = "30"

	}
	var strHour string
	if hour < 10 {
		strHour = "0" + strconv.Itoa(hour)
	} else {
		strHour = strconv.Itoa(hour)
	}

	timeTo := fmt.Sprintf("T%s:%s:00", strHour, strMin)

	return timeTo
}

// Will check if token has expired by passing the expiry time
func checkIfTokenExpired(expiryTime int64) bool {
	t := time.Now().Unix()
	if expiryTime > t {
		fmt.Printf("Token has expired '%s'. ", ConvertTime(expiryTime))

		// Get new token
		tokenStruct, err := getToken()
		if err != nil {
			log.Fatal("ERROR: ", err)
		}

		// Set variables with updated token details
		token = tokenStruct.Token
		tokenExpiry = tokenStruct.Exp
		os.Setenv("glowToken", token)
		os.Setenv("glowTokenExpiry", strconv.Itoa(int(tokenExpiry)))
		return false
	}
	return true
}

// Gets Environment variables
func getEnvVar(strIn string) string {
	strRet := os.Getenv(strIn)
	if strRet == "" {
		log.Fatal("Missing environment variable ", strIn)
	}
	fmt.Printf("%s = %s\n", strIn, strRet)
	return strRet
}

// Converts from Unix timestamp to readable time
func ConvertTime(timeIn int64) time.Time {
	unixTime := time.Unix(timeIn, 0)
	return unixTime
}

// Sets up all the required environment variables
func setRequiredEnvironmentVariables() {
	// Pull glow username and password from environment variables
	glowUsername = getEnvVar("glowUsername")
	glowPassword = getEnvVar("glowPassword")
	influxDbToken = getEnvVar("influxDbToken")
	influxDbUrl = getEnvVar("influxDbUrl")
	influxDbOrg = getEnvVar("influxDbOrg")
	influxDbBucket = getEnvVar("influxDbBucket")

	defaultInterval, _ = strconv.Atoi(os.Getenv("defaultInterval"))
	if defaultInterval == 0 {
		defaultInterval = 30
	}
}
