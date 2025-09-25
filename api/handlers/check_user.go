package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// CheckUserHandler checks if user exists in database
//
//	@Summary		Check if user exists
//	@Description	Check if the authenticated user has registered before
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CheckUserResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		Bearer
//	@Router			/user/check [get]
func CheckUserHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if user exists in database using Firebase ID
	userExists := CheckUserExists(uid)

	response := CheckUserResponse{
		Register: !userExists, // true if user needs to register
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("User existence check for UID %s: register=%t", uid, !userExists)
}
