package server

import (
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
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Published  string   `json:"published"`
	Pdfurl     string   `json:"pdfurl"`
	Categories []string `json:"categories"`
}

func generate(c echo.Context) error {
	log.Print("Generate request received")
	ctx := c.Request().Context()
	resultChan, cancel, err := arxiv.Search(ctx, &arxiv.Query{
		Terms:         "deep learning",
		MaxPageNumber: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	var articles []Article

	for resultPage := range resultChan {
		if err := resultPage.Err; err != nil {
			continue
		}

		feed := resultPage.Feed

		for _, entry := range feed.Entry {
			categories := []string{}
			for _, category := range entry.Category {
				categories = append(categories, string(category.Term))
			}

			articles = append(articles, Article{
				ID:         entry.ID,
				Title:      entry.Title,
				Published:  string(entry.Updated),
				Pdfurl:     entry.Link[1].Href,
				Categories: categories,
			})
		}
		if resultPage.PageNumber >= 1 {
			cancel()
		}
	}
	return c.JSON(http.StatusOK, articles)
}

func healthcheck(c echo.Context) error {
	log.Print("Healthcheck request received")
	return c.JSON(http.StatusOK, map[string]string{"status": "RUNNING"})
}
