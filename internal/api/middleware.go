package api

import (
	"context"
	"net/http"

	"health-api/internal/models"
)

type contextKey string

const (
	roleContextKey   contextKey = "role"
	userIDContextKey contextKey = "user_id"
)

func RoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleHeader := r.Header.Get("X-Role")
		if roleHeader == "" {
			http.Error(w, "Missing X-Role header", http.StatusUnauthorized)
			return
		}

		var role models.Role
		switch roleHeader {
		case "physician":
			role = models.RolePhysician
		case "patient":
			role = models.RolePatient
		case "admin":
			role = models.RoleAdmin
		default:
			http.Error(w, "Invalid role. Must be one of: physician, patient, admin", http.StatusBadRequest)
			return
		}

		userIDHeader := r.Header.Get("X-User-ID")
		if userIDHeader == "" && role != models.RoleAdmin {
			http.Error(w, "Missing X-User-ID header", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), roleContextKey, role)
		if userIDHeader != "" {
			ctx = context.WithValue(ctx, userIDContextKey, userIDHeader)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRoleFromContext(ctx context.Context) (models.Role, bool) {
	role, ok := ctx.Value(roleContextKey).(models.Role)
	return role, ok
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}