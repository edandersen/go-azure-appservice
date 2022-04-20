package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Hello from Go and Gin",
			"link":  "/json",
		})
	})

	router.GET("/json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"foo": "bar",
		})
	})

	// Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
