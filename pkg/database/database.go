package database

import (
	"fmt"
	"log"
	"todo-list/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := "host=db user=postgres password=yourpassword dbname=todo_list port=5432 sslmode=disable"

	// Try to connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// If the database doesn't exist, create it
		if err.Error() == `FATAL: database "todo_list" does not exist (SQLSTATE 3D000)` {
			createDatabase()
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("Failed to reconnect to the database: %v", err)
			}
		} else {
			log.Fatalf("Failed to connect to the database: %v", err)
		}
	}

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Todo{},
		&models.Session{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB instance: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	} else {
		log.Println("Database connection closed successfully.")
	}
}

// Creates the "todo_list" database if it does not exist
func createDatabase() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL to create database: %v", err)
	}

	sql := "CREATE DATABASE todo_list;"
	if err := db.Exec(sql).Error; err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	fmt.Println("Database todo_list created successfully.")
}
