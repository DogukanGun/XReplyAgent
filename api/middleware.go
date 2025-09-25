package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"cg-mentions-bot/api/handlers"
	"firebase.google.com/go/v4/auth"
)

func AuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("Missing Authorization header")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "Missing Authorization header"}); err != nil {
					log.Printf("Error encoding response: %v", err)
				}
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Printf("Invalid Authorization header format")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "Invalid Authorization header format"}); err != nil {
					log.Printf("Error encoding response: %v", err)
				}
				return
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			if idToken == "" {
				log.Printf("Empty ID token")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "Empty ID token"}); err != nil {
					log.Printf("Error encoding response: %v", err)
				}
				return
			}

			token, err := authClient.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				log.Printf("Invalid Firebase ID token: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "Invalid Firebase ID token"}); err != nil {
					log.Printf("Error encoding response: %v", err)
				}
				return
			}

			ctx := context.WithValue(r.Context(), handlers.UidKey, token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
