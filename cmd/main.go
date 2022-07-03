package main

import (
	"github.com/DipandaAser/tg-bot-storage/internal/config"
	"github.com/DipandaAser/tg-bot-storage/internal/handlers"
	internalManager "github.com/DipandaAser/tg-bot-storage/internal/manager"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	err := internalManager.InitManager(config.GetDefaultConfig().Tokens...)
	if err != nil {
		log.Fatal(err)
		return
	}

	go internalManager.GetDefaultManager().StartUploaderManager()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "api-key"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
	}))

	//Health check endpoint
	router.GET("/health", handlers.Health)

	apiRouter := router.Group("/api")
	apiRouter.Use(handlers.ApiKeyMiddleware)

	//Files
	apiRouter.GET("/files", handlers.DownloadFile)
	apiRouter.POST("/files", handlers.UploadFile)

	//Bot management
	apiRouter.POST("/bot", handlers.AddBotToManager)

	if err = router.Run(":7000"); err != nil {
		log.Fatal(err)
	}
}
