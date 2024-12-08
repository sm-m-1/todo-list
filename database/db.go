package database

import (
	"fmt"
	"log"
	"time"
	"todo-list/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=yourpassword dbname=todo_list port=5432 sslmode=disable"

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
		&models.Todo{},
		&models.User{},
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

// GORMStore implements the scs.Store interface
type GORMStore struct {
	db       *gorm.DB
	lifetime time.Duration
}

// NewGORMStore creates a new GORM-based session store
func NewGORMStore(db *gorm.DB, lifetime time.Duration) *GORMStore {
	return &GORMStore{
		db:       db,
		lifetime: lifetime,
	}
}

// Commit stores the session data in the database.
func (s *GORMStore) Commit(key string, data []byte, expiry time.Time) error {
	session := models.Session{
		Token:  key,
		Data:   data,
		Expiry: expiry,
	}

	// Upsert the session data (insert if new, update if existing)
	return s.db.Save(&session).Error
}

// Fetch retrieves the session data from the database.
func (s *GORMStore) Fetch(key string) ([]byte, error) {
	var session models.Session
	if err := s.db.Where("token = ?", key).First(&session).Error; err != nil {
		return nil, err
	}

	if time.Now().After(session.Expiry) {
		// Session has expired, delete it and return nil
		s.db.Delete(&session)
		return nil, nil
	}

	return session.Data, nil
}

// Find retrieves the session data by its key.
func (s *GORMStore) Find(key string) ([]byte, bool, error) {
	var session models.Session
	if err := s.db.Where("token = ?", key).First(&session).Error; err != nil {
		return nil, false, err
	}

	if time.Now().After(session.Expiry) {
		// Session has expired, delete it and return nil
		s.db.Delete(&session)
		return nil, false, nil
	}

	return session.Data, true, nil
}

// Delete removes the session data from the database.
func (s *GORMStore) Delete(key string) error {
	return s.db.Where("token = ?", key).Delete(&models.Session{}).Error
}
