package server

import (
	"net/http"
)

// Middleware to check user role
func (s *Server) RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the role from the context (this was set in getAuthCallBack)
			role, ok := r.Context().Value("userRole").(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the user has the required role
			if role != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// If the role matches, proceed with the request
			next.ServeHTTP(w, r)
		})
	}
}
