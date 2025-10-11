package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/high-effort-low-stress/go-bank-api/controllers"
	"github.com/high-effort-low-stress/go-bank-api/database"
	"github.com/high-effort-low-stress/go-bank-api/repositories"
	"github.com/high-effort-low-stress/go-bank-api/services"
	"github.com/joho/godotenv"
)

var PORT_ENV = "PORT"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	db := database.DB

	emailService, err := services.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize EmailService: %v", err)
	}

	// Dependencies
	onboardingRequestRepository := repositories.NewOnboardingRequestRepository(db)
	onboardingService := services.NewOnboardingService(onboardingRequestRepository, emailService)
	onboardingController := controllers.NewOnboardingController(onboardingService)

	server := gin.Default()

	apiV1 := server.Group("/api/v1")
	{
		onboarding := apiV1.Group("/onboarding")
		{
			onboarding.POST("/start", onboardingController.StartOnboarding)
		}
	}

	server.Run(os.Getenv(PORT_ENV))
}
