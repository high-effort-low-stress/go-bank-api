package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/high-effort-low-stress/go-bank-api/internal/database"
	"github.com/high-effort-low-stress/go-bank-api/internal/notification"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/controllers"
	onboarding_repositories "github.com/high-effort-low-stress/go-bank-api/internal/onboarding/repositories"
	onboarding_services "github.com/high-effort-low-stress/go-bank-api/internal/onboarding/services"
	user_repositories "github.com/high-effort-low-stress/go-bank-api/internal/users/repositories"
	user_services "github.com/high-effort-low-stress/go-bank-api/internal/users/services"
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

	emailService, err := notification.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize EmailService: %v", err)
	}

	// Dependencies
	onboardingRequestRepository := onboarding_repositories.NewOnboardingRequestRepository(db)
	userRepository := user_repositories.NewUserRepository(db)

	onboardingService := onboarding_services.NewOnboardingService(onboardingRequestRepository, emailService, nil)
	verifyEmailTokenService := onboarding_services.NewVerifyEmailTokenService(onboardingRequestRepository)
	createUserService := user_services.NewCreateUserService(userRepository)
	completeOnboardingService := onboarding_services.NewCompleteOnboardingService(onboardingRequestRepository, createUserService)
	onboardingController := controllers.NewOnboardingController(onboardingService, verifyEmailTokenService, completeOnboardingService)

	server := gin.Default()

	apiV1 := server.Group("/api/v1")
	{
		onboarding := apiV1.Group("/onboarding")
		{
			onboarding.POST("/start", onboardingController.StartOnboarding)
			onboarding.POST("/verify", onboardingController.VerifyEmail)
			onboarding.POST("/complete", onboardingController.CompleteOnboarding)
		}
	}

	server.Run(os.Getenv(PORT_ENV))
}
