package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"social_media/internal/domain"
	"social_media/internal/repository"
	"social_media/internal/service"
	"social_media/internal/handler"
	"social_media/pkg/jwt"
	"social_media/router"
)

func main() {
	// Load environment variables.
	if err := godotenv.Load("configs/config.env"); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	appPort := os.Getenv("APP_PORT")
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the User model. (In production, use migration tools.)
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize repository, services, handlers, and JWT manager.
	userRepo := repository.NewUserRepository(db)
	jwtManager := jwt.NewJWTManager(jwtSecret, time.Hour*24) // token valid for 24 hours
	authService := service.NewAuthService(userRepo, jwtManager)
	profileService := service.NewProfileService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)

	// Set up router.
	r := router.SetupRouter(authHandler, profileHandler, jwtManager)
	log.Printf("Server starting on port %s...", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
