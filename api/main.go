// Package main XReplyAgent API
//
//	@title			XReplyAgent API
//	@version		1.0
//	@description	API for XReplyAgent with Firebase authentication
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:3000
//	@BasePath	/api
//
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/api/option"

	"cg-mentions-bot/api/handlers"
	_ "cg-mentions-bot/docs"
)

func main() {
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("api/firebase-service-account.json"))
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"message": "Server is running"}`); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Use(AuthMiddleware(authClient))
		r.Post("/execute-app", handlers.ExecuteAppHandler)
		r.Get("/profile", handlers.GetProfileHandler)
	})

	log.Printf("Server starting on port 3000...")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
