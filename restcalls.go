package main

import (
	"encoding/json"
	"fmt"
	"io"
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

