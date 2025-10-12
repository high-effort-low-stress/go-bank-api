package repositories_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/repositories"
	"github.com/high-effort-low-stress/go-bank-api/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db                *gorm.DB
	repo              repositories.OnboardingRequestRepository
	postgresContainer *postgres.PostgresContainer
	ctx               context.Context
)

var sqldir string = filepath.Join("..", "..", "..", "sql")

func TestMain(m *testing.M) {
	ctx = context.Background()
	var err error

	dbName := "public"
	dbUser := "postgres"
	dbPassword := "password"

	postgresContainer, err = postgres.Run(ctx,
		"postgres:17.5-alpine",
		postgres.WithInitScripts(filepath.Join(sqldir, "migrations", "01-create-onboarding.sql")),
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

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "options=--search_path%3Donboarding")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	db, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	repo = repositories.NewOnboardingRequestRepository(db)

	os.Exit(m.Run())
}

func TestOnboardingRequestCreate(t *testing.T) {
	t.Run("#method Create", func(t *testing.T) {
		requestToCreate := &models.OnboardingRequest{
			FullName:              "Integration Test User",
			Email:                 "integration@example.com",
			DocumentNumber:        "98765432100",
			VerificationTokenHash: "integration-test-hash",
			TokenExpiresAt:        time.Now().Add(1 * time.Hour),
			Status:                models.StatusPending,
		}
		err := repo.Create(requestToCreate)

		require.NoError(t, err)
		assert.NotEmpty(t, requestToCreate.PublicID, "PublicID should be set by BeforeCreate hook")
	})
}

func TestOnboardingRequestFindByDocumentOrEmail(t *testing.T) {
	t.Run("#method FindByDocumentOrEmail", func(t *testing.T) {
		expectedFullName := "Jane Doe"
		expectedEmail := "jane.doe@example.com"
		expectedPublicId := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		expectedDocument := "12345678901"

		insertDB, err := testutil.ReadFileContent(filepath.Join(sqldir, "test", "add-onboarding-requests.sql"))
		require.NoError(t, err)

		db.Exec(insertDB)

		foundRequest, err := repo.FindByDocumentOrEmail(expectedDocument, "some-other-email@test.com")
		foundRequest2, err := repo.FindByDocumentOrEmail("1234567810", expectedEmail)

		require.NoError(t, err)
		assert.Equal(t, expectedEmail, foundRequest.Email)
		assert.Equal(t, expectedFullName, foundRequest.FullName)
		assert.Equal(t, expectedPublicId, foundRequest.PublicID)

		assert.Equal(t, expectedEmail, foundRequest2.Email)
		assert.Equal(t, expectedFullName, foundRequest2.FullName)
		assert.Equal(t, expectedPublicId, foundRequest2.PublicID)
	})
}
