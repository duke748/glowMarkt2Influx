package main

import "time"

type authBody struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ApplicationId string `json:"applicationId"`
}

type authorisation struct {
	Valid bool   `json:"valid"`
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
	// UserGroups              []interface{} `json:"userGroups"`
	// FunctionalGroupAccounts []interface{} `json:"functionalGroupAccounts"`
	// AccountID               string        `json:"accountId"`
	// IsTempAuth              bool          `json:"isTempAuth"`
	Name string `json:"name"`
}

type vEConsumptionDataSlice struct {
	Timestamp time.Time
	Kwh       float64
}

type veConsumptionData struct {
	Classifier string `json:"classifier"`
	// Data       [][]float64 `json:"data"`
	// Data []struct {
	// 	Timestamp time.Time
	// 	Kwh       float64
	// } `json:"data"`
	Data  []vEConsumptionDataSlice `json:"data"`
	Name  string                   `json:"name"`
	Query struct {
		From     string `json:"from"`
		Function string `json:"function"`
		Period   string `json:"period"`
		To       string `json:"to"`
	} `json:"query"`
	ResourceID     string `json:"resourceId"`
	ResourceTypeID string `json:"resourceTypeId"`
	Status         string `json:"status"`
	Units          string `json:"units"`
}

type virtualEntities []struct {
	// Clone         bool   `json:"clone"`
	// Active        bool   `json:"active"`
	// ApplicationID string `json:"applicationId"`
	// VeTypeID      string `json:"veTypeId"`
	// PostalCode    string `json:"postalCode"`
	Resources []struct {
		ResourceID     string `json:"resourceId"`
		ResourceTypeID string `json:"resourceTypeId"`
		Name           string `json:"name"`
	} `json:"resources"`
	// OwnerID    string        `json:"ownerId"`
	// Name       string        `json:"name"`
	// VeChildren []interface{} `json:"veChildren"`
	// VeID       string        `json:"veId"`
	// UpdatedAt  time.Time     `json:"updatedAt"`
	// CreatedAt  time.Time     `json:"createdAt"`
}
