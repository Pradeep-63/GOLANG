// utils/common.go
package utils

import (
	"theransticslabs/m/models"

	"gorm.io/gorm"
)

// FindUserByEmail searches for a user by email.
// Returns the user and any error encountered.
func FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	// Corrected query with two conditions
	result := db.Preload("Role").Where("email = ? AND is_deleted = ?", email, false).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// StringInSlice checks if a string exists in a slice of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
