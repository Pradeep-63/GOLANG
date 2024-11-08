package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"theransticslabs/m/utils"
)

type Product struct {
	Name        string  `json:"name" validate:"required" form:"name"`
	Description string  `json:"description" form:"description"`
	Image       string  `json:"image" form:"image"`
	Price       float64 `json:"price" validate:"required,gt=0" form:"price"`
}

type EncryptedData struct {
	Data string `json:"data" validate:"required" form:"data"`
}

func EncryptProductDetails(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req Product

	// Use the common request parser
	// Define allowed fields for this request
	allowedFields := []string{"name", "description", "image", "price"}

	err := utils.ParseRequestBody(r, &req, allowedFields) // Assuming this validates required fields
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Marshal struct to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgErrorEncodingProductDetails, nil)
		return
	}

	// Encrypt the JSON string
	encryptedData, err := utils.Encrypt(string(jsonData))
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgErrorEncryptingProductDetails, map[string]string{"error": err.Error()})
		return
	}
	log.Println(req)
	// Success response
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgProductDetailsEncryptedSuccessfully, encryptedData)
}

// VerifyProduct validates and verifies an encrypted product data request
func VerifyProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into EncryptedData struct
	var req EncryptedData

	// Define allowed fields for this request
	allowedFields := []string{"data"}

	// Use the common request parser to validate allowed fields
	if err := utils.ParseRequestBody(r, &req, allowedFields); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// Validate that 'data' is present and non-empty in the request
	if req.Data == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgMissingOrEmptyEncryptedDataParameter, nil)
		return
	}

	// Attempt to decrypt the encrypted data
	decryptedData, err := utils.Decrypt(req.Data)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgFailedToDecryptData, map[string]string{"error": err.Error()})
		return
	}

	// Unmarshal decrypted JSON into a Product struct
	var product Product
	if err := json.Unmarshal([]byte(decryptedData), &product); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidProductDataFormat, map[string]string{"error": err.Error()})
		return
	}

	// Validate required fields in the decrypted product data
	// Validate product name
	if product.Name == "" || !utils.IsValidProductName(product.Name) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidProductName, nil)
		return
	}

	// Validate product price
	if product.Price <= 0 || !utils.IsValidPrice(strconv.FormatFloat(product.Price, 'f', 2, 64)) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgProductPriceMustBeGreaterThanZero, product)
		return
	}

	// Validate product description if provided
	if product.Description != "" && len(product.Description) > 1000 {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgProductDescriptionTooLong, nil)
		return
	}

	// Validate product image as either base64 or URL
	if product.Image != "" && !(utils.IsValidBase64Image(product.Image) || utils.IsValidImageURL(product.Image)) {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgInvalidImageFormat, nil)
		return
	}

	// Return the decrypted and validated product data
	utils.JSONResponse(w, http.StatusOK, true, utils.MsgProductVerifiedSuccessfully, product)
}
