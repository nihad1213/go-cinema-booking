package handler

import (
	"context"
	"net/http"

	"cinema-booking/internal/domain"
)

type ctxKey string

const userKey ctxKey = "user"

func UserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), userKey, r.Header.Get("X-User-ID"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func currentUser(r *http.Request) string {
	u, _ := r.Context().Value(userKey).(string)
	return u
}

func requireUser(r *http.Request) (string, error) {
	u := currentUser(r)
	if u == "" {
		return "", domain.ErrUnauthorized
	}
	return u, nil
}
