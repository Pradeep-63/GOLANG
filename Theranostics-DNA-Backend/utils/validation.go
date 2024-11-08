// utils/validation.go
package utils

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// IsValidEmail checks if the provided email has a valid format.
func IsValidEmail(email string) bool {
	// Simple regex for email validation
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email) && len(email) <= 255 // Check email length
}

// IsValidFirstName checks if the first name is valid: alphabetic and length 3-50.
func IsValidFirstName(firstName string) bool {
	return isAlphabetic(firstName) && len(firstName) >= 3 && len(firstName) <= 50
}

// IsValidLastName checks if the last name is valid: alphabetic and length 0-50.
func IsValidLastName(lastName string) bool {
	if lastName == "" {
		return true // Optional last name
	}
	return isAlphabetic(lastName) && len(lastName) >= 3 && len(lastName) <= 50
}

// isAlphabetic checks if the string contains only alphabetic characters or single spaces between words
func isAlphabetic(s string) bool {
	spaceAllowed := true // To keep track if space is allowed (only one space between words)

	for i, char := range s {
		if unicode.IsLetter(char) {
			spaceAllowed = true // Reset the flag when a letter is found
		} else if char == ' ' {
			// If a space is found, check if it's allowed (only one space between words)
			if !spaceAllowed || (i > 0 && s[i-1] == ' ') {
				return false // Consecutive spaces or space not allowed
			}
			spaceAllowed = false // Space found, disallow another consecutive space
		} else {
			return false // Any other character is invalid
		}
	}
	return true
}

// IsValidPassword checks if the password meets the criteria:
// 8-20 characters, at least one uppercase, one lowercase, one number, and one special character.
func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case isSpecialChar(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// isSpecialChar checks if the character is a special character.
func isSpecialChar(char rune) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?`~"
	for _, sc := range specialChars {
		if char == sc {
			return true
		}
	}
	return false
}

// IsValidKitType checks if the kit type is valid
func IsValidKitType(kitType string) bool {
	return kitType == "blood" || kitType == "saliva"
}

// IsValidSupplierName checks if the supplier name is alphabetic and within the required length
func IsValidSupplierName(name string) bool {
	if len(name) < 3 || len(name) > 50 {
		return false
	}
	return isAlphabetic(name)
}

// IsValidContactNumber checks if the contact number is numeric and between 10-15 digits
func IsValidContactNumber(number string) bool {
	if len(number) < 10 || len(number) > 15 {
		return false
	}
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(number)
}

// IsValidProductName checks if the product name meets the requirements
func IsValidProductName(name string) bool {
	if len(name) < 3 || len(name) > 100 {
		return false
	}
	// Check for valid characters (letters, numbers, spaces, and basic punctuation)
	validChar := regexp.MustCompile(`^[a-zA-Z0-9\s\-_,.&()]+$`)
	return validChar.MatchString(name)
}

// IsValidBase64Image checks if a base64-encoded string is a valid image format.
func IsValidBase64Image(str string) bool {
	if str == "" {
		return true // Empty string is valid for optional base64
	}

	// Check if base64 string starts with a valid image MIME type
	if !strings.HasPrefix(str, "data:image/") {
		return false
	}

	// Find the comma separator and decode the actual base64 data
	base64DataIdx := strings.Index(str, ",")
	if base64DataIdx == -1 {
		return false
	}

	// Decode base64 data
	_, err := base64.StdEncoding.DecodeString(str[base64DataIdx+1:])
	return err == nil
}

// IsValidImageURL checks if the provided string is a valid image URL format.
func IsValidImageURL(url string) bool {
	// Basic URL validation with regex
	validURL := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	if !validURL.MatchString(url) {
		return false
	}

	// Check if URL ends with a common image file extension
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	for _, ext := range imageExtensions {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}

// IsValidPrice checks if the price string is a valid positive number
func IsValidPrice(price string) bool {
	// First check if it's a valid number
	val, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return false
	}

	// Check if it's positive and has maximum 2 decimal places
	if val <= 0 {
		return false
	}

	// Check decimal places
	parts := strings.Split(price, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return false
	}

	return true
}

// IsValidQuantity checks if the quantity string is a valid positive integer
func IsValidQuantity(quantity string) bool {
	val, err := strconv.Atoi(quantity)
	if err != nil {
		return false
	}
	return val > 0
}
