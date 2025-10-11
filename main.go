package main

import (
	"github.com/gin-gonic/gin"
	"github.com/high-effort-low-stress/go-bank-api/controllers"
)

var PORT = ":8080"

func main() {
	server := gin.Default()

	apiV1 := server.Group("/api/v1")
	{
		onboarding := apiV1.Group("/onboarding")
		{
			onboarding.POST("/start", controllers.StartOnboarding)
		}
	}

	server.Run(PORT)
}
