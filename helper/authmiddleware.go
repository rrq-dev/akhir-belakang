package helper

import (
	"akhir-belakang/model"
	"context"
	"log"
	"net/http"
)

// RoleMiddleware validates the user's role against allowed roles
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			tokenString, err := GetTokenFromHeader(r)
			if err != nil {
				log.Printf("Token retrieval error: %v", err)
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Parse and validate the JWT token
			claims := &model.Claims{}
			if err := ParseAndValidateToken(tokenString, claims); err != nil {
				log.Printf("Token validation error: %v", err)
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Check if the user's role is allowed
			if !isRoleAllowed(claims.Role, allowedRoles) {
				log.Printf("Access denied: insufficient permissions for role '%s'", claims.Role)
				http.Error(w, "Access denied: insufficient permissions", http.StatusForbidden)
				return
			}

			// Add userID and role to the request context
			ctx := context.WithValue(r.Context(), model.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, model.RoleKey, claims.Role)

			// Continue to the next handler
			log.Printf("Access granted: userID=%d, role=%s", claims.UserID, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isRoleAllowed checks if the given role is in the list of allowed roles
func isRoleAllowed(role string, allowedRoles []string) bool {
	for _, r := range allowedRoles {
		if r == role {
			return true
		}
	}
	return false
}
 
