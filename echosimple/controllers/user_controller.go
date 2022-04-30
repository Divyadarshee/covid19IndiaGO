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
// @Description Get the covid19 cases details for the given GPS location.
// @Tags user
// @Accept */*
// @Produce json
// @Param latitude query float64 true "Latitude"
// @Param longitude query float64 true "Longitude"
// @Success 200 {object} map[string]interface{}
// @Router /cases [get]
func GetCases(c echo.Context) error {

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// creates query params binder that stops binding at first error
	err := echo.QueryParamsBinder(c).
		Float64("latitude", &models.GPS.Latitude).
		Float64("longitude", &models.GPS.Longitude).
		BindError() // returns first binding error
	//defer cancel()

	var stateCode echo.Map
	stateCode = GetStateFromGPS(models.GPS.Latitude, models.GPS.Longitude)

	//fmt.Println(GetStateFromGPS(22.5499978, 88.333332))
	//fmt.Println(GetStateFromGPS(20.397373, 72.832802))
	//var state_code string = GetStateFromGPS(20.397373, 72.832802)

	url := "https://data.covid19india.org/v4/min/data.min.json"
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, reerr := http.NewRequest("GET", url, nil)
	if reerr != nil {
		fmt.Println("error")
	}

	resp, reserr := client.Do(req)
	fmt.Println(reserr)
	body, readErr := ioutil.ReadAll(resp.Body)
	fmt.Println(readErr)
	var data map[string]interface{}

	aerr := json.Unmarshal(body, &data)
	fmt.Println(aerr)
	var lastUpdatedStateQuery string = fmt.Sprintf("%s.meta.last_updated", stateCode["StateCode"])
	lastUpdatedState, err := jmespath.Search(lastUpdatedStateQuery, data)
	var casesStateQuery string = fmt.Sprintf("%s.total", stateCode["StateCode"])
	casesState, err := jmespath.Search(casesStateQuery, data)
	var lastUpdatedIndiaQuery string = "TT.meta.last_updated"
	lastUpdatedIndia, err := jmespath.Search(lastUpdatedIndiaQuery, data)
	var casesIndiaQuery string = "TT.total"
	casesIndia, err := jmespath.Search(casesIndiaQuery, data)
	fmt.Println(err)
	//fmt.Println(lastUpdatedState)
	//fmt.Println(casesState)
	var caseTotalState responses.Cases
	var caseTotalIndia responses.Cases
	jsonStrState, err := json.Marshal(casesState)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: "error", CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStrState, &caseTotalState); err != nil {
		fmt.Println(err)
	}
	_, isPresent := stateCode["Error"]
	if isPresent == true {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: stateCode["Error"].(string), CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}
	caseTotalState.LastUpdated = lastUpdatedState.(string)
	caseTotalState.Location = stateCode["State"].(string)
	//fmt.Println(lastUpdatedIndia)
	//fmt.Println(casesIndia)

	jsonStrIndia, err := json.Marshal(casesIndia)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.CaseResponse{Status: http.StatusInternalServerError, Message: "error", CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStrIndia, &caseTotalIndia); err != nil {
		fmt.Println(err)
	}
	fmt.Println(caseTotalIndia.Vaccinated2)
	caseTotalIndia.LastUpdated = lastUpdatedIndia.(string)
	caseTotalIndia.Location = "India"
	return c.JSON(http.StatusOK, responses.CaseResponse{Status: http.StatusOK, Message: "success", CasesInState: caseTotalState, CasesInIndia: caseTotalIndia})
}

func GetStateFromGPS(lat float64, long float64) echo.Map {
	url := fmt.Sprintf("https://eu1.locationiq.com/v1/reverse.php?key=pk.c8772d5c9e1d9046e5995be8c9edcaa4&lat=%f&lon=%f&format=json", lat, long)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return echo.Map{"Error": err.Error()}
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return echo.Map{"Error": err.Error()}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return echo.Map{"Error": err.Error()}
	}
	fmt.Println(string(body))
	//fmt.Printf("type of body is %T", body)
	//fmt.Println()
	var data map[string]interface{}

	aerr := json.Unmarshal(body, &data)
	fmt.Println(data)
	fmt.Println(aerr)

	result, err := jmespath.Search("address.state", data)
	fmt.Println(err)
	fmt.Println(result)
	for key, value := range state_data.StateCodes {
		if result == key {
			fmt.Printf("%q is the key for the value %q\n", key, value)
			return echo.Map{"State": key, "StateCode": value}
		}
	}
	return nil
}

//import (
//	"context"
//	"github.com/go-playground/validator/v10"
//	"github.com/labstack/echo/v4"
//	"go-swag-sample/echosimple/configs"
//	"go-swag-sample/echosimple/models"
//	"go-swag-sample/echosimple/responses"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"net/http"
//	"time"
//)

//var userCollection *mongo.Collection = configs.GetCollections(configs.DB, "users")
//var validate = validator.New()
//
//// CreateUser godoc
//// swagger:route GET /admin/company/
//// @Summary Create a user.
//// @Description Create a user with name, location and title.
//// @Tags user
//// @Accept */*
//// @Produce json
//// @Success 201 {object} map[string]interface{}
//// @Router /user [post]
//func CreateUser(c echo.Context) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	var user models.User
//	defer cancel()
//
//	//validate the request body
//	if err := c.Bind(&user); err != nil {
//		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	//use validator library to validate required fields
//	if validationErr := validate.Struct(&user); validationErr != nil {
//		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
//	}
//
//	newUser := models.User{
//		Id:       primitive.NewObjectID(),
//		Name:     user.Name,
//		Location: user.Location,
//		Title:    user.Title,
//	}
//
//	result, err := userCollection.InsertOne(ctx, newUser)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
//}
//
//// GetAUser godoc
//// @Summary Get a user.
//// @Description get the user for the given userid.
//// @Tags user
//// @Accept */*
//// @Produce json
//// @Success 200 {object} map[string]interface{}
//// @Router /user/{userId} [get]
//func GetAUser(c echo.Context) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	userId := c.Param("userId")
//	var user models.User
//	defer cancel()
//
//	objId, _ := primitive.ObjectIDFromHex(userId)
//
//	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
//
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
//}
//
//// EditAUser godoc
//// @Summary Edit a user details.
//// @Description Edit a user's name, location and/or title.
//// @Tags user
//// @Accept */*
//// @Produce json
//// @Success 200 {object} map[string]interface{}
//// @Router /user/{userId} [put]
//func EditAUser(c echo.Context) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	userId := c.Param("userId")
//	var user models.User
//	defer cancel()
//
//	objId, _ := primitive.ObjectIDFromHex(userId)
//
//	//validate the request body
//	if err := c.Bind(&user); err != nil {
//		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	//use validator library to validate required fields
//	if validationErr := validate.Struct(&user); validationErr != nil {
//		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
//	}
//
//	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
//
//	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
//
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	//get update user details
//	var updatedUser models.User
//	if result.MatchedCount == 1 {
//		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(updatedUser)
//
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//		}
//	}
//	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
//}
//
//// DeleteAUser godoc
//// @Summary Delete a user.
//// @Description Delete a user.
//// @Tags user
//// @Accept */*
//// @Produce json
//// @Success 200 {object} map[string]interface{}
//// @Router /user/{userId} [delete]
//func DeleteAUser(c echo.Context) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	userId := c.Param("userId")
//	defer cancel()
//
//	objId, _ := primitive.ObjectIDFromHex(userId)
//
//	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
//
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": result}})
//}
//
//// GetAllUsers godoc
//// @Summary Show all the users.
//// @Description get the details of all users.
//// @Tags root
//// @Accept */*
//// @Produce json
//// @Success 200 {object} map[string]interface{}
//// @Router /users [get]
//func GetAllUsers(c echo.Context) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	var users []models.User
//	defer cancel()
//
//	results, err := userCollection.Find(ctx, bson.M{})
//
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//
//	defer results.Close(ctx)
//	for results.Next(ctx) {
//		var singleUser models.User
//		if err = results.Decode(&singleUser); err != nil {
//			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
//		}
//
//		users = append(users, singleUser)
//	}
//
//	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
//}
