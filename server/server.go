package server

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"virunus.com/cornhole/config"
)

var (
	//go:embed templates
	templates embed.FS

	//go:embed static
	static embed.FS
)

func Serve(Configuration *config.Config) {
	t, err := loadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.SetHTMLTemplate(t)

	initializeRoutes(router)

	err = router.Run(fmt.Sprint(`localhost:`, Configuration.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func showIndexPage(c *gin.Context)  {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "Home Page",
		},
	)
}

func initializeRoutes(router *gin.Engine) {
	router.StaticFS("/static", http.FS(static)) // static files
	router.GET("/favicon.ico", func(context *gin.Context) { // custom favicon handler
		context.FileFromFS("/static/favicon.ico", http.FS(static))
	})

	router.GET("/", showIndexPage)
}

func loadTemplates() (*template.Template, error) {
	parsedTemplate, err := template.ParseFS(templates, "**/*.html")
	if err != nil {
		return nil, err
	}

	return parsedTemplate, nil
}
