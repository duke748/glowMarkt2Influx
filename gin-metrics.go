package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type returnMetricsStruct struct {
	NumberOfDays int `json:"numberOfDays"`
}

func returnMetricsData(c *gin.Context) {
	body := returnMetricsStruct{}

	// using BindJson method to serialize body with struct
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// fmt.Println(body)

	days := body.NumberOfDays
	c.JSON(http.StatusAccepted, gin.H{"message": "Accepted. Retrieving data"})

	for i := -1; i > -days; i-- {
		fmt.Println(i)

		// date range
		dateStart := time.Now().AddDate(0, 0, i).Format("2006-01-02") + "T00:00:00"
		dateEnd := time.Now().AddDate(0, 0, i).Format("2006-01-02") + "T23:59:59"

		getMeterReadings(electricityConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Electricity", false)
		getMeterReadings(gasConsumptionId, "PT30M&function=sum&from="+dateStart+"&to="+dateEnd, "Gas", false)

	}

}
