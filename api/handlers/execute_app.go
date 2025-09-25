package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ExecuteAppHandler executes an application
//
//	@Summary		Execute an applicat
//	@Summary		Execute an application
//	@Description	Execute an application with optional parameters
//	@Tags			apps
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ExecuteAppRequest	true	"Execute App Request"
//	@Success		200		{object}	ExecuteAppResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		Bearer
//	@Router			/execute-app [post]
func ExecuteAppHandler(w http.ResponseWriter, r *http.Request) {
	var req ExecuteAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON in request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON in request body"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	if req.AppName == "" {
		log.Printf("Missing app_name in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing app_name field"}); err != nil {
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

	response := ExecuteAppResponse{
		UID:     uid,
		AppName: req.AppName,
		Status:  "executed",
		Message: "App executed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("Successfully executed app %s for user %s", req.AppName, uid)
}
