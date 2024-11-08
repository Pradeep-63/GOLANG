// routes/routes.go
package routes

import (
	"net/http"

	"theransticslabs/m/controllers"
	"theransticslabs/m/middlewares"
	"theransticslabs/m/utils"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Apply the Logging Middleware to all routes
	router.Use(middlewares.LoggingMiddleware)

	// Define Routes
	router.HandleFunc(utils.RouteWelcome, controllers.WelcomeHandler).Methods("GET")
	router.HandleFunc(utils.RouteLogin, controllers.LoginHandler).Methods("POST")                   // Login Route
	router.HandleFunc(utils.RouteForgetPassword, controllers.ForgetPasswordHandler).Methods("POST") // Login Route
	router.HandleFunc(utils.RouteEncryptProductDetails, controllers.EncryptProductDetails).Methods("POST")
	router.HandleFunc(utils.RouteVerifyProductDetails, controllers.VerifyProduct).Methods("POST")
	router.HandleFunc(utils.RouteProductPaymentDetails, controllers.OrderCreateHandler).Methods("POST")
	router.HandleFunc(utils.RoutePaymentSuccessPaypal, controllers.HandlePaymentSuccess).Methods("GET")

	// Protected Routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middlewares.AuthMiddleware)
	protected.Use(middlewares.CreatePermissionMiddleware())
	// Add protected routes here
	// e.g., protected.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	protected.HandleFunc(utils.RouteResetPassword, controllers.ResetPasswordHandler).Methods("PATCH")
	protected.HandleFunc(utils.RouteUpdateUser, controllers.UpdateUserInfoHandler).Methods("PATCH")
	// Private User Profile Route
	protected.HandleFunc(utils.RouteGetUserProfile, controllers.GetUserProfileHandler).Methods("GET")
	// Admin Users List Route
	protected.HandleFunc(utils.RouteGetAdminUserList, controllers.GetAdminUsersHandler).Methods("GET")
	protected.HandleFunc(utils.RouteCreateAdminUser, controllers.CreateUserHandler).Methods("POST")
	protected.HandleFunc(utils.RouteUpdateAdminUserProfile, controllers.UpdateUserProfileHandler).Methods("PATCH")
	protected.HandleFunc(utils.RouteUpdateAdminUserStatus, controllers.UpdateUserStatusHandler).Methods("PATCH")
	protected.HandleFunc(utils.RouteUpdateAdminUserPassword, controllers.UpdateUserPasswordHandler).Methods("PATCH")

	protected.HandleFunc(utils.RouteDeleteAdminUser, controllers.DeleteUserHandler).Methods("DELETE")
	protected.HandleFunc(utils.RouteKitInfo, controllers.CreateKitHandler).Methods("POST")
	protected.HandleFunc(utils.RouteKitInfo, controllers.GetKitsListHandler).Methods("GET")
	protected.HandleFunc(utils.RouteKitInfoID, controllers.UpdateKitHandler).Methods("PATCH")
	protected.HandleFunc(utils.RouteKitInfoID, controllers.DeleteKitHandler).Methods("DELETE")
	protected.HandleFunc(utils.RouteLogout, controllers.LogoutHandler).Methods("DELETE")

	// Handle 404
	router.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)

	return router
}
