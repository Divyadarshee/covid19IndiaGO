package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmespath/go-jmespath"
	"github.com/labstack/echo/v4"
	"go-swag-sample/echosimple/configs"
	state_data "go-swag-sample/echosimple/data"
	"go-swag-sample/echosimple/models"
	"go-swag-sample/echosimple/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"time"
)

var caseCollection *mongo.Collection = configs.GetCollections(configs.DB, "cases")
var validate = validator.New()

// GetCases godoc
// @Summary Get latest covid19 cases.
// @Description Get the covid19 cases details for the given GPS coordinates.
// @Tags user
// @Accept */*
// @Produce json
// @Param latitude query float64 true "Latitude"
// @Param longitude query float64 true "Longitude"
// @Success 200 {object} map[string]interface{}
// @Router /cases [get]
func GetCases(c echo.Context) error {

	var stateCases models.Cases
	var indiaCases models.Cases
	var stateCode echo.Map

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// creates query params binder that stops binding at first error
	bindingErr := echo.QueryParamsBinder(c).
		MustFloat64("latitude", &models.GPS.Latitude).
		MustFloat64("longitude", &models.GPS.Longitude).
		BindError() // returns first binding error

	defer cancel()

	if bindingErr != nil {
		return c.JSON(http.StatusBadRequest, responses.CaseResponse{Status: http.StatusBadRequest, Message: bindingErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
	}

	if validationErr := validate.Struct(&models.GPS); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.CaseResponse{Status: http.StatusBadRequest, Message: validationErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
	}

	// Get the state code along with the name of the state
	stateCode = GetStateFromGPS(models.GPS.Latitude, models.GPS.Longitude)

	_, isError := stateCode["Error"]
	if isError == true {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: stateCode["Error"].(string), CasesInState: stateCases, CasesInIndia: indiaCases})
	}

	stateDBQueryErr := caseCollection.FindOne(ctx, bson.M{"statecode": stateCode["StateCode"]}).Decode(&stateCases)
	indiaDBQueryErr := caseCollection.FindOne(ctx, bson.M{"statecode": "TT"}).Decode(&indiaCases)

	if stateDBQueryErr != nil {
		return c.JSON(http.StatusNotFound, responses.CaseResponse{Status: http.StatusNotFound, Message: stateDBQueryErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
	}

	if indiaDBQueryErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: indiaDBQueryErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
	}

	return c.JSON(http.StatusOK, responses.CaseResponse{Status: http.StatusOK, Message: "success", CasesInState: stateCases, CasesInIndia: indiaCases})
}

func GetStateFromGPS(latitude float64, longitude float64) echo.Map {
	reverseGeocodingUrl := fmt.Sprintf(state_data.ReverseGeocodingUrl, latitude, longitude)
	method := "GET"

	client := &http.Client{}
	reverseGeocodingRequest, reverseGeocodingRequestErr := http.NewRequest(method, reverseGeocodingUrl, nil)

	if reverseGeocodingRequestErr != nil {
		fmt.Println(reverseGeocodingRequestErr)
		return echo.Map{"Error": reverseGeocodingRequestErr.Error()}
	}
	reverseGeocodingResponse, reverseGeocodingResponseErr := client.Do(reverseGeocodingRequest)
	if reverseGeocodingResponseErr != nil {
		fmt.Println(reverseGeocodingResponseErr)
		return echo.Map{"Error": reverseGeocodingResponseErr.Error()}
	}
	defer reverseGeocodingResponse.Body.Close()

	// reading the body of the response
	body, ReadErr := ioutil.ReadAll(reverseGeocodingResponse.Body)
	if ReadErr != nil {
		fmt.Println(ReadErr)
		return echo.Map{"Error": ReadErr.Error()}
	}
	fmt.Println(string(body))
	//fmt.Printf("type of body is %T", body)
	//fmt.Println()
	var data map[string]interface{}

	parsingErr := json.Unmarshal(body, &data)
	if parsingErr != nil {
		fmt.Println(parsingErr)
		return echo.Map{"Error": parsingErr.Error()}
	}

	// jquery using jmespath to get state name from the given gps location
	stateName, stateNameErr := jmespath.Search("address.state", data)
	if stateNameErr != nil {
		fmt.Println(stateNameErr)
		return echo.Map{"Error": stateNameErr.Error()}
	}

	// Converting state name to state code
	for key, value := range state_data.StateCodes {
		if stateName == key {
			fmt.Printf("%q is the key for the value %q\n", key, value)
			return echo.Map{"State": key, "StateCode": value}
		}
	}

	return echo.Map{"Error": "GPS location not in India"}
}
