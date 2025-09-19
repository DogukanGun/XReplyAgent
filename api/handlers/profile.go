package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// GetProfileHandler retrieves user profile
//
//	@Summary		Get user profile
//	@Description	Get the authenticated user's profile information
//	@Tags			profile
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ProfileResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		Bearer
//	@Router			/profile [get]
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(UidKey).(string)
	if !ok {
		log.Printf("UID not found in request context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	response := ProfileResponse{
		UID:       uid,
		Email:     "user@example.com",
		TwitterID: "@user_twitter",
		Message:   "Profile retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("Successfully retrieved profile for user %s", uid)
}
