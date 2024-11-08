// config/config.go
package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	LoadEnv()
	cfg := AppConfig

	// First, connect to the default 'postgres' database to check/create target database
	defaultDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)

	defaultDB, err := gorm.Open(postgres.Open(defaultDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress default logging
	})
	if err != nil {
		log.Fatalf("Failed to connect to default database: %v", err)
	}

	// Check if the target database exists
	var exists bool
	query := "SELECT 1 FROM pg_database WHERE datname = ?"
	err = defaultDB.Raw(query, cfg.DBName).Scan(&exists).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	if !exists {
		// Create the database
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		if err := defaultDB.Exec(createDBQuery).Error; err != nil {
			log.Fatalf("Failed to create database %s: %v", cfg.DBName, err)
		}
		log.Printf("Database %s created successfully", cfg.DBName)
	} else {
		log.Printf("Database %s already exists", cfg.DBName)
	}

	// Now, connect to the target database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("%s environment: %s", "Database connection established", cfg.Environment)
}

// CreateEnumTypes creates the necessary enum types in the database
func CreateEnumTypes() error {
	enumTypes := map[string]string{
		"dna_payment_status": "CREATE TYPE payment_status AS ENUM ('Pending', 'Completed', 'Failed')",
		"dna_order_status":   "CREATE TYPE order_status AS ENUM ('Pending', 'Shipped', 'Delivered', 'Cancelled')",
	}

	for typeName, query := range enumTypes {
		var exists bool
		checkQuery := fmt.Sprintf("SELECT 1 FROM pg_type WHERE typname = '%s'", typeName)

		// Check if the enum type already exists
		if err := DB.Raw(checkQuery).Scan(&exists).Error; err != nil {
			return fmt.Errorf("failed to check if enum type %s exists: %w", typeName, err)
		}

		if exists {
			log.Printf("Enum type %s already exists", typeName)
			continue
		}

		// Create the enum type if it does not exist
		if err := DB.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to create enum type %s: %w", typeName, err)
		}
		log.Printf("Enum type %s created successfully", typeName)
	}
	return nil
}
