package webapi

import (
	"context"
	"log/slog"
	"net/http"
	"slices"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/domain/user"
)

func AclMiddleware(requiredRoles ...user.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value(ctxUserRole).(user.Role)
			if !slices.Contains(requiredRoles, role) {
				slog.Warn("ACL check failed", "userRole", role, "requiredRoles", requiredRoles)
				ErrorResponse(w, http.StatusForbidden, "forbidden")
				return
			}
			slog.Info("ACL check passed", "userRole", role)
			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(srvc *app.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sessionID string
			if c, err := r.Cookie("session-id"); err == nil {
				sessionID = c.Value
			} else {
				slog.Warn("No session cookie found", "err", err)
				ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			user, err := srvc.GetBySessionID(sessionID)
			if err != nil {
				slog.Warn("Invalid session ID", "err", err)
				ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			ctx := context.WithValue(r.Context(), ctxUserID, user.UUID.String())
			ctx = context.WithValue(ctx, ctxSessionID, sessionID)
			ctx = context.WithValue(ctx, ctxUserRole, user.Role)

			slog.Info("Authenticated user", "userID", user.UUID.String(), "role", user.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
