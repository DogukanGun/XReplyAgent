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
//	@host		localhost:3002
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
	"os"

	httpSwagger "github.com/swaggo/http-swagger"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/api/option"

	_ "cg-mentions-bot/api/docs"
	"cg-mentions-bot/api/handlers"
	"cg-mentions-bot/internal/utils/db"
)

func main() {
	ctx := context.Background()

	// Initialize Firebase
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("api/firebase-service-account.json"))
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v", err)
	}

	// Initialize MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default for local development
	}

	dbClient, err := db.ConnectToDB(mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Set database client for handlers
	handlers.SetDBClient(dbClient)

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

	r.Route("/api", func(r chi.Router) {
		r.Use(AuthMiddleware(authClient))
		r.Route("/user", func(r chi.Router) {
			r.Post("/check", handlers.CheckUserHandler)
			r.Post("/register", handlers.RegisterUserHandler)
			r.Post("/session", handlers.NewSessionHandler)
			r.Get("/profile", handlers.GetProfileHandler)
		})
		r.Post("/execute", handlers.ExecuteAppHandler)
		r.Post("/agent/ask", handlers.AgentAskHandler)
	})

	// Swagger endpoint (simple: require SWAGGER_URL at runtime)
	swaggerURL := os.Getenv("SWAGGER_URL")
	if swaggerURL == "" {
		log.Fatalf("SWAGGER_URL is required, e.g., http://<host>:3002/swagger/doc.json or https://domain/swagger/doc.json")
	}
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
	))

	log.Printf("Server starting on port 3002...")
	err = http.ListenAndServe(":3002", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
