package router

import (
	"net/http"
	"tenant/internal/container"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "tenant/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init(e *echo.Echo, hc *container.HandlerComponent, sc *container.SharedComponent) {
	e.Pre(middleware.Rewrite(map[string]string{
		"/v1/*": "/$1",
	}))

	// Utility endpoints
	e.GET("/docs/index.html", echoSwagger.WrapHandler)
	e.GET("/docs/doc.json", echoSwagger.WrapHandler)
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/ping", ping)
	publicRouter(e, hc)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "X-User-ID", "x-token", "X-Requested-With", "x-device-id", "x-service-authorization", "x-service-timestamp"},
		AllowCredentials: true,
	}))
}

// ping write pong to http.ResponseWriter.
func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
