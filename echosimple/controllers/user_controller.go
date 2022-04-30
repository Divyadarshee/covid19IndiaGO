package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmespath/go-jmespath"
	"github.com/labstack/echo/v4"
	"go-swag-sample/echosimple/configs"
	state_data "go-swag-sample/echosimple/data"
	"go-swag-sample/echosimple/models"
	"go-swag-sample/echosimple/responses"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollections(configs.DB, "users")
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

	var stateCode echo.Map
	var caseTotalState responses.Cases
	var caseTotalIndia responses.Cases
	var sourceData map[string]interface{}
	var lastUpdatedStateQuery string
	var casesStateQuery string
	var lastUpdatedIndiaQuery string = "TT.meta.last_updated"
	var casesIndiaQuery string = "TT.total"

	// creates query params binder that stops binding at first error
	bindingErr := echo.QueryParamsBinder(c).
		MustFloat64("latitude", &models.GPS.Latitude).
		MustFloat64("longitude", &models.GPS.Longitude).
		BindError() // returns first binding error

	if bindingErr != nil {
		return c.JSON(http.StatusBadRequest, responses.CaseResponse{Status: http.StatusBadRequest, Message: bindingErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	if validationErr := validate.Struct(&models.GPS); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.CaseResponse{Status: http.StatusBadRequest, Message: validationErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// Get the state code along with the name of the state
	stateCode = GetStateFromGPS(models.GPS.Latitude, models.GPS.Longitude)

	// Get call to the source of covid19 data: https://data.covid19india.org/v4/min/data.min.json
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	sourceRequest, sourceRequestErr := http.NewRequest("GET", state_data.Covid19SourceUrl, nil)

	if sourceRequestErr != nil {
		fmt.Println("error")
	}

	sourceResponse, sourceResponseErr := client.Do(sourceRequest)
	if sourceResponseErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: sourceResponseErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// reading the body of the response
	body, readErr := ioutil.ReadAll(sourceResponse.Body)
	if readErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: readErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// parsing json
	parsingErr := json.Unmarshal(body, &sourceData)
	if parsingErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: parsingErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// jquery using jmespath to get state specific total cases, vaccination and last updated details
	lastUpdatedStateQuery = fmt.Sprintf("%s.meta.last_updated", stateCode["StateCode"])
	lastUpdatedState, lastUpdatedStateErr := jmespath.Search(lastUpdatedStateQuery, sourceData)
	if lastUpdatedStateErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: lastUpdatedStateErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}
	casesStateQuery = fmt.Sprintf("%s.total", stateCode["StateCode"])
	casesState, casesStateErr := jmespath.Search(casesStateQuery, sourceData)
	if casesStateErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: casesStateErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// jquery using jmespath to get India's total cases, vaccination and last updated details
	lastUpdatedIndia, lastUpdatedIndiaErr := jmespath.Search(lastUpdatedIndiaQuery, sourceData)
	if lastUpdatedIndiaErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: lastUpdatedIndiaErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	casesIndia, casesIndiaErr := jmespath.Search(casesIndiaQuery, sourceData)
	if casesIndiaErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: casesIndiaErr.Error(), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// modifying jquery results to specific response model
	jsonStrState, err := json.Marshal(casesState)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: "error", CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStrState, &caseTotalState); err != nil {
		fmt.Println(err)
	}

	_, isError := stateCode["Error"]
	if isError == true {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: stateCode["Error"].(string), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	caseTotalState.LastUpdated = lastUpdatedState.(string)
	caseTotalState.Location = stateCode["State"].(string)

	jsonStrIndia, err := json.Marshal(casesIndia)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: "error", CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStrIndia, &caseTotalIndia); err != nil {
		fmt.Println(err)
	}

	caseTotalIndia.LastUpdated = lastUpdatedIndia.(string)
	caseTotalIndia.Location = "India"

	return c.JSON(http.StatusOK, responses.CaseResponse{Status: http.StatusOK, Message: "success", CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
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
