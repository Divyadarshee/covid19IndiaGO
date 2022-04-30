package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "go-swag-sample/docs/echosimple" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	"go-swag-sample/echosimple/configs"
	"go-swag-sample/echosimple/routes"
)

// @title Covid19 Cases and Vaccinations in India
// @version 1.0
// @description This is a covid19 cases data server which when given the GPS coordinates of a location returns
// @description the cases details as in confirmed, deceased, recovered, teseted along with vaccination details as in
// @description single and double dose in coordinates provided State and in India in total
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// run database
	configs.ConnectDB()

	//routes
	routes.CaseRoute(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error { // GET function to the route = "/" path and an handler
		return c.JSON(200, &echo.Map{"data": "Get Covid cases details here"}) // function that returns a JSON of "Hello from Echo & mongoDB".
		// echo.Map is a shortcut for map[string]interface{} useful for JSON returns
	})
	port := fmt.Sprintf(":%s", configs.GetPort())
	e.Logger.Fatal(e.Start(port)) // Start function is used to run the application on port 6000

}
