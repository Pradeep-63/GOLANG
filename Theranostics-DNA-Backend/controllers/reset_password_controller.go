package controllers

import (
	"net/http"
	"strings"

	"theransticslabs/m/config"
	"theransticslabs/m/middlewares"
	"theransticslabs/m/utils"

	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Verify the token (this is done by the AuthMiddleware)

	// 2. Get the user from the context (set by AuthMiddleware)
	user, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserNotFound, nil)
		return
	}

	// 3. Parse the request body
	var req ResetPasswordRequest
	// Define allowed fields for this request
	allowedFields := []string{"old_password", "new_password"}

	// Use the common request parser for both JSON and form data, and validate allowed fields
	err := utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Trim whitespace
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.OldPassword = strings.TrimSpace(req.OldPassword)

	// Validate input
	if req.NewPassword == "" || req.OldPassword == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgOldNewPasswordRequired, nil)
		return
	}

	// 4. Validate the new password
	if !utils.IsValidPassword(req.NewPassword) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidPasswordFormat, nil)
		return
	}

	// 5. Check if old password and new password are the same
	if req.OldPassword == req.NewPassword {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgPasswordSame, nil)
		return
	}

	// 6. Verify the old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.OldPassword)); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidOldPassword, nil)
		return
	}

	// 7. Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToHashNewPassword, nil)
		return
	}

	// Update user's hashed password and remove the token
	user.HashPassword = hashedPassword
	user.Token = ""
	if err := config.DB.Save(user).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToUpdateNewPassword, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgResetPasswordSuccessfully, nil)
}
