package server

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"virunus.com/cornhole/config"
)

var (
	//go:embed templates
	templates embed.FS

	//go:embed assets
	staticFs embed.FS
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func embedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func Serve(Configuration *config.Config) {
	t, err := loadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.SetHTMLTemplate(t)

	// set up dynamic routes
	initializeRoutes(router)

	// set up assets files
	router.NoRoute(static.Serve("/", embedFolder(staticFs, "assets/webroot")))

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
	router.GET("/", showIndexPage)
}

func loadTemplates() (*template.Template, error) {
	parsedTemplate, err := template.ParseFS(templates, "**/*.html")
	if err != nil {
		return nil, err
	}

	return parsedTemplate, nil
}
