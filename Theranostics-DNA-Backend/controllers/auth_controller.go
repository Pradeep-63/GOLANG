// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"theransticslabs/m/config"
	"theransticslabs/m/middlewares"
	"theransticslabs/m/utils"

	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest represents the expected payload for login.
type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// LoginResponse represents the response after a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles user login requests.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Define allowed fields for this request
	allowedFields := []string{"email", "password"}

	// Use the common request parser for both JSON and form data, and validate allowed fields
	err := utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}
	// Validate input
	if req.Email == "" || req.Password == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgMissingEmailOrPassword, nil)
		return
	}

	// Trim spaces and convert email to lowercase
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if !utils.IsValidEmail(req.Email) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidEmailFormat, nil)
		return
	}

	if !utils.IsValidPassword(req.Password) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidPasswordFormat, nil)
		return
	}

	// Find the user by email using the common function
	user, err := utils.FindUserByEmail(config.DB, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgInvalidCredentials, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password))
	if err != nil {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgInvalidCredentials, nil)
		return
	}

	// Check if user is active and not deleted
	if !user.ActiveStatus || user.IsDeleted {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserInactive, nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(*user)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgTokenCreationFailed, nil)
		return
	}

	// Save the token to the user's Token field
	user.Token = token
	err = config.DB.Save(&user).Error
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgTokenSaveFailed, nil)
		return
	}

	// Prepare the response
	response := LoginResponse{
		Token: token,
	}

	// Send the response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgLoginSuccess, response)
}

// LogoutHandler removes the JWT token from the user's record in the database.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Verify the token (this is done by the AuthMiddleware)

	// 2. Get the user from the context (set by AuthMiddleware)
	user, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserNotFound, nil)
		return
	}

	// Begin a transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear the token in the database
	user.Token = "" // Set the token field to an empty string
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

	// Send a success response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgLogoutSuccess, nil)
}
