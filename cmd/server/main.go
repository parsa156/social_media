package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"social_media/internal/repository"
	"social_media/internal/service"
	"social_media/internal/handler"
	"social_media/pkg/jwt"
	"social_media/router"
)

func main() {
	// Load environment variables from configs/config.env
	if err := godotenv.Load("configs/config.env"); err != nil {
		log.Printf("Warning: no config file found, using environment variables")
	}

	// Retrieve configuration values.
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Build the PostgreSQL DSN.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// Initialize the database connection.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// NOTE: Do not use AutoMigrate in production.
	// Migrations are assumed to be handled externally via migration files.
	log.Println("Database connection established.")

	// Initialize repositories.
	userRepo := repository.NewUserRepository(db)
	convoRepo := repository.NewConversationRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Initialize the JWT Manager.
	jwtManager := jwt.NewJWTManager(jwtSecret, time.Hour*24) // Token valid for 24 hours.

	// Initialize services.
	authService := service.NewAuthService(userRepo, jwtManager)
	profileService := service.NewProfileService(userRepo)
	// Pass userRepo to conversation service so it can look up recipients by phone or username.
	convoService := service.NewConversationService(convoRepo, messageRepo, userRepo)

	// Initialize handlers.
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	convoHandler := handler.NewConversationHandler(convoService)

	// Setup the router with public and protected endpoints.
	r := router.SetupRouter(authHandler, profileHandler, convoHandler, jwtManager)

	// Start the server.
	log.Printf("Server starting on port %s...", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
