package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type catchupStruct struct {
	// json tag to serialize json body
	SetCatchup bool `json:setCatchup`
}

func endpointTest(c gin.Context) {
	body := catchupStruct{}
	// using BindJson method to serialize body with struct
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(body)
	c.JSON(http.StatusAccepted, &body)
}

func postCatchup(c *gin.Context) {

	body := catchupStruct{}

	// using BindJson method to serialize body with struct
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// fmt.Println(body)

	catchupRequired = body.SetCatchup
	c.JSON(http.StatusAccepted, gin.H{"message": "Accepted. catchupRequired set to " + strconv.FormatBool(catchupRequired)})

}

func getCatchup(c *gin.Context) {
	c.JSON(200, gin.H{"catchupRequired": catchupRequired})
}
