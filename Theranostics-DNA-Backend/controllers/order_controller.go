// controllers/order_controller.go

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"theransticslabs/m/config"
	"theransticslabs/m/emails"
	"theransticslabs/m/models"
	"theransticslabs/m/utils"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type OrderRequest struct {
	FirstName          string `json:"first_name" form:"first_name" validate:"required,max=50,min=3"`
	LastName           string `json:"last_name" form:"last_name" validate:"omitempty,max=50,min=3"`
	Email              string `json:"email" form:"email" validate:"required,email,max=100"`
	PhoneNumber        string `json:"phone_number" form:"phone_number" validate:"required,max=15,min=10"`
	Country            string `json:"country" form:"country" validate:"required,max=50,min=3"`
	StreetAddress      string `json:"street_address" form:"street_address" validate:"required,max=255,min=5"`
	TownCity           string `json:"town_city" form:"town_city" validate:"required,max=100,min=5"`
	Region             string `json:"region" form:"region" validate:"omitempty,max=100,min=3"`
	Postcode           string `json:"postcode" form:"postcode" validate:"omitempty,max=20,min=3"`
	ProductName        string `json:"product_name" form:"product_name" validate:"required,max=100,min=3"`
	ProductDescription string `json:"product_description" form:"product_description" validate:"omitempty"`
	ProductImage       string `json:"product_image" form:"product_image" validate:"omitempty,base64"`
	ProductPrice       string `json:"product_price" form:"product_price" validate:"required,numeric,gt=0"`
	Quantity           string `json:"quantity" form:"quantity" validate:"required,numeric,min=1"`
}

type PaymentResponse struct {
	OrderID    uint   `json:"order_id"`
	PaymentID  uint   `json:"payment_id"`
	PaymentURL string `json:"payment_url"`
}

// OrderCreateHandler processes the complete order flow
func OrderCreateHandler(w http.ResponseWriter, r *http.Request) {

	var req OrderRequest
	allowedFields := []string{"first_name", "last_name", "email", "phone_number", "country",
		"street_address", "town_city", "region", "postcode", "product_name",
		"product_description", "product_image", "product_price", "quantity"}

	if err := utils.ParseRequestBody(r, &req, allowedFields); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}
	log.Println(req)

	if err := validateOrderRequest(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToStartTransaction, nil)
		return
	}
	defer tx.Rollback()

	// 2. Process customer
	customer, err := processCustomer(tx, &req)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, fmt.Sprintf(utils.MsgFailedToProcessCustomer, err.Error()), nil)
		return
	}

	// 3. Create order
	order, err := processOrderDetails(tx, customer, &req)
	if err != nil {
		tx.Rollback()
		utils.JSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	// 4. Initialize PayPal payment
	paymentURL, paymentID, err := initializePayPalPayment(tx, order, customer)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, fmt.Sprintf(utils.MsgFailedToInitializePayment, err.Error()), nil)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToCommitTransaction, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgOrderCreatedSuccessfully, PaymentResponse{
		OrderID:    order.ID,
		PaymentID:  paymentID,
		PaymentURL: paymentURL,
	})

}

// Additional validation functions
func validateOrderRequest(req *OrderRequest) error {
	// Existing validations
	if !utils.IsValidFirstName(req.FirstName) {
		return fmt.Errorf(utils.MsgInvalidFirstName)
	}
	if req.LastName != "" && !utils.IsValidLastName(req.LastName) {
		return fmt.Errorf(utils.MsgInvalidLastName)
	}
	if !utils.IsValidEmail(req.Email) {
		return fmt.Errorf(utils.MsgInvalidEmailFormat)
	}
	if !utils.IsValidContactNumber(req.PhoneNumber) {
		return fmt.Errorf(utils.MsgInvalidPhoneNumber)
	}
	if len(req.Country) < 3 || len(req.Country) > 50 {
		return fmt.Errorf(utils.MsgInvalidCountryName)
	}
	if len(req.StreetAddress) < 5 || len(req.StreetAddress) > 255 {
		return fmt.Errorf(utils.MsgInvalidStreetAddress)
	}
	if len(req.TownCity) < 5 || len(req.TownCity) > 100 {
		return fmt.Errorf(utils.MsgInvalidTownCity)
	}
	if req.Region != "" && (len(req.Region) < 3 || len(req.Region) > 100) {
		return fmt.Errorf(utils.MsgInvalidRegion)
	}
	if req.Postcode != "" && (len(req.Postcode) < 3 || len(req.Postcode) > 20) {
		return fmt.Errorf(utils.MsgInvalidPostcode)
	}

	// Product validations
	if !utils.IsValidProductName(req.ProductName) {
		return fmt.Errorf(utils.MsgInvalidProductName)
	}

	// Product description validation (optional)
	if req.ProductDescription != "" && len(req.ProductDescription) > 1000 {
		return fmt.Errorf(utils.MsgProductDescriptionTooLong)
	}

	// Product image validation (optional, base64)
	if req.ProductImage != "" && !(utils.IsValidBase64Image(req.ProductImage) || utils.IsValidImageURL(req.ProductImage)) {
		return fmt.Errorf(utils.MsgInvalidProductImage)
	}

	// Product price validation
	if !utils.IsValidPrice(req.ProductPrice) {
		return fmt.Errorf(utils.MsgInvalidProductPrice)
	}

	// Quantity validation
	if !utils.IsValidQuantity(req.Quantity) {
		return fmt.Errorf(utils.MsgInvalidQuantity)
	}

	return nil
}

func processCustomer(tx *gorm.DB, req *OrderRequest) (*models.Customer, error) {
	var customer models.Customer
	result := tx.Where("email = ?", req.Email).First(&customer)

	if result.Error == gorm.ErrRecordNotFound {
		customer = models.Customer{
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Email:         req.Email,
			PhoneNumber:   req.PhoneNumber,
			Country:       req.Country,
			StreetAddress: req.StreetAddress,
			TownCity:      req.TownCity,
			Region:        req.Region,
			Postcode:      req.Postcode,
		}
		if err := tx.Create(&customer).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &customer, nil
}

func processOrderDetails(tx *gorm.DB, customer *models.Customer, req *OrderRequest) (*models.Order, error) {
	price, err := strconv.ParseFloat(req.ProductPrice, 64)
	if err != nil {
		return nil, fmt.Errorf(utils.MsgInvalidPriceFormat)
	}

	quantity, err := strconv.Atoi(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf(utils.MsgInvalidQuantityFormat)
	}

	order := models.Order{
		CustomerID:         customer.ID,
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		ProductImage:       req.ProductImage,
		ProductPrice:       price,
		Quantity:           quantity,
		TotalPrice:         price * float64(quantity),
		PaymentStatus:      "Pending",
		OrderStatus:        "Pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func initializePayPalPayment(tx *gorm.DB, order *models.Order, customer *models.Customer) (string, uint, error) {
	// Create payment record
	payment := &models.Payment{
		OrderID:       order.ID,
		PaymentStatus: "Pending",
		Amount:        order.TotalPrice,
	}
	if err := tx.Create(payment).Error; err != nil {
		return "", 0, err
	}

	// Initialize PayPal payment
	accessToken, err := utils.GetPayPalAccessToken()
	if err != nil {
		return "", 0, err
	}

	paypalOrder, err := utils.CreatePayPalOrder(payment.ID, order.TotalPrice, accessToken, order, customer)
	if err != nil {
		log.Println(err)

		return "", 0, err
	}

	// Update payment with PayPal transaction ID
	payment.TransactionID = paypalOrder.ID
	if err := tx.Save(payment).Error; err != nil {
		return "", 0, err
	}

	// Get PayPal approval URL
	var approvalURL string
	for _, link := range paypalOrder.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	return approvalURL, payment.ID, nil
}

func HandlePaymentSuccess(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("payment_id")
	paypalOrderID := r.URL.Query().Get("token")

	if paymentID == "" || paypalOrderID == "" {
		utils.JSONResponse(w, http.StatusBadRequest, false, utils.MsgMissingPaymentInformation, nil)
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToStartTransactionAgain, nil)
		return
	}
	defer tx.Rollback()

	// Verify and capture PayPal payment
	if err := captureAndVerifyPayment(paypalOrderID); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, fmt.Sprintf(utils.MsgPaymentVerificationFailed, err.Error()), nil)
		return
	}

	// Update payment and order status
	if err := updatePaymentAndOrderStatus(tx, paymentID); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	// Generate invoice and send emails
	if err := handleSuccessfulPayment(tx, paymentID); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgFailedToCompletePaymentProcess, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, utils.MsgPaymentCompletedSuccessfully, nil)
}

func captureAndVerifyPayment(paypalOrderID string) error {
	accessToken, err := utils.GetPayPalAccessToken()
	if err != nil {
		return err
	}
	return utils.CapturePayPalPayment(paypalOrderID, accessToken)
}

func updatePaymentAndOrderStatus(tx *gorm.DB, paymentID string) error {
	var payment models.Payment
	if err := tx.First(&payment, paymentID).Error; err != nil {
		return fmt.Errorf(utils.MsgPaymentNotFound)
	}

	payment.PaymentStatus = "Completed"
	if err := tx.Save(&payment).Error; err != nil {
		return fmt.Errorf(utils.MsgFailedToUpdatePayment)
	}

	var order models.Order
	if err := tx.First(&order, payment.OrderID).Error; err != nil {
		return fmt.Errorf(utils.MsgOrderNotFound)
	}

	order.PaymentStatus = "Completed"
	order.OrderStatus = "Processing"
	return tx.Save(&order).Error
}

func handleSuccessfulPayment(tx *gorm.DB, paymentID string) error {
	var (
		payment  models.Payment
		order    models.Order
		customer models.Customer
	)

	// Get all necessary data
	if err := tx.First(&payment, paymentID).Error; err != nil {
		return err
	}
	if err := tx.First(&order, payment.OrderID).Error; err != nil {
		return err
	}
	if err := tx.First(&customer, order.CustomerID).Error; err != nil {
		return err
	}

	// Generate invoice
	invoicePath, err := generateInvoice(&payment, &customer)
	if err != nil {
		return err
	}

	// Save invoice
	invoice := models.Invoice{
		PaymentID:   payment.ID,
		InvoiceLink: invoicePath,
		Price:       payment.Amount,
		InvoiceID:   strconv.FormatUint(uint64(payment.ID), 10),
	}
	if err := tx.Create(&invoice).Error; err != nil {
		return err
	}

	// Send emails
	return sendConfirmationEmails(&customer, &order, &invoice)
}

func generateInvoice(payment *models.Payment, customer *models.Customer) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add invoice content
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %d", payment.ID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s %s", customer.FirstName, customer.LastName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Amount: $%.2f", payment.Amount))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02")))

	// Save PDF
	invoicePath := filepath.Join("public/invoices", fmt.Sprintf("invoice_%d.pdf", payment.ID))
	if err := os.MkdirAll(filepath.Dir(invoicePath), 0755); err != nil {
		return "", err
	}

	// Generate the PDF and save to file
	if err := pdf.OutputFileAndClose(invoicePath); err != nil {
		return "", err
	}

	return filepath.Join("invoices", fmt.Sprintf("invoice_%d.pdf", payment.ID)), nil
}

func sendConfirmationEmails(customer *models.Customer, order *models.Order, invoice *models.Invoice) error {
	// Send customer email
	customerEmail := emails.CustomerOrderConfirmationEmail(
		customer.FirstName,
		customer.LastName,
		order.ProductName,
		order.Quantity,
		order.TotalPrice,
		invoice.InvoiceLink,
		config.AppConfig.AppUrl,
	)
	if err := config.SendEmail([]string{customer.Email}, "Order Confirmation", customerEmail); err != nil {
		return err
	}

	// Send admin notification if configured
	if adminEmail := os.Getenv("ADMIN_EMAIL"); adminEmail != "" {
		adminEmail := emails.NewOrderNotificationEmail(
			customer.FirstName,
			customer.LastName,
			customer.Email,
			order.ProductName,
			order.Quantity,
			order.TotalPrice,
			invoice.InvoiceLink,
		)
		if err := config.SendEmail([]string{adminEmail}, "New Order Received", adminEmail); err != nil {
			return err
		}
	}

	return nil
}
