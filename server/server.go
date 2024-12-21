package server

import (
	"ronbunkun/config"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo   *echo.Echo
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		Config: cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
