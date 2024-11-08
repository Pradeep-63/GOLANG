// manage_inventory_controller
package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"theransticslabs/m/config"
	"theransticslabs/m/middlewares"
	"theransticslabs/m/models"
	"theransticslabs/m/utils"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// KitRequest represents the expected request body structure
type KitRequest struct {
	Type                  string      `json:"type" form:"type"`
	SupplierName          string      `json:"supplier_name" form:"supplier_name"`
	SupplierContactNumber string      `json:"supplier_contact_number" form:"supplier_contact_number"`
	SupplierAddress       string      `json:"supplier_address" form:"supplier_address"`
	Quantity              interface{} `json:"quantity" form:"quantity"`
}

type KitsListResponse struct {
	Page         int         `json:"page"`
	PerPage      int         `json:"per_page"`
	Sort         string      `json:"sort"`
	SortColumn   string      `json:"sort_column"`
	SearchText   string      `json:"search_text"`
	TotalRecords int64       `json:"total_records"`
	TotalPages   int         `json:"total_pages"`
	Type         string      `json:"type"`
	Records      []KitDetail `json:"records"`
}

type KitDetail struct {
	ID                    uint            `json:"id"`
	Type                  string          `json:"type"`
	Quantity              int             `json:"quantity"`
	SupplierName          string          `json:"supplier_name"`
	SupplierAddress       string          `json:"supplier_address"`
	SupplierContactNumber string          `json:"supplier_contact_number"`
	Status                bool            `json:"status"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedBy             KitsUserProfile `json:"created_by"`
}

type KitsUserProfile struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// KitUpdateRequest represents the PATCH request structure
type KitUpdateRequest struct {
	Type                  *string     `json:"type" form:"type"`
	SupplierName          *string     `json:"supplier_name" form:"supplier_name"`
	SupplierContactNumber *string     `json:"supplier_contact_number" form:"supplier_contact_number"`
	SupplierAddress       *string     `json:"supplier_address" form:"supplier_address"`
	Quantity              interface{} `json:"quantity" form:"quantity"`
	Status                *bool       `json:"status" form:"status"`
}

// CreateKitHandler handles the creation of a new kit
func CreateKitHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the user from the context (set by AuthMiddleware)
	user, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserNotAuthenticated, nil)
		return
	}

	// Parse the request body
	var req KitRequest
	// Define allowed fields for this request
	allowedFields := []string{"type", "supplier_name", "supplier_contact_number", "supplier_address", "quantity"}

	// Use the common request parser for both JSON and form data, and validate allowed fields
	err := utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Trim any spaces
	req.Type = strings.TrimSpace(req.Type)
	req.SupplierName = strings.TrimSpace(req.SupplierName)
	req.SupplierContactNumber = strings.TrimSpace(req.SupplierContactNumber)
	req.SupplierAddress = strings.TrimSpace(req.SupplierAddress)

	// 2. Validate required fields and their constraints
	if err := validateKitRequest(req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Parse quantity
	quantity, err := parseQuantity(req.Quantity)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// 5. Create the Kit model
	newKit := models.Kit{
		Type:     req.Type,
		Quantity: quantity,
		ExtraInfo: models.ExtraInfo{
			SupplierName:          req.SupplierName,
			SupplierContactNumber: req.SupplierContactNumber,
			SupplierAddress:       req.SupplierAddress,
		},
		CreatedBy: user.ID, // Extracted from the token
		Status:    true,
		IsDeleted: false,
	}

	// 6. Store kit in the database
	if err := config.DB.Create(&newKit).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgKitCreatedSuccessfully, nil)
}

// Helper function to convert quantity to int
func parseQuantity(value interface{}) (int, error) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return 0, errors.New(utils.MsgQuantityCannotBeEmpty)
		}
		quantity, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New(utils.MsgInvalidQuantityFormat)
		}
		return quantity, nil
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, errors.New(utils.MsgInvalidQuantityType)
	}
}

// validateKitRequest performs validation on all fields
func validateKitRequest(req KitRequest) error {
	if req.Type == "" || req.SupplierName == "" {
		return errors.New(utils.MsgAllFieldsOfkitRequired)
	}

	// Validate Kit Type
	if !utils.IsValidKitType(req.Type) {
		return errors.New(utils.MsgInvalidKitType)
	}

	// Validate Supplier Name
	if !utils.IsValidSupplierName(req.SupplierName) {
		return errors.New(utils.MsgValidationSupplierName)
	}

	// Validate Contact Number only if provided
	if req.SupplierContactNumber != "" && !utils.IsValidContactNumber(req.SupplierContactNumber) {
		return errors.New(utils.MsgValidationContactNumber)
	}

	// Validate Address only if provided
	if req.SupplierAddress != "" && (len(req.SupplierAddress) < 5 || len(req.SupplierAddress) > 100) {
		return errors.New(utils.MsgValidationAddress)
	}

	// Validate Quantity
	quantity, err := parseQuantity(req.Quantity)
	if err != nil {
		return err
	}
	if quantity < 0 {
		return errors.New(utils.MsgQuantityCannotBeNegative)
	}
	if quantity > 999999 { // Add a reasonable upper limit
		return errors.New(utils.MsgQuantityExceedsMaxValue)
	}

	return nil
}

// GetKitsListHandler handles requests to fetch the kits list.
func GetKitsListHandler(w http.ResponseWriter, r *http.Request) {

	// Define allowed query parameters
	allowedFields := []string{"page", "per_page", "sort", "sort_column", "search_text", "status", "type"}

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
	validSortColumns := []string{"type", "supplier_name", "supplier_address", "supplier_contact_number", "created_by", "created_at", "quantity"}
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

	// Optional 'type' with validation
	kitType := strings.ToLower(query.Get("type"))
	if kitType != "" && !utils.IsValidKitType(kitType) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidKitType, nil)
		return
	}

	// Initialize GORM query with Debug for detailed logging
	db := config.DB.Debug().
		Model(&models.Kit{}).
		Joins("JOIN users ON kits.created_by = users.id").
		Where("kits.is_deleted = ?", false)

	// Apply status filter
	if status == "active" {
		db = db.Where("kits.active_status = ?", true)
	} else if status == "inactive" {
		db = db.Where("kits.active_status = ?", false)
	}

	// Apply type filter
	if kitType != "" {
		db = db.Where("kits.type = ?", kitType)
	}

	// Apply search filter if searchText is provided
	if searchText != "" {
		searchPattern := "%" + searchText + "%"
		db = db.Where(
			"kits.extra_info ->> 'supplier_name' ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ? OR CAST(kits.created_at AS TEXT) ILIKE ? OR CAST(kits.quantity AS TEXT) ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
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

	// Apply sorting with special handling for JSONB fields
	orderClause := ""
	switch sortColumn {
	case "supplier_name":
		orderClause = fmt.Sprintf("extra_info->>'supplier_name' %s", sort)
	case "supplier_address":
		orderClause = fmt.Sprintf("extra_info->>'supplier_address' %s", sort)
	case "supplier_contact_number":
		orderClause = fmt.Sprintf("extra_info->>'supplier_contact_number' %s", sort)
	default:
		orderClause = fmt.Sprintf("%s %s", sortColumn, sort)
	}
	db = db.Order(orderClause)

	// Apply pagination
	offset := (page - 1) * perPage
	db = db.Limit(perPage).Offset(offset)

	// Fetch records
	var kits []models.Kit
	if err := db.Find(&kits).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToFetchRecords, nil)
		return
	}
	// Fetch associated users separately
	userIDs := make([]uint, len(kits))
	for i, kit := range kits {
		userIDs[i] = kit.CreatedBy
	}

	var users []models.User
	if err := config.DB.Where("id IN ?", userIDs).Find(&users).Error; err != nil {

		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToFetchRecords, nil)
		return
	}

	// Create a map of user IDs to users for quick lookup
	userMap := make(map[uint]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Prepare kit details with user information
	var kitDetails []KitDetail
	for _, k := range kits {
		user, exists := userMap[k.CreatedBy]
		if !exists {
			continue
		}

		detail := KitDetail{
			ID:                    k.ID,
			Type:                  k.Type,
			Quantity:              k.Quantity,
			SupplierName:          k.ExtraInfo.SupplierName,
			SupplierAddress:       k.ExtraInfo.SupplierAddress,
			SupplierContactNumber: k.ExtraInfo.SupplierContactNumber,
			Status:                k.Status,
			CreatedAt:             k.CreatedAt,
			CreatedBy: KitsUserProfile{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
		}
		kitDetails = append(kitDetails, detail)
	}

	// Prepare the response
	response := KitsListResponse{
		Page:         page,
		PerPage:      perPage,
		Sort:         sort,
		SortColumn:   sortColumn,
		SearchText:   searchText,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		Type:         kitType,
		Records:      kitDetails,
	}

	// Send the response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgKitsListFetchedSuccessfully, response)
}

// validateKitUpdateRequest validates the update request
func validateKitUpdateRequest(req KitUpdateRequest) error {
	// Validate Type if provided
	if req.Type != nil {
		if *req.Type == "" {
			return errors.New(utils.MsgTypeCannotBeEmptyIfProvided)
		}
		if !utils.IsValidKitType(*req.Type) {
			return errors.New(utils.MsgInvalidKitType)
		}
	}

	// Validate SupplierName if provided
	if req.SupplierName != nil {
		if *req.SupplierName == "" {
			return errors.New(utils.MsgSupplierNameCannotBeEmptyIfProvided)
		}
		if !utils.IsValidSupplierName(*req.SupplierName) {
			return errors.New(utils.MsgValidationSupplierName)
		}
	}

	// Validate SupplierContactNumber if provided
	if req.SupplierContactNumber != nil && *req.SupplierContactNumber != "" {
		if !utils.IsValidContactNumber(*req.SupplierContactNumber) {
			return errors.New(utils.MsgValidationContactNumber)
		}
	}

	// Validate Address if provided
	if req.SupplierAddress != nil && *req.SupplierAddress != "" {
		if len(*req.SupplierAddress) < 5 || len(*req.SupplierAddress) > 100 {
			return errors.New(utils.MsgValidationAddress)
		}
	}

	// Validate Quantity if provided
	if req.Quantity != nil {
		quantity, err := parseQuantity(req.Quantity)
		if err != nil {
			return err
		}
		if quantity < 0 {
			return errors.New(utils.MsgQuantityCannotBeNegative)
		}
		if quantity > 999999 {
			return errors.New(utils.MsgQuantityExceedsMaxValue)
		}
	}

	return nil
}

// UpdateKitHandler handles PATCH requests to update kit details
func UpdateKitHandler(w http.ResponseWriter, r *http.Request) {

	// Get kit ID from URL parameters
	vars := mux.Vars(r)
	kitID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidKitID, nil)
		return
	}

	// Parse the request body
	var req KitUpdateRequest
	allowedFields := []string{"type", "supplier_name", "supplier_contact_number", "supplier_address", "quantity", "status"}

	err = utils.ParseRequestBody(r, &req, allowedFields)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Trim spaces for string fields if they're provided
	if req.Type != nil {
		*req.Type = strings.TrimSpace(*req.Type)
	}
	if req.SupplierName != nil {
		*req.SupplierName = strings.TrimSpace(*req.SupplierName)
	}
	if req.SupplierContactNumber != nil {
		*req.SupplierContactNumber = strings.TrimSpace(*req.SupplierContactNumber)
	}
	if req.SupplierAddress != nil {
		*req.SupplierAddress = strings.TrimSpace(*req.SupplierAddress)
	}

	// Validate the request
	if err := validateKitUpdateRequest(req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Fetch existing kit
	var kit models.Kit
	if err := config.DB.Where("id = ? AND is_deleted = ?", kitID, false).First(&kit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgKitNotFound, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Update fields if provided
	if req.Type != nil {
		kit.Type = *req.Type
	}
	if req.Quantity != nil {
		quantity, _ := parseQuantity(req.Quantity) // Already validated
		kit.Quantity = quantity
	}
	if req.Status != nil {
		kit.Status = *req.Status
	}

	// Update ExtraInfo fields
	extraInfo := kit.ExtraInfo
	if req.SupplierName != nil {
		extraInfo.SupplierName = *req.SupplierName
	}
	if req.SupplierContactNumber != nil {
		extraInfo.SupplierContactNumber = *req.SupplierContactNumber
	}
	if req.SupplierAddress != nil {
		extraInfo.SupplierAddress = *req.SupplierAddress
	}
	kit.ExtraInfo = extraInfo

	// Save the updates
	if err := tx.Save(&kit).Error; err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Commit the transaction
	tx.Commit()

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgKitUpdatedSuccessfully, nil)
}

// DeleteKitHandler handles the soft deletion of kits
func DeleteKitHandler(w http.ResponseWriter, r *http.Request) {

	// Get kit ID from URL parameters
	vars := mux.Vars(r)
	kitID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidKitID, nil)
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Check if kit exists and is not already deleted
	var kit models.Kit
	if err := tx.Where("id = ? AND is_deleted = ?", kitID, false).First(&kit).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgKitAlreadyDeleted, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Perform soft delete
	updates := map[string]interface{}{
		"is_deleted": true,
		"status":     false, // Also set status to inactive
	}

	if err := tx.Model(&kit).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
		return
	}

	// Commit the transaction
	tx.Commit()

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgKitDeletedSuccessfully, nil)
}
