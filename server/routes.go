package server

import (
	"net/http"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureRoutes(server *Server) {
	server.Echo.Use(middleware.Recover())
	server.Echo.Pre(middleware.RemoveTrailingSlash())
	server.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: server.Echo.Logger.Output(),
	}))

	server.Echo.GET("/health", healthcheck)
}

func healthcheck(c echo.Context) error {
	log.Print("Healthcheck request received")
	return c.JSON(http.StatusOK, map[string]string{"status": "RUNNING"})
}
