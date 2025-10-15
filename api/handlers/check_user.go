package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// CheckUserHandler checks if user exists in database and validates device
//
//	@Summary		Check if user exists and device matches
//	@Description	Check if the authenticated user has registered before and if device has changed
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CheckUserRequest	true	"Check User Request"
//	@Success		200		{object}	CheckUserResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		Bearer
//	@Router			/user/check [post]
func CheckUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON in request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON in request body"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	if req.DeviceIdentifier == "" {
		log.Printf("Missing device_identifier in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing device_identifier field"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

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
	user, userExists := GetUserByFirebaseID(uid)
	deviceChanged := false

	if userExists && user.DeviceIdentifier != "" && user.DeviceIdentifier != req.DeviceIdentifier {
		deviceChanged = true
	}

	response := CheckUserResponse{
		Register:      !userExists, // true if user needs to register
		DeviceChanged: deviceChanged,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("User existence check for UID %s: register=%t, device_changed=%t", uid, !userExists, deviceChanged)
}
