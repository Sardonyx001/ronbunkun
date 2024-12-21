package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureRoutes(server *Server) {
	server.Echo.Use(middleware.Recover())
	server.Echo.Pre(middleware.RemoveTrailingSlash())
}

func todo(c echo.Context) error {
	return c.JSON(http.StatusOK, "Path: "+c.Path())
}
