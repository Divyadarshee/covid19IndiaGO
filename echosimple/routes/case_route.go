package routes

import (
	"github.com/labstack/echo/v4"
	"go-swag-sample/echosimple/controllers"
)

func CaseRoute(e *echo.Echo) {

	e.GET("/cases", controllers.GetCases)
}
