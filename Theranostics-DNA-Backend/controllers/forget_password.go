// controllers/forget_password.go
package controllers

import (
	"net/http"
	"strings"

	"gorm.io/gorm"

	"theransticslabs/m/config"
	"theransticslabs/m/emails"
	"theransticslabs/m/utils"
)

// PasswordForgetRequest represents the expected request body
type PasswordForgetRequest struct {
	Email string `json:"email" form:"email"`
}

func ForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req PasswordForgetRequest
	// Define allowed fields for this request
	allowedFields := []string{"email"}

	// Use the common request parser for both JSON and form data, and validate allowed fields
	err := utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Validate input
	if req.Email == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgMissingEmail, nil)
		return
	}

	// Trim spaces and convert email to lowercase
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidEmailFormat, nil)
		return
	}

	// Find the user by email using the common function
	user, err := utils.FindUserByEmail(config.DB, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgEmailNotExist, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Generate a secure random password
	newPassword := utils.GenerateSecurePassword()

	// Start a transaction
	tx := config.DB.Begin()

	// Defer a rollback in case of failure
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Send the email
	emailBody := emails.ResetPasswordEmail(user.FirstName, user.LastName, req.Email, newPassword, config.AppConfig.AppUrl)
	if err := config.SendEmail([]string{req.Email}, "Password Reset", emailBody); err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedSentEmail, nil)
		return

	}

	// Hash and update the new password in the database
	hashedPassword, err := utils.HashPassword(newPassword) // Implement HashPassword function
	if err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return

	}

	// Update user's hashed password and remove the token
	user.HashPassword = hashedPassword
	user.Token = ""
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Send the response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgForgetSuccessfully, nil)
}
