package config

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// SessionMiddleware ensures the user is authenticated
func SessionMiddleware(next http.HandlerFunc, sessionManager *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionUsername := sessionManager.GetString(r.Context(), "username")
		if sessionUsername == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// set the userID in the request context so that it can be used to create TODO items.
		sessionUserID := sessionManager.Get(r.Context(), "userID")
		ctx := context.WithValue(r.Context(), "userID", sessionUserID)
		next(w, r.WithContext(ctx))
	}
}
