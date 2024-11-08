// middlewares/auth_middleware.go
package middlewares

import (
	"context"
	"net/http"
	"strings"

	"theransticslabs/m/config"
	"theransticslabs/m/models"
	"theransticslabs/m/utils"

	"gorm.io/gorm"
)

// Define a type for context keys to avoid collisions
type contextKey string

const userContextKey = contextKey("user")

// AuthMiddleware verifies the JWT token, ensures the user exists and is active,
// and sets the user information in the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgAuthHeaderMissing, nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgInvalidAuthHeaderFormat, nil)
			return
		}

		tokenString := parts[1]

		// Validate JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgInvalidOrExpiredToken, nil)
			return
		}

		// Extract user information from claims
		userID, ok := claims["id"]
		if !ok {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgInvalidTokenClaims, nil)
			return
		}

		// Fetch the user from the database
		var user models.User
		result := config.DB.Preload("Role").First(&user, userID)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserDoesNotExist, nil)
				return
			}
			utils.JSONResponse(w, http.StatusInternalServerError, false, utils.MsgInternalServerError, nil)
			return
		}

		// Check if user is active and not deleted
		if !user.ActiveStatus || user.IsDeleted {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserAccountInactiveOrDeleted, nil)
			return
		}

		// Compare the token in the request with the one stored in the database
		if user.Token != tokenString {
			utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUnauthorizedUser, nil)
			return
		}

		// Set user information in the context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the user information from the context.
// It returns the user and a boolean indicating whether the user was found.
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(models.User)
	if !ok {
		return nil, false
	}
	return &user, true
}
