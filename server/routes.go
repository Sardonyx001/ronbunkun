package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marvin-hansen/arxiv/v1"

	"github.com/samber/lo"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

func ConfigureRoutes(server *Server) {
	server.Echo.Use(middleware.Recover())
	server.Echo.Use(middleware.CORS())
	server.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: server.Echo.Logger.Output(),
	}))

	server.Echo.GET("/health", healthcheck)

	api := server.Echo.Group("/api")

	api.GET("/translate", llmHandler)
	api.POST("/generate", generate)
}

func healthcheck(c echo.Context) error {
	log.Print("Healthcheck request received")
	return c.JSON(http.StatusOK, map[string]string{"status": "RUNNING"})
}

type Article struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Published  string   `json:"published"`
	Pdfurl     string   `json:"pdfurl"`
	Categories []string `json:"categories"`
	Summary    string   `json:"summary"`
}

type GenerateData struct {
	Keywords  []string `json:"keywords"`
	MaxPapers int      `json:"maxPapers"`
}

func generate(c echo.Context) error {
	log := c.Logger()

	log.Print("Generate request received")
	data := new(GenerateData)
	if err := c.Bind(data); err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	searchFields := lo.Map(data.Keywords, func(keyword string, _ int) *arxiv.Field {
		return &arxiv.Field{All: keyword}
	})

	ctx := c.Request().Context()
	resultChan, cancel, err := arxiv.Search(ctx, &arxiv.Query{
		Filters: []*arxiv.Filter{
			{
				Op:     arxiv.OpOR,
				Fields: searchFields,
			},
		},
		MaxPageNumber: 5,
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

		for i, entry := range feed.Entry {

			categories := lo.Map(entry.Category, func(cat *arxiv.Class, idx int) string {
				return string(cat.Term)
			})

			articles = append(articles, Article{
				ID:         entry.ID,
				Title:      entry.Title,
				Published:  string(entry.Updated),
				Pdfurl:     entry.Link[1].Href,
				Categories: categories,
				Summary:    entry.Summary.Body,
			})

			if i >= data.MaxPapers-1 {
				cancel()
				break
			}
		}
		if resultPage.PageNumber >= 5 || len(articles)-1 >= data.MaxPapers {
			cancel()
		}
	}

	urls := lo.Map(articles, func(article Article, _ int) string {
		return article.Pdfurl
	})

	pdfs, err := loadPDFs(urls)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	_ = pdfs

	return c.JSON(http.StatusOK, articles)
}

func loadPDFs(urls []string) (map[string][]byte, error) {
	pdfs := make(map[string][]byte)
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			log.Print("did not find : " + url)
			continue
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download file: %s", url)
		}

		content, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		pdfs[url] = content
		log.Printf("Loaded PDF from %s", url)
	}
	return pdfs, nil
}

type LLMRequest struct {
	Model  string `json:"model"`
	System string `json:"system"`
	Prompt string `json:"prompt"`
}

type LLMResponse struct {
	Text string `json:"text"`
}

func llmHandler(c echo.Context) error {
	system := "You are an ai translator that translates english papers into japanese. Print only the translation."
	prompt := "Translate the following: "
	textToTranslate := "Because Temporal Logic (TL) is typically used for detecting time-based complex patterns over streams in real time, we aimed at taking advantage of TL after reconstructing records in log files. We have selected Apache web server logs as a case study since Web servers are one of the most preferred targets for hackers and other cyber criminals to intrude systems because of their publicity. Web logs are set of timely recorded events occurred between web servers and clients. In general, Web log files keep each record in the form of request and response together in one line. Reconstructing web server activities as streams based on records in web log files gives us the capability of implementing TL based on streaming data. This way it becomes trivial to achieve fast and better forensics investigation because there are state- of-the-art streaming technologies such as StreamBase, Esper, etc. (“EsperTech - Esper,” n.d., “StreamBase | Complex Event Processing, Event Stream Processing, StreamBase Streaming Platform,” n.d.). We preferred Esper platform and Event Processing Language (EPL) as a standard language to define misuse patterns. Esper provides .NET and Java packages that are easy to implement either for a standalone application or enterprise framework along with EPL, a declarative language for dealing with high frequency time- based event data (“EsperTech - Esper,” n.d.). To the best of our knowledge, there is no platform performing post-mortem log analysis using MSFOMTL, EPL, and reconstruction approach. In addition, cyber security professionals and investigators lack a standard format or language to store and share their previous experiences on log analysis. Besides performance advantage and temporal logical capabilities, our approach would base a platform and a library to store, share, and adjust previously identified patterns of misuses for further analysis. Through the paper, we describe previous work in Section 2, along with some background information in Section 3. Section 4 describes misuse patterns informally that could reside and could be detected from web log records. Then, in Section 5, we define formal versions of those patterns using a special case of TL. Section 6 describes EPL queries that are mapped from TL formulae given in the previous section."

	s, err := llmTranslateAndSummarize(c, system, prompt, textToTranslate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, s)
}

func llmTranslateAndSummarize(c echo.Context, system string, prompt string, text string) ([]string, error) {
	log := c.Logger()
	log.SetOutput(os.Stdout)

	llm, err := ollama.New(ollama.WithModel("schroneko/gemma-2-2b-jpn-it"))
	if err != nil {
		return nil, err
	}
	log.Info("llm instance created")

	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	completion, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0))
	if err != nil {
		return nil, err
	}
	log.Info("llm content generated")

	// Using chains
	translatePrompt := prompts.NewPromptTemplate(
		system+"Translate the following text from {{.inputLanguage}} to {{.outputLanguage}}. {{.text}}",
		[]string{"inputLanguage", "outputLanguage", "text"},
	)
	llmChain := chains.NewLLMChain(llm, translatePrompt)

	outputValues, err := chains.Call(ctx, llmChain, map[string]any{
		"inputLanguage":  "English",
		"outputLanguage": "Japanese",
		"text":           text,
	})
	if err != nil {
		return nil, err
	}
	log.Info("llm chain called")

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		return nil, fmt.Errorf("invalid chain return")
	}

	return []string{completion.Choices[0].Content, out}, err
}
