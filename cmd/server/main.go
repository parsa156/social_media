package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v4/pgxpool"

	"social_media/internal/repository"
	"social_media/internal/service"
	"social_media/internal/handler"
	"social_media/pkg/jwt"
	"social_media/router"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(".env"); err != nil {
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

	// Build the PostgreSQL connection string for pgx.
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Database connection established using pgx.")

	// Initialize repositories using the pgx pool.
	userRepo := repository.NewUserRepository(pool)
	convoRepo := repository.NewConversationRepository(pool)
	messageRepo := repository.NewMessageRepository(pool)
	roomRepo := repository.NewRoomRepository(pool)
	roomMembershipRepo := repository.NewRoomMembershipRepository(pool)
	roomMessageRepo := repository.NewRoomMessageRepository(pool)

	// Initialize the JWT Manager.
	jwtManager := jwt.NewJWTManager(jwtSecret, time.Hour*24) // Token valid for 24 hours.

	// Initialize services.
	authService := service.NewAuthService(userRepo, jwtManager)
	profileService := service.NewProfileService(userRepo)
	convoService := service.NewConversationService(convoRepo, messageRepo, userRepo)
	roomService := service.NewRoomService(roomRepo, roomMembershipRepo, roomMessageRepo)

	// Initialize handlers.
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	convoHandler := handler.NewConversationHandler(convoService)
	roomHandler := handler.NewRoomHandler(roomService)

	// Setup the router with public and protected endpoints.
	r := router.SetupRouter(authHandler, profileHandler, convoHandler, roomHandler, jwtManager)

	// Start the server.
	log.Printf("Server starting on port %s...", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
