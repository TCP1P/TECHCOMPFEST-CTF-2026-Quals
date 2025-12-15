package middleware

import (
	"app-api/internal/utils"
	"net/http"

	"go.uber.org/zap"
)

func AdminMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the user is an admin
			if !isAdmin(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			userId, _ := utils.GetUserIDFromContext(r.Context())
			logger.Info("Admin access granted", zap.Uint("userID", userId))
			next.ServeHTTP(w, r)
		})
	}
}

// isAdmin checks if the user has admin privileges
func isAdmin(r *http.Request) bool {
	role, err := utils.GetUserRoleFromContext(r.Context())
	if err != nil {
		return false
	}
	return role == "admin"
}
