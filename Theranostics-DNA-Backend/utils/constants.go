// utils/constantsgo
package utils

const (
	// Route Names

	// Public
	RouteWelcome               = "/"
	RouteLogin                 = "/login"
	RouteForgetPassword        = "/user/forgot-password"
	RouteEncryptProductDetails = "/encrypt-product"
	RouteVerifyProductDetails  = "/verify-product"
	RouteProductPaymentDetails = "/order-payment"
	RoutePaymentSuccessPaypal  = "/payment/status"

	// Private
	RouteLogout                  = "/logout"
	RouteResetPassword           = "/user/reset-password"
	RouteUpdateUser              = "/user"
	RouteGetUserProfile          = "/user/profile"
	RouteGetAdminUserList        = "/staff"
	RouteCreateAdminUser         = "/staff"
	RouteUpdateAdminUserProfile  = "/staff/{id}/details"
	RouteUpdateAdminUserPassword = "/staff/{id}/password"
	RouteUpdateAdminUserStatus   = "/staff/{id}/status"
	RouteDeleteAdminUser         = "/staff/{id}"
	RouteKitInfo                 = "/kits"
	RouteKitInfoID               = "/kits/{id}"

	// API Messages
	MsgWelcome                                 = "Welcome to the project!"
	MsgDatabaseConnected                       = "Database connection established successfully."
	MsgDatabaseMigrated                        = "Database migration completed successfully."
	MsgServerStarted                           = "Server is running on port %s."
	MsgInternalServerError                     = "An internal server error has occurred."
	MsgEndpointNotFound                        = "The requested endpoint could not be found."
	MsgUserNotFound                            = "User not found in the current context."
	MsgInvalidRequest                          = "The request is invalid."
	MsgUserNotAuthenticated                    = "User is not authenticated."
	MsgInvalidJSONData                         = "Invalid JSON data provided."
	MsgInvalidFormData                         = "Invalid form data provided."
	MsgUnsupportedContentType                  = "Unsupported content type. Only JSON and form-urlencoded data are allowed."
	MsgOnlyEmailPasswordAllowed                = "Only email and password fields are permitted."
	MsgOnlyAllowedFieldsAllowed                = "Only specific fields are permitted."
	MsgAccessDeniedForOtherThanAdminSuperAdmin = "You are not authorized to view this page."

	// Seeding Messages
	MsgRolesSeededSuccessfully = "Roles seeded successfully."
	MsgRoleCreated             = "Role '%s' created successfully."
	MsgRoleAlreadyExists       = "Role '%s' already exists."
	MsgFailedToCreateRole      = "Failed to create role '%s': %v."
	MsgFailedToCheckRole       = "Error checking role '%s': %v."
	MsgSeedingCompleted        = "Seeding process completed successfully."

	// User Seeding Messages
	MsgUsersSeededSuccessfully = "Users seeded successfully."
	MsgUserCreated             = "User '%s' created successfully."
	MsgUserAlreadyExists       = "User '%s' already exists."
	MsgFailedToCreateUser      = "Failed to create user '%s': %v."
	MsgFailedToCheckUser       = "Error checking user '%s': %v."

	// Middlewares Messages
	MsgAuthHeaderMissing            = "Authorization header missing"
	MsgInvalidAuthHeaderFormat      = "Invalid authorization header format"
	MsgInvalidTokenClaims           = "Invalid token claims."
	MsgUserDoesNotExist             = "User does not exist."
	MsgUserAccountInactiveOrDeleted = "User account is inactive or deleted."
	MsgUnauthorizedUser             = "Unauthorized user."

	// Login Messages
	MsgLoginSuccess               = "Login Successful."
	MsgInvalidCredentials         = "Invalid email or password."
	MsgInvalidEmailFormat         = "The provided email format is invalid."
	MsgInvalidPasswordFormat      = "Password must be 8-20 characters long and include at least one uppercase letter, one lowercase letter, one number, and one special character."
	MsgMissingEmailOrPassword     = "Both email and password are required."
	MsgTokenCreationFailed        = "Failed to create the JWT token."
	MsgAuthorizationHeaderMissing = "Authorization header is missing."
	MsgInvalidAuthorizationHeader = "The format of the authorization header is invalid."
	MsgInvalidOrExpiredToken      = "The token is either invalid or expired."
	MsgUserInactive               = "The account is currently inactive."
	MsgTokenSaveFailed            = "Failed to save the token. Please try again."
	MsgLogoutSuccess              = "Successfully logged out."

	// Forget Password Messages
	MsgMissingEmail       = "Email is required."
	MsgEmailNotExist      = "The provided email does not exist."
	MsgFailedSentEmail    = "Failed to send email. Please try again later."
	MsgForgetSuccessfully = "A new password has been sent successfully. Please check your email."

	// Reset Password Messages
	MsgPasswordSame              = "The new password must differ from the old password."
	MsgInvalidOldPassword        = "The provided old password is invalid."
	MsgFailedToHashNewPassword   = "Failed to hash the new password."
	MsgFailedToUpdateNewPassword = "Failed to update the new password."
	MsgResetPasswordSuccessfully = "The password has been reset successfully."
	MsgOldNewPasswordRequired    = "Both old and new passwords are required."

	// Update User Messages
	MsgAtleastPassSomeData     = "First name must be provided."
	MsgInvalidRoleId           = "The provided role ID is invalid."
	MsgUserUpdateSuccessfully  = "User information updated successfully."
	MsgUserFetchedSuccessfully = "User profile fetched successfully."

	// Fetch Admin User Details Messages
	MsgAccessDeniedForOtherUser    = "You are not authorized to view this page."
	MsgStatusInvalid               = "The status filter is invalid. Allowed values: all, active, inactive."
	MsgFailedToCountRecords        = "Failed to count the records."
	MsgFailedToFetchRecords        = "Failed to fetch the records."
	MsgUserListFetchedSuccessfully = "Users fetched successfully."
	MsgAdminRoleNotFound           = "Admin role not found."

	// Create Admin User Messages
	MsgFirstNameValidation     = "First name must be alphabetic and between 3 to 20 characters long."
	MsgLastNameValidation      = "Last name must be alphabetic and between 3 to 20 characters long (if provided)."
	MsgEmailValidation         = "Email format is invalid or exceeds the maximum length."
	MsgEmailAlreadyInUse       = "The email is already in use."
	MsgFailedToSecurePassword  = "Failed to secure the password."
	MsgFailedCreateUser        = "Failed to create the user."
	MsgFailedToSendEmail       = "Failed to send email."
	MsgUserCreatedSuccessfully = "User created successfully."
	MsgFirstNameEmailRequired  = "Both first name and email are required."

	// Update Admin User Messages
	MsgPasswordValidation          = "Password must be between 8 to 20 characters and include one uppercase letter, one lowercase letter, one number, and one special character."
	MsgFailedToUpdateAdminUser     = "Failed to update the user."
	MsgAdminUserUpdateSuccessfully = "User updated successfully."

	// Delete Admin User Messages
	MsgUserAlreadyDeleted     = "User not found or has already been deleted."
	MsgFailedToDeleteUser     = "Failed to delete the user."
	MsgUserDeleteSuccessfully = "User deleted successfully."

	// Kit Related Messages
	MsgKitCreatedSuccessfully      = "Kit added successfully."
	MsgInvalidKitType              = "Invalid kit type. Must be either 'blood' or 'saliva'."
	MsgValidationSupplierName      = "Supplier name should be alphabetic and between 3 to 50 characters."
	MsgValidationContactNumber     = "Contact number should be numeric and between 10 to 15 digits."
	MsgValidationAddress           = "Address must be between 5 to 100 characters."
	MsgAllFieldsOfkitRequired      = "All fields are required."
	MsgKitsListFetchedSuccessfully = "Kits list fetched successfully."

	MsgKitUpdatedSuccessfully              = "Kit details updated successfully."
	MsgKitDeletedSuccessfully              = "Kit deleted successfully."
	MsgKitNotFound                         = "Kit not found."
	MsgKitAlreadyDeleted                   = "Kit not found or already deleted."
	MsgInvalidQueryParameters              = "Invalid query parameters."
	MsgInvalidPageParameter                = "Invalid 'page' parameter."
	MsgInvalidPerPageParameter             = "Invalid 'per_page' parameter."
	MsgInvalidSortParameter                = "Invalid 'sort' parameter."
	MsgInvalidSortColumnParameter          = "Invalid 'sort_column' parameter."
	MsgQuantityCannotBeEmpty               = "Quantity cannot be empty."
	MsgInvalidQuantityFormat               = "Invalid quantity format."
	MsgInvalidQuantityType                 = "Invalid quantity type."
	MsgQuantityCannotBeNegative            = "Quantity cannot be negative."
	MsgQuantityExceedsMaxValue             = "Quantity exceeds maximum allowed value."
	MsgTypeCannotBeEmptyIfProvided         = "Type cannot be empty if provided."
	MsgSupplierNameCannotBeEmptyIfProvided = "Supplier name cannot be empty if provided."
	MsgInvalidKitID                        = "Invalid kit ID."

	// Payment Related Message
	MsgFailedToStartTransaction       = "Failed to start transaction"
	MsgFailedToProcessCustomer        = "Failed to process customer: %s"
	MsgFailedToInitializePayment      = "Failed to initialize payment: %s"
	MsgFailedToCommitTransaction      = "Failed to commit transaction"
	MsgOrderCreatedSuccessfully       = "Order created successfully"
	MsgInvalidFirstName               = "Invalid first name: must be 3-50 characters long and contain only letters"
	MsgInvalidLastName                = "Invalid last name: must be 3-50 characters long and contain only letters"
	MsgInvalidPhoneNumber             = "Invalid phone number: must be 10-15 digits"
	MsgInvalidCountryName             = "Invalid country name: must be 3-50 characters"
	MsgInvalidStreetAddress           = "Invalid street address: must be 5-255 characters"
	MsgInvalidTownCity                = "Invalid town/city: must be 5-100 characters"
	MsgInvalidRegion                  = "Invalid region: must be 3-100 characters"
	MsgInvalidPostcode                = "Invalid postcode: must be 3-20 characters"
	MsgInvalidProductName             = "Invalid product name: must be 3-100 characters and contain only letters, numbers, and basic punctuation"
	MsgProductDescriptionTooLong      = "Product description too long: must not exceed 1000 characters"
	MsgInvalidProductImage            = "Invalid product image: must be a valid base64 encoded string"
	MsgInvalidProductPrice            = "Invalid product price: must be a positive number with a maximum of 2 decimal places"
	MsgInvalidQuantity                = "Invalid quantity: must be a positive integer"
	MsgInvalidPriceFormat             = "Invalid price format"
	MsgMissingPaymentInformation      = "Missing payment information"
	MsgFailedToStartTransactionAgain  = "Failed to start transaction"
	MsgPaymentVerificationFailed      = "Payment verification failed: %s"
	MsgFailedToCompletePaymentProcess = "Failed to complete payment process"
	MsgPaymentCompletedSuccessfully   = "Payment completed successfully"
	MsgPaymentNotFound                = "payment not found"
	MsgFailedToUpdatePayment          = "Failed to update payment"
	MsgOrderNotFound                  = "Order not found"
	MsgInvalidImageFormat             = "Invalid image format."

	MsgErrorEncodingProductDetails          = "Error encoding product details"
	MsgErrorEncryptingProductDetails        = "Error encrypting product details"
	MsgProductDetailsEncryptedSuccessfully  = "Product details encrypted successfully"
	MsgMissingOrEmptyEncryptedDataParameter = "Missing or empty encrypted data parameter"
	MsgFailedToDecryptData                  = "Failed to decrypt data"
	MsgInvalidProductDataFormat             = "Invalid product data format"
	MsgProductNameIsRequired                = "Product name is required"
	MsgProductPriceMustBeGreaterThanZero    = "Product price must be greater than zero"
	MsgProductVerifiedSuccessfully          = "Product verified successfully"
)
