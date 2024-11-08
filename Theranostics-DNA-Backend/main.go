// main.go
package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"

	"theransticslabs/m/config"
	"theransticslabs/m/models"
	"theransticslabs/m/routes"
	"theransticslabs/m/seeds"
	"theransticslabs/m/utils"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	allowedOrigins := strings.Split(config.AppConfig.AllowedOrigins, ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins, // Replace with your frontend URL(s)
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Initialize the Database
	config.InitDB()

	// Create enum types
	if err := config.CreateEnumTypes(); err != nil {
		log.Fatalf("Failed to create enum types: %v", err)
	}

	// Auto-Migrate the models
	err := config.DB.AutoMigrate(&models.Role{}, &models.User{}, &models.Kit{}, &models.Customer{}, &models.Order{}, &models.Payment{}, &models.Invoice{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println(utils.MsgDatabaseMigrated)

	// Run the Seeders
	seeds.SeedAll()

	// Initialize the Router and Routes
	router := routes.SetupRoutes()

	// Serve static files (images and invoices)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))
	router.PathPrefix("/invoices/").Handler(http.StripPrefix("/invoices/", http.FileServer(http.Dir("./public/invoices"))))

	// Wrap your router with CORS middleware
	handler := c.Handler(router)

	// Start the Server
	log.Printf(utils.MsgServerStarted, config.AppConfig.ServerPort)
	if err := http.ListenAndServe(":"+config.AppConfig.ServerPort, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
