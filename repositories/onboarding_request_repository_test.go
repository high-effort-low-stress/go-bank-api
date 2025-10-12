package repositories_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/models"
	"github.com/high-effort-low-stress/go-bank-api/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	repo repositories.OnboardingRequestRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "public"
	dbUser := "postgres"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17.5-alpine",
		postgres.WithInitScripts(filepath.Join("..", "sql", "migrations", "01_create-onboarding.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Obtém a string de conexão dinâmica do contêiner
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "options=--search_path%3Donboarding")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// Conecta ao banco de dados e inicializa o repositório
	db, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	repo = repositories.NewOnboardingRequestRepository(db)

	os.Exit(m.Run())
}

func TestOnboardingRequestRepository(t *testing.T) {
	// Dados base para os testes
	requestToCreate := &models.OnboardingRequest{
		FullName:              "Integration Test User",
		Email:                 "integration@example.com",
		DocumentNumber:        "98765432100",
		VerificationTokenHash: "integration-test-hash",
		TokenExpiresAt:        time.Now().Add(1 * time.Hour),
		Status:                models.StatusPending,
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(requestToCreate)

		require.NoError(t, err)
		assert.NotEmpty(t, requestToCreate.PublicID, "PublicID should be set by BeforeCreate hook")
	})

	t.Run("FindByDocumentOrEmail", func(t *testing.T) {
		// Depende do sucesso do teste "Create"
		require.NotEmpty(t, requestToCreate.PublicID, "Cannot run Find test if Create failed")

		// Act
		foundRequest, err := repo.FindByDocumentOrEmail(requestToCreate.DocumentNumber, "some-other-email@test.com")
		foundRequest2, err := repo.FindByDocumentOrEmail("1234567810", requestToCreate.Email)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, requestToCreate.Email, foundRequest.Email)
		assert.Equal(t, requestToCreate.FullName, foundRequest.FullName)
		assert.Equal(t, requestToCreate.PublicID, foundRequest.PublicID)

		assert.Equal(t, requestToCreate.Email, foundRequest2.Email)
		assert.Equal(t, requestToCreate.FullName, foundRequest2.FullName)
		assert.Equal(t, requestToCreate.PublicID, foundRequest2.PublicID)
	})
}
