package routes

import (
	"github.com/labstack/echo/v4"
	"go-swag-sample/echosimple/controllers"
)

func CaseRoute(e *echo.Echo) {

	e.GET("/cases", controllers.GetCases)

	// All routes related to users comes here
	//e.POST("/user", controllers.CreateUser)
	//
	//e.GET("/user/:userId", controllers.GetAUser)
	//
	//e.PUT("/user/:userId", controllers.EditAUser)
	//
	//e.DELETE("/user/:userId", controllers.DeleteAUser)
	//
	//e.GET("/users", controllers.GetAllUsers)
}
