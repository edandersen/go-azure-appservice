package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/fsnotify.v1"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Hello from Go and Gin running on Azure App Service",
			"link":  "/json",
		})
	})

	router.GET("/json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"foo": "bar",
		})
	})

	router.Static("/public", "./public")

	// creates a new file watcher for App_offline.htm
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	// watch for App_offline.htm and exit the program if present
	// This allows continuous deployment on App Service as the .exe will not be
	// terminated otherwise
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if strings.HasSuffix(event.Name, "app_offline.htm") {
					fmt.Println("Exiting due to app_offline.htm being present")
					os.Exit(0)
				}
			}
		}
	}()

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err := watcher.Add(currentDir); err != nil {
		fmt.Println("ERROR", err)
	}

	// Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
