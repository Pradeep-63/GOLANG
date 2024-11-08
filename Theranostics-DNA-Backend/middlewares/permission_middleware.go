// middlewares/permission_middleware.go
package middlewares

import (
	"net/http"
	"regexp"
	"strings"
	"theransticslabs/m/utils"
)

// RoutePermission defines which roles can access a specific route
type RoutePermission struct {
	Route  string   // Route path from utils.Routes
	Roles  []string // Allowed roles for this route
	Method string   // HTTP method (GET, POST, etc.)
}

// Define allowed roles
var AllowedRoles = []string{
	"super-admin",
	"admin",
	"user",
}

// RoutePermissions maps routes to their permitted roles
// Using the route constants from utils.Routes for consistency
var RoutePermissions = []RoutePermission{

	{
		Route:  "/api" + utils.RouteLogout, // "/api/staff"
		Roles:  []string{"super-admin", "admin", "user"},
		Method: http.MethodDelete,
	},
	{
		Route:  "/api" + utils.RouteResetPassword, // "/api/staff"
		Roles:  []string{"super-admin", "admin", "user"},
		Method: http.MethodPatch,
	},
	{
		Route:  "/api" + utils.RouteUpdateUser, // "/api/staff"
		Roles:  []string{"super-admin", "admin", "user"},
		Method: http.MethodPatch,
	},
	{
		Route:  "/api" + utils.RouteGetUserProfile, // "/api/user/profile"
		Roles:  []string{"super-admin", "admin", "user"},
		Method: http.MethodGet,
	},
	{
		Route:  "/api" + utils.RouteGetAdminUserList, // "/api/staff"
		Roles:  []string{"super-admin"},
		Method: http.MethodGet,
	},
	{
		Route:  "/api" + utils.RouteCreateAdminUser, // "/api/staff"
		Roles:  []string{"super-admin"},
		Method: http.MethodPost,
	},
	{
		Route:  "/api" + utils.RouteUpdateAdminUserProfile, // "/api/staff/{id}/details"
		Roles:  []string{"super-admin"},
		Method: http.MethodPatch,
	},
	{
		Route:  "/api" + utils.RouteUpdateAdminUserPassword, // "/api/staff/{id}/password"
		Roles:  []string{"super-admin"},
		Method: http.MethodPatch,
	},
	{
		Route:  "/api" + utils.RouteUpdateAdminUserStatus, // "/api/staff/{id}/status"
		Roles:  []string{"super-admin"},
		Method: http.MethodPatch,
	},
	{
		Route:  "/api" + utils.RouteDeleteAdminUser, // "/api/staff/{id}"
		Roles:  []string{"super-admin"},
		Method: http.MethodDelete,
	},
	{
		Route:  "/api" + utils.RouteKitInfo, // "/api/kits"
		Roles:  []string{"super-admin", "admin"},
		Method: "", // Empty means allow all methods
	},
	{
		Route:  "/api" + utils.RouteKitInfoID, // "/api/kits/{id}"
		Roles:  []string{"super-admin", "admin"},
		Method: "", // Empty means allow all methods
	},
}

// CheckPermission checks if a user's role has permission for the given route and method
func CheckPermission(userRole, route, method string) bool {
	// Convert role to lowercase for consistent comparison
	userRole = strings.ToLower(userRole)

	// First check if the role is valid
	isValidRole := false
	for _, role := range AllowedRoles {
		if strings.ToLower(role) == userRole {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return false
	}

	// Check permissions for the route
	for _, permission := range RoutePermissions {
		// Convert route pattern to regex for matching
		routePattern := strings.Replace(permission.Route, "{id}", "[^/]+", -1)

		matched, _ := regexp.MatchString("^"+routePattern+"$", route)

		if matched {
			// If method is specified, it must match
			if permission.Method != "" && permission.Method != method {
				continue
			}

			// Check if user's role is allowed
			for _, allowedRole := range permission.Roles {
				if strings.ToLower(allowedRole) == userRole {
					return true
				}
			}
		}
	}

	return false
}

// CreatePermissionMiddleware creates a middleware that checks permissions
func CreatePermissionMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (set by AuthMiddleware)
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				utils.JSONResponse(w, http.StatusUnauthorized, false, utils.MsgUserNotAuthenticated, nil)
				return
			}

			// Check if user has permission for this route
			hasPermission := CheckPermission(user.Role.Name, r.URL.Path, r.Method)
			if !hasPermission {
				utils.JSONResponse(w, http.StatusForbidden, false, utils.MsgAccessDeniedForOtherUser, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
