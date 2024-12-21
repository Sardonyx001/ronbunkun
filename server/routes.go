package server

import (
	"bytes"
	"fmt"
	"net/http"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marvin-hansen/arxiv/v1"
)

func ConfigureRoutes(server *Server) {
	server.Echo.Use(middleware.Recover())
	server.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: server.Echo.Logger.Output(),
	}))

	server.Echo.GET("/health", healthcheck)

	api := server.Echo.Group("/api")

	api.GET("/generate", generate)
}

type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

func generate(c echo.Context) error {
	log.Print("Generate request received")
	ctx := c.Request().Context()
	resultChan, cancel, err := arxiv.Search(ctx, &arxiv.Query{
		Terms:         "deep learning",
		MaxPageNumber: 5,
	})
	if err != nil {
		log.Fatal(err)
	}

	var responseBuf bytes.Buffer

	var articles []Article

	for resultPage := range resultChan {
		if err := resultPage.Err; err != nil {
			fmt.Fprintf(&responseBuf, "#%d err: %v", resultPage.PageNumber, err)
			continue
		}

		fmt.Fprintf(&responseBuf, "#%d\n", resultPage.PageNumber)
		feed := resultPage.Feed
		fmt.Fprintf(&responseBuf, "\tTitle: %s\n\tID: %s\n\tAuthor: %#v\n\tUpdated: %#v\n", feed.Title, feed.ID, feed.Author, feed.Updated)

		for i, entry := range feed.Entry {
			fmt.Fprintf(&responseBuf, "\n\t\tEntry: #%d Title: %s ID: %s\n\t\tSummary: %s\n\t\tContent: %#v\n\t\tUpdated: %#v\n\t\tLinks: %#v\n",
				i, entry.Title, entry.ID, entry.Summary.Body, entry.Content, entry.Updated, entry.Link,
			)
		}
		if resultPage.PageNumber >= 2 {
			cancel()
		}
	}
	return c.JSON(http.StatusOK, articles)
}

func healthcheck(c echo.Context) error {
	log.Print("Healthcheck request received")
	return c.JSON(http.StatusOK, map[string]string{"status": "RUNNING"})
}
