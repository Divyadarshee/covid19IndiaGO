package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
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

var pool = configs.CreateNewPool()
var caseCollection *mongo.Collection = configs.GetCollections(configs.DB, "cases")
var validate = validator.New()

// GetCases godoc
// @Summary Get latest covid19 cases.
// @Description Get the covid19 cases details for the given GPS coordinates.
// @Tags cases
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

	//return 400 if coordinates are out of India
	_, isError := stateCode["Error"]
	if isError == true {
		return c.JSON(http.StatusBadRequest, responses.CaseResponse{Status: http.StatusBadRequest, Message: stateCode["Error"].(string), CasesInState: stateCases, CasesInIndia: indiaCases})
	}

	//redis
	conn := pool.Get()
	defer conn.Close()

	//check if state cache exists
	if CheckCache(conn, stateCode["StateCode"].(string)) != true {
		fmt.Printf("Fetching data for state code: %s from DB \n", stateCode["StateCode"])
		//fetching data from DB
		stateDBQueryErr := caseCollection.FindOne(ctx, bson.M{"statecode": stateCode["StateCode"]}).Decode(&stateCases)
		if stateDBQueryErr != nil {
			return c.JSON(http.StatusNotFound, responses.CaseResponse{Status: http.StatusNotFound, Message: stateDBQueryErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
		}

		//caching
		SetCache(conn, stateCode["StateCode"].(string), stateCases)

	} else {
		stateCache, stateCacheErr := RetrieveCache(conn, stateCode["StateCode"].(string), stateCases)
		if stateCacheErr != nil {
			fmt.Printf("Error while retrieving cache: %s \n", stateCacheErr)

			//fetching data from DB
			stateDBQueryErr := caseCollection.FindOne(ctx, bson.M{"statecode": stateCode["StateCode"]}).Decode(&stateCases)
			if stateDBQueryErr != nil {
				return c.JSON(http.StatusNotFound, responses.CaseResponse{Status: http.StatusNotFound, Message: stateDBQueryErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
			}

			//caching
			SetCache(conn, stateCode["StateCode"].(string), stateCache)

		} else {
			stateCases = stateCache
		}
	}

	//check if India cache exists
	if CheckCache(conn, "TT") != true {
		fmt.Println("Fetching data for state code: TT from DB")
		//fetching data from DB
		indiaDBQueryErr := caseCollection.FindOne(ctx, bson.M{"statecode": "TT"}).Decode(&indiaCases)
		if indiaDBQueryErr != nil {
			return c.JSON(http.StatusNotFound, responses.CaseResponse{Status: http.StatusNotFound, Message: indiaDBQueryErr.Error(), CasesInState: stateCases, CasesInIndia: indiaCases})
		}

		//caching
		SetCache(conn, "TT", indiaCases)

	} else {
		indiaCache, indiaCacheErr := RetrieveCache(conn, "TT", indiaCases)
		if indiaCacheErr != nil {
			fmt.Printf("Error while retrieving cache: %s \n", indiaCacheErr)

			//fetching data from DB
			indiaDBQueryErr := caseCollection.FindOne(ctx, bson.M{"statecode": "TT"}).Decode(&indiaCases)
			if indiaDBQueryErr != nil {
				return c.JSON(http.StatusNotFound, responses.CaseResponse{Status: http.StatusNotFound, Message: indiaDBQueryErr.Error(), CasesInState: stateCases, CasesInIndia: stateCases})
			}

			//caching
			SetCache(conn, "TT", indiaCases)
		} else {
			indiaCases = indiaCache
		}
	}

	return c.JSON(http.StatusOK, responses.CaseResponse{Status: http.StatusOK, Message: "success", CasesInState: stateCases, CasesInIndia: indiaCases})
}

func CheckCache(conn redis.Conn, stateCode string) bool {

	CacheExists, CacheExistsErr := redis.Bool(conn.Do("EXISTS", stateCode))
	if CacheExistsErr != nil {
		fmt.Printf("Error checking if key %s exists: %v \n", stateCode, CacheExistsErr)
	}
	return CacheExists
}

func SetCache(conn redis.Conn, stateCode string, stateCases models.Cases) {
	//caching data for further queries
	_, createStateCacheErr := conn.Do("HSET", redis.Args{}.Add(stateCode).AddFlat(stateCases)...)
	if createStateCacheErr != nil {
		fmt.Printf("Error caching data: %s \n", createStateCacheErr.Error())
	}

	//setting cache expiry
	_, setStateCacheExipryErr := conn.Do("EXPIRE", stateCode, "1800")
	if setStateCacheExipryErr != nil {
		fmt.Printf("Error while setting expiry for state code %s cache: %s", stateCode, setStateCacheExipryErr.Error())
	}
}

func ConvertCachetoStruct(stateCode string, stateCases models.Cases, stateCache []interface{}) models.Cases {

	fmt.Printf("Fetching data for state code: %s from cache \n", stateCode)
	if err := redis.ScanStruct(stateCache, &stateCases); err != nil {
		fmt.Println(err.Error())
	}

	return stateCases
}

func RetrieveCache(conn redis.Conn, stateCode string, stateCases models.Cases) (models.Cases, error) {
	//retrieving cached data
	stateCache, stateCacheErr := redis.Values(conn.Do("HGETALL", stateCode))
	if stateCacheErr != nil {
		return stateCases, stateCacheErr
	} else {
		stateCases = ConvertCachetoStruct(stateCode, stateCases, stateCache)
	}
	return stateCases, nil
}

func GetStateFromGPS(latitude float64, longitude float64) echo.Map {
	reverseGeocodingUrl := fmt.Sprintf(state_data.ReverseGeocodingUrl, latitude, longitude)
	method := "GET"

	client := &http.Client{}
	reverseGeocodingRequest, reverseGeocodingRequestErr := http.NewRequest(method, reverseGeocodingUrl, nil)

	if reverseGeocodingRequestErr != nil {
		return echo.Map{"Error": reverseGeocodingRequestErr.Error()}
	}
	reverseGeocodingResponse, reverseGeocodingResponseErr := client.Do(reverseGeocodingRequest)
	if reverseGeocodingResponseErr != nil {
		return echo.Map{"Error": reverseGeocodingResponseErr.Error()}
	}
	defer reverseGeocodingResponse.Body.Close()

	// reading the body of the response
	body, ReadErr := ioutil.ReadAll(reverseGeocodingResponse.Body)
	if ReadErr != nil {
		return echo.Map{"Error": ReadErr.Error()}
	}
	fmt.Println(string(body))
	var data map[string]interface{}

	parsingErr := json.Unmarshal(body, &data)
	if parsingErr != nil {
		return echo.Map{"Error": parsingErr.Error()}
	}

	// jquery using jmespath to get state name from the given gps location
	stateName, stateNameErr := jmespath.Search("address.state", data)
	if stateNameErr != nil {
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
