// utils/paypal.go

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"theransticslabs/m/config"
	"theransticslabs/m/models"
)

// Define the response structure for PayPal access token
type PayPalAccessTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

// Define the response structure for PayPal order
type PayPalOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Links  []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

// PayPalCaptureResponse represents the response from PayPal's capture endpoint
type PayPalCaptureResponse struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PaymentSource struct {
		PayPal struct {
			EmailAddress string `json:"email_address"`
			AccountID    string `json:"account_id"`
		} `json:"paypal"`
	} `json:"payment_source"`
}

func GetPayPalAccessToken() (string, error) {
	url := config.AppConfig.PaypalAPIUrl + "/v1/oauth2/token"
	clientID := config.AppConfig.PaypalClientID
	secret := config.AppConfig.PaypalClientSecret

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res PayPalAccessTokenResponse
	json.NewDecoder(resp.Body).Decode(&res)

	if res.AccessToken == "" {
		return "", errors.New("failed to get PayPal token")
	}
	return res.AccessToken, nil
}

func CreatePayPalOrder(paymentID uint, amount float64, accessToken string, orderDetails *models.Order, customer *models.Customer) (PayPalOrderResponse, error) {
	url := config.AppConfig.PaypalAPIUrl + "/v2/checkout/orders"

	returnURL := fmt.Sprintf("%s/payment/status?payment_id=%d", config.AppConfig.ApiUrl, paymentID)
	cancelURL := fmt.Sprintf("%s/payment/cancel?payment_id=%d", config.AppConfig.AppUrl, paymentID)

	// Calculate unit price and subtotal
	unitPrice := amount / float64(orderDetails.Quantity)
	subtotal := unitPrice * float64(orderDetails.Quantity)

	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"application_context": map[string]interface{}{
			"return_url":          returnURL,
			"cancel_url":          cancelURL,
			"shipping_preference": "NO_SHIPPING",
			"user_action":         "PAY_NOW",
			"brand_name":          "Theranostics DNA",

			"payment_method": map[string]interface{}{
				"payer_selected":            "PAYPAL",
				"payee_preferred":           "IMMEDIATE_PAYMENT_REQUIRED",
				"standard_entry_class_code": "WEB",
			},
		},

		"purchase_units": []map[string]interface{}{
			{
				"reference_id": strconv.FormatUint(uint64(paymentID), 10),
				"description":  fmt.Sprintf("Order #%d", orderDetails.ID),
				"custom_id":    fmt.Sprintf("ORDER_%d", orderDetails.ID),
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", amount),
					"breakdown": map[string]interface{}{
						"item_total": map[string]string{
							"currency_code": "USD",
							"value":         fmt.Sprintf("%.2f", subtotal),
						},
					},
				},
				"items": []map[string]interface{}{
					{
						"name": orderDetails.ProductName,
						"description": fmt.Sprintf("Product Details:\n- Name: %s\n- Quantity: %d\n- Price per unit: $%.2f\n- Subtotal: $%.2f",
							orderDetails.ProductName,
							orderDetails.Quantity,
							unitPrice,
							subtotal,
						),
						"quantity": strconv.Itoa(orderDetails.Quantity),
						"unit_amount": map[string]string{
							"currency_code": "USD",
							"value":         fmt.Sprintf("%.2f", unitPrice),
						},
					},
				},
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return PayPalOrderResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PayPalOrderResponse{}, err
	}
	defer resp.Body.Close()

	var res PayPalOrderResponse
	json.NewDecoder(resp.Body).Decode(&res)
	return res, nil
}

// CapturePayPalPayment captures a previously authorized PayPal payment
func CapturePayPalPayment(orderID string, accessToken string) error {
	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", config.AppConfig.PaypalAPIUrl, orderID)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return fmt.Errorf("failed to create capture request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Prefer", "return=representation")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send capture request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read capture response: %w", err)
	}

	// Check for successful status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to capture payment. Status: %d, Body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var captureResponse PayPalCaptureResponse
	if err := json.Unmarshal(body, &captureResponse); err != nil {
		return fmt.Errorf("failed to parse capture response: %w", err)
	}

	// Verify capture status
	if captureResponse.Status != "COMPLETED" {
		return fmt.Errorf("payment capture failed. Status: %s", captureResponse.Status)
	}

	return nil
}
