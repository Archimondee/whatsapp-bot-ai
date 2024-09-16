package main

import (
	"62-gpt/app/http/controllers"
	"62-gpt/app/pkg"
	"62-gpt/app/services"
	"62-gpt/config"
	"62-gpt/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthStatus struct {
	ServerStatus string
}

func main() {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
		return
	}
	utils.ConnectDatabase(loadedConfig.DBDriver, loadedConfig.DBSource, "")

	whatsapp := pkg.NewWhatsapp()
	gptServices := pkg.NewGptService()
	whatsappService := services.NewWhatsappService()
	whatsappController := controllers.NewWhatsappController(whatsappService, gptServices, whatsapp)

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.ResponseData("success", "Server running well", &HealthStatus{
			ServerStatus: "ok",
		}))
	})
	apiV1 := r.Group("/api/v1")

	apiV1.GET("/webhook", whatsappController.GetWhatsappWebhook)
	apiV1.POST("/webhook", whatsappController.PostWhatsappWebhook)

	err = r.Run(loadedConfig.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
