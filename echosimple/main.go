package main

import (
	"fmt"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go-swag-sample/echosimple/configs"
	"go-swag-sample/echosimple/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go-swag-sample/docs/echosimple" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
)

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
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

	// Routes
	// run database
	configs.ConnectDB()
	//e.GET("/", HealthCheck)
	//e.GET("/swagger/*", echoSwagger.WrapHandler)
	//
	//// Start server
	//e.Logger.Fatal(e.Start(":3000"))

	//routes
	routes.CaseRoute(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error { // GET function to the route = "/" path and an handler
		return c.JSON(200, &echo.Map{"data": "Hello from Echo & mongoDB"}) // function that returns a JSON of "Hello from Echo & mongoDB".
		// echo.Map is a shortcut for map[string]interface{} useful for JSON returns
	})
	port := fmt.Sprintf(":%s", configs.GetPort())
	e.Logger.Fatal(e.Start(port)) // Start function is used to run the application on port 6000

}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
