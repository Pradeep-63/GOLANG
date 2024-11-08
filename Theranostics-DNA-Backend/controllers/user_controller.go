// controllers/user_controller.go
package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"theransticslabs/m/config"
	"theransticslabs/m/emails"
	"theransticslabs/m/middlewares"
	"theransticslabs/m/models"
	"theransticslabs/m/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UpdateUserRequest struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
}

// UserProfile represents the user data to be sent in the response.
type UserProfile struct {
	ID           uint        `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Email        string      `json:"email"`
	Role         RoleProfile `json:"role"`
	ActiveStatus bool        `json:"active_status"`
	CreatedAt    time.Time   `json:"created_at"`
}

// RoleProfile represents the role data to be sent in the response.
type RoleProfile struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// GetUserProfileResponse represents the response structure.
type GetUserProfileResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    UserProfile `json:"data"`
}

// AdminUsersResponse represents the structured response for admin users list.
type AdminUsersResponse struct {
	Page         int           `json:"page"`
	PerPage      int           `json:"per_page"`
	Sort         string        `json:"sort"`
	SortColumn   string        `json:"sort_column"`
	SearchText   string        `json:"search_text"`
	Status       string        `json:"status"`
	TotalRecords int64         `json:"total_records"`
	TotalPages   int           `json:"total_pages"`
	Records      []UserProfile `json:"records"`
}

// CreateUserRequest represents the expected input for creating a user.
type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,max=50" form:"first_name"`
	LastName  string `json:"last_name" validate:"required,max=50" form:"last_name"`
	Email     string `json:"email" validate:"required,email,max=100" form:"email"`
}

type UpdateAdminUserRequest struct {
	FirstName    string `json:"first_name,omitempty" form:"first_name"`
	LastName     string `json:"last_name,omitempty" form:"last_name"`
	Email        string `json:"email,omitempty" form:"email"`
	Password     string `json:"password,omitempty" form:"password"`
	ActiveStatus bool   `json:"active_status,omitempty" form:"active_status"`
}

// Request structs for different update operations
type UpdateUserProfileRequest struct {
	FirstName string  `json:"first_name" form:"first_name"`
	LastName  *string `json:"last_name" form:"last_name"`
	Email     string  `json:"email" form:"email"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password" form:"password"`
}

type UpdateUserStatusRequest struct {
	ActiveStatus bool `json:"active_status" form:"active_status"`
}

func UpdateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. This is a private route (enforced by AuthMiddleware)

	// Get the user from the context (set by AuthMiddleware)
	user, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserNotFound, nil)
		return
	}

	// Parse the request body
	var req UpdateUserRequest
	// Define allowed fields for this request
	allowedFields := []string{"first_name", "last_name"}

	// Use the common request parser for both JSON and form data, and validate allowed fields
	err := utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Trim the first name
	req.FirstName = strings.TrimSpace(req.FirstName)

	// 5. Check if at least the first name is provided
	if req.FirstName == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgAtleastPassSomeData, nil)
		return
	}

	// Validate first name
	if !utils.IsValidFirstName(req.FirstName) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgFirstNameValidation, nil)
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Update first name
	user.FirstName = req.FirstName

	// Handle last name
	if req.LastName != "" {
		req.LastName = strings.TrimSpace(req.LastName)
		if !utils.IsValidLastName(req.LastName) {
			utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgLastNameValidation, nil)
			return
		}
		user.LastName = req.LastName
	} else {
		// Set last name to nil (or zero value)
		user.LastName = "" // Assuming the LastName field is a string
	}

	// 3. Update the user details
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

	// Prepare the response user profile
	userProfile := UserProfile{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role: RoleProfile{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgUserUpdateSuccessfully, userProfile)
}

// GetUserProfileHandler handles requests to fetch user profile details.
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context (set by AuthMiddleware)
	user, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserNotFound, nil)
		return
	}

	// Prepare the user profile
	userProfile := UserProfile{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role: RoleProfile{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
	}

	// Send the encrypted response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgUserFetchedSuccessfully, userProfile)
}

// GetAdminUsersHandler handles requests to fetch the admin users list.
func GetAdminUsersHandler(w http.ResponseWriter, r *http.Request) {

	// Define allowed query parameters
	allowedFields := []string{"page", "per_page", "sort", "sort_column", "search_text", "status"}

	// Parse query parameters with default values
	query := r.URL.Query()
	if !utils.AllowFields(query, allowedFields) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidQueryParameters, nil)
		return
	}

	// Default and validation for 'page'
	page := 1
	if val := query.Get("page"); val != "" {
		if p, err := strconv.Atoi(val); err == nil && p > 0 {
			page = p
		} else {
			utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidPageParameter, nil)
			return
		}
	}

	// Default and validation for 'per_page'
	perPage := 10
	if val := query.Get("per_page"); val != "" {
		if pp, err := strconv.Atoi(val); err == nil && pp > 0 {
			perPage = pp
		} else {
			utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidPerPageParameter, nil)
			return
		}
	}

	// Default and validation for 'sort'
	sort := "desc"
	if val := strings.ToLower(query.Get("sort")); val == "asc" || val == "desc" {
		sort = val
	} else if val != "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidSortParameter, nil)
		return
	}

	// Default and validation for 'sort_column'
	sortColumn := "created_at"
	validSortColumns := []string{"first_name", "last_name", "email", "active_status", "created_at"}
	if val := strings.ToLower(query.Get("sort_column")); val != "" {
		if utils.StringInSlice(val, validSortColumns) {
			sortColumn = val
		} else {
			utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidSortColumnParameter, nil)
			return
		}
	}

	// Optional 'search_text'
	searchText := strings.TrimSpace(query.Get("search_text"))

	// Default and validation for 'status'
	status := "all"
	if val := strings.ToLower(query.Get("status")); val == "active" || val == "inactive" || val == "all" {
		status = val
	} else if val != "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgStatusInvalid, nil)
		return
	}

	// Initialize GORM query with Debug for detailed logging
	db := config.DB.Debug().
		Model(&models.User{}).
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "admin").
		Where("users.is_deleted = ?", false)

	// Apply status filter
	if status == "active" {
		db = db.Where("users.active_status = ?", true)
	} else if status == "inactive" {
		db = db.Where("users.active_status = ?", false)
	}

	// Apply search filter
	if searchText != "" {
		searchPattern := "%" + searchText + "%"
		db = db.Where("roles.name = ?", "admin"). // Ensure only admin users are selected
								Where(
				"users.first_name ILIKE ? OR users.last_name ILIKE ? OR users.email ILIKE ? OR CAST(users.created_at AS TEXT) ILIKE ?",
				searchPattern, searchPattern, searchPattern, searchPattern,
			)
	}

	// Get total records count
	var totalRecords int64
	if err := db.Count(&totalRecords).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToCountRecords, nil)
		return
	}

	// Calculate total pages
	var totalPages int
	if totalRecords == 0 {
		totalPages = 0
	} else {
		totalPages = int((totalRecords + int64(perPage) - 1) / int64(perPage))
	}

	// Apply sorting
	db = db.Order(sortColumn + " " + sort)

	// Apply pagination
	offset := (page - 1) * perPage
	db = db.Limit(perPage).Offset(offset)

	// Fetch records
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToFetchRecords, nil)
		return
	}

	// Prepare user profiles
	var userProfiles []UserProfile
	for _, u := range users {
		profile := UserProfile{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Role: RoleProfile{
				ID:   u.Role.ID,
				Name: u.Role.Name,
			},
			ActiveStatus: u.ActiveStatus,
			CreatedAt:    u.CreatedAt,
		}
		userProfiles = append(userProfiles, profile)
	}

	// Prepare the response
	response := AdminUsersResponse{
		Page:         page,
		PerPage:      perPage,
		Sort:         sort,
		SortColumn:   sortColumn,
		SearchText:   searchText,
		Status:       status,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		Records:      userProfiles,
	}

	// Send the response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgUserListFetchedSuccessfully, response)
}

// CreateUserHandler handles the creation of a new admin user.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse and Validate Request Body
	var req CreateUserRequest

	// Define allowed fields for this request
	allowedFields := []string{"email", "first_name", "last_name"}

	// Use the common request parser for both JSON and form data, and validate allowed fields
	err := utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// 2. Validate Input Data
	if req.Email == "" || req.FirstName == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgFirstNameEmailRequired, nil)
		return
	}

	// Trim
	req.FirstName = strings.TrimSpace(req.FirstName)
	if req.FirstName == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgAtleastPassSomeData, nil)
		return
	}

	// Validate first name
	if !utils.IsValidFirstName(req.FirstName) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgFirstNameValidation, nil)
		return
	}

	if req.LastName != "" {
		req.LastName = strings.TrimSpace(req.LastName)
		if !utils.IsValidLastName(req.LastName) {
			utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgLastNameValidation, nil)
			return
		}

	} else {
		// Set last name to nil (or zero value)
		req.LastName = "" // Assuming the LastName field is a string
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if !utils.IsValidEmail(req.Email) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgEmailValidation, nil)
		return
	}

	// 3. Check if the Email Already Exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.JSONResponse(w, http.StatusConflict, false, utils.MsgEmailAlreadyInUse, nil)
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// 4. Retrieve Admin Role ID
	var adminRole models.Role
	result := config.DB.Where("name = ? AND is_deleted = ?", "admin", false).First(&adminRole)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgAdminRoleNotFound, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// 5. Generate and Hash Password
	password := utils.GenerateSecurePassword()
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToSecurePassword, nil)
		return
	}

	// 6. Create a Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 7. Create the User
	newUser := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		HashPassword: hashedPassword,
		RoleID:       adminRole.ID, // Use adminRole.ID directly
		ActiveStatus: true,
		IsDeleted:    false,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedCreateUser, nil)
		return
	}

	// 8. Send Email to the User
	appUrl := config.AppConfig.AppUrl // Ensure this is set in your configuration
	emailBody := emails.WelcomeEmail(newUser.FirstName, newUser.LastName, newUser.Email, password, appUrl)
	if err := config.SendEmail([]string{newUser.Email}, "Welcome to Our Platform", emailBody); err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedSentEmail, nil)
		return
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedCreateUser, nil)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, true, utils.MsgUserCreatedSuccessfully, nil)
}

// DeleteUserHandler marks the user as deleted (sets is_deleted to true) based on the user ID.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	// Extract user ID from URL
	userID := mux.Vars(r)["id"]

	// Find the user based on ID and check if it's already deleted
	var existingUser models.User
	if err := config.DB.Where("id = ? AND is_deleted = ?", userID, false).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgUserAlreadyDeleted, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Mark user as deleted (set is_deleted to true)
	existingUser.IsDeleted = true

	// Create a transaction for updating the user
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(&existingUser).Error; err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToDeleteUser, nil)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgUserDeleteSuccessfully, nil)
}

// Common function to fetch existing user
func getExistingUser(userID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("id = ? AND is_deleted = ?", userID, false).First(&user).Error
	return &user, err
}

// Helper functions for update operations
func updateUserProfile(user *models.User, req *UpdateUserProfileRequest) error {
	if req.FirstName != "" {
		if !utils.IsValidFirstName(req.FirstName) {
			return errors.New(utils.MsgFirstNameValidation)
		}
		user.FirstName = req.FirstName
	}

	if req.LastName != nil { // Only update if LastName is provided
		if !utils.IsValidLastName(*req.LastName) {
			return errors.New(utils.MsgLastNameValidation)
		}
		user.LastName = *req.LastName
	}
	// Check if Email is valid and update
	if req.Email != "" {
		if !utils.IsValidEmail(req.Email) {
			return errors.New(utils.MsgEmailValidation)
		}

		// Check for existing email in the database
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", req.Email, user.ID).First(&existingUser).Error; err == nil {
			return errors.New(utils.MsgEmailAlreadyInUse)
		}
		user.Email = req.Email
	}

	return nil
}

func saveUserAndNotify(tx *gorm.DB, user *models.User, req *UpdateUserProfileRequest) error {
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return errors.New(utils.MsgFailedToUpdateAdminUser)
	}

	emailBody := emails.UserDetailsUpdatedEmail(req.FirstName, *req.LastName, req.Email, config.AppConfig.AppUrl)
	if err := config.SendEmail([]string{user.Email}, "Your Profile Details Updated", emailBody); err != nil {
		tx.Rollback()
		return errors.New("Failed to send profile update notification")
	}

	return tx.Commit().Error
}

// UpdateUserProfileHandler handles updating user's profile information
func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {

	var req UpdateUserProfileRequest
	if err := utils.ParseRequestBody(r, &req, []string{"first_name", "last_name", "email"}); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	existingUser, err := getExistingUser(mux.Vars(r)["id"])
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgUserNotFound, nil)
		return
	}

	// Validate and update profile fields
	if err := updateUserProfile(existingUser, &req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	tx := config.DB.Begin()
	if err := saveUserAndNotify(tx, existingUser, &req); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgAdminUserUpdateSuccessfully, nil)
}

func updateUserPassword(user *models.User, password string) error {
	if !utils.IsValidPassword(password) {
		return errors.New(utils.MsgPasswordValidation)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New(utils.MsgFailedToSecurePassword)
	}
	user.HashPassword = hashedPassword
	user.Token = ""

	return nil
}

func saveUserAndSendPasswordNotification(tx *gorm.DB, user *models.User, password string) error {
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return errors.New(utils.MsgFailedToUpdateAdminUser)
	}

	emailBody := emails.PasswordUpdatedEmail(user.FirstName, user.LastName, user.Email, password, config.AppConfig.AppUrl)
	if err := config.SendEmail([]string{user.Email}, "Your Password has been Updated", emailBody); err != nil {
		tx.Rollback()
		return errors.New("Failed to send password update notification")
	}

	return tx.Commit().Error
}

// UpdateUserPasswordHandler handles password updates
func UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {

	var req UpdateUserPasswordRequest
	if err := utils.ParseRequestBody(r, &req, []string{"password"}); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	existingUser, err := getExistingUser(mux.Vars(r)["id"])
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgUserNotFound, nil)
		return
	}

	// Validate and update password
	if err := updateUserPassword(existingUser, req.Password); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	tx := config.DB.Begin()
	if err := saveUserAndSendPasswordNotification(tx, existingUser, req.Password); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgAdminUserUpdateSuccessfully, nil)
}

func updateUserStatus(tx *gorm.DB, user *models.User, newStatus bool) error {
	if user.ActiveStatus == newStatus {
		return errors.New("Status is already set to the requested value")
	}

	user.ActiveStatus = newStatus
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return errors.New(utils.MsgFailedToUpdateAdminUser)
	}

	// Send status change notification
	statusMessage := "activated"
	if !newStatus {
		statusMessage = "deactivated"
	}
	emailBody := emails.UserStatusChangedEmail(user.FirstName, user.LastName, statusMessage, config.AppConfig.AppUrl)
	if err := config.SendEmail([]string{user.Email}, "Your Account Status has Changed", emailBody); err != nil {
		tx.Rollback()
		return errors.New("Failed to send status change notification")
	}

	return tx.Commit().Error
}

// UpdateUserStatusHandler handles user status updates
func UpdateUserStatusHandler(w http.ResponseWriter, r *http.Request) {

	var req UpdateUserStatusRequest
	if err := utils.ParseRequestBody(r, &req, []string{"active_status"}); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	existingUser, err := getExistingUser(mux.Vars(r)["id"])
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgUserNotFound, nil)
		return
	}

	tx := config.DB.Begin()
	if err := updateUserStatus(tx, existingUser, req.ActiveStatus); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgAdminUserUpdateSuccessfully, nil)
}
