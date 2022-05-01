package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmespath/go-jmespath"
	state_data "go-swag-sample/echosimple/data"
	"go-swag-sample/echosimple/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ping database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the DataBase")
	return client
}

// DB Client instance
var DB *mongo.Client = ConnectDB()

// getting Database collections
func GetCollections(client *mongo.Client, collectionString string) *mongo.Collection {
	collection := client.Database("golangdb").Collection(collectionString)
	return collection
}

// clearing any existing collections before populating
func ClearCollections() {
	var userCollection *mongo.Collection = GetCollections(DB, "cases")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection.Drop(ctx)
	fmt.Println("Collections cleared")
}

// populating the Database with covid19 details
func PopulateDB() {
	var caseCollection *mongo.Collection = GetCollections(DB, "cases")
	var sourceData map[string]interface{}
	var lastUpdatedStateQuery string
	var casesStateQuery string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var cases models.Cases
	defer cancel()

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
		fmt.Println(sourceResponseErr.Error())
	}

	// reading the body of the response
	body, readErr := ioutil.ReadAll(sourceResponse.Body)
	if readErr != nil {
		fmt.Println(readErr.Error())
	}

	// parsing json
	parsingErr := json.Unmarshal(body, &sourceData)
	if parsingErr != nil {
		fmt.Println(parsingErr.Error())
	}

	for name, code := range state_data.StateCodes {
		lastUpdatedStateQuery = fmt.Sprintf("%s.meta.last_updated", code)
		lastUpdatedState, lastUpdatedStateErr := jmespath.Search(lastUpdatedStateQuery, sourceData)
		if lastUpdatedStateErr != nil {
			fmt.Println(lastUpdatedStateErr.Error())
		}
		casesStateQuery = fmt.Sprintf("%s.total", code)
		casesState, casesStateErr := jmespath.Search(casesStateQuery, sourceData)
		if casesStateErr != nil {
			fmt.Println(casesStateErr.Error())
		}

		// modifying jquery results to specific response model
		jsonStrState, marshalErr := json.Marshal(casesState)
		if marshalErr != nil {
			fmt.Println(marshalErr.Error())
		}

		// Convert json string to struct
		if unmarshalErr := json.Unmarshal(jsonStrState, &cases); unmarshalErr != nil {
			fmt.Println(unmarshalErr)
		}

		cases.StateName = name
		cases.StateCode = code
		cases.LastUpdated = lastUpdatedState.(string)

		result, InsertErr := caseCollection.InsertOne(ctx, cases)

		if InsertErr != nil {
			fmt.Println(InsertErr.Error())
		}

		//fmt.Println(result)
	}
}
