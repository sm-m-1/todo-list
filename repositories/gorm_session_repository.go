package repositories

import (
	"errors"
	"fmt"
	"time"
	"todo-list/models"

	"gorm.io/gorm"
)

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

	// Use GORM's `Save` method to handle both insert and update logic
	if err := s.db.Save(&session).Error; err != nil {
		// Return error for any unexpected system-level issue
		return err
	}
	fmt.Println("saved token in commit key: ", key)

	return nil
}

// Find retrieves the session data by its key.
func (s *GORMStore) Find(key string) ([]byte, bool, error) {
	fmt.Println("trying to find token for the key: ", key)
	var session models.Session
	err := s.db.Where("token = ?", key).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Session token not found
			return nil, false, nil
		}
		// System-level error (e.g., database connection issue)
		return nil, false, err
	}

	// Check if the session has expired
	if time.Now().After(session.Expiry) {
		// Session expired, delete it and return false for found
		_ = s.db.Delete(&session).Error // Cleanup, but don't block on error
		return nil, false, nil
	}

	// Return the session data and indicate it was found
	return session.Data, true, nil
}

// Delete removes the session data from the database.
func (s *GORMStore) Delete(key string) error {
	// Attempt to find and delete the session token
	err := s.db.Where("token = ?", key).Delete(&models.Session{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Return error only for unexpected system-level issues
		return err
	}

	// If the token does not exist, it's a no-op, so return nil
	return nil
}
