package handlers

import (
	"cg-mentions-bot/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

// RegisterUserHandler registers a new user
//
//	@Summary		Register new user
//	@Description	Register a new user with Twitter ID and username
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterUserRequest	true	"Register User Request"
//	@Success		200		{object}	RegisterUserResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		Bearer
//	@Router			/user/register [post]
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON in request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON in request body"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	if req.TwitterID == "" {
		log.Printf("Missing twitter_id in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing twitter_id field"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	if req.Username == "" {
		log.Printf("Missing username in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing username field"}); err != nil {
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

	firebaseID, ok := r.Context().Value(UidKey).(string)
	if !ok {
		log.Printf("Firebase ID not found in request context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// Create or get wallets for the user
	walletService := services.NewWalletService(DBClient)
	walletKeys, err := walletService.CreateOrGetWallet(req.TwitterID)
	if err != nil {
		log.Printf("Failed to create wallets: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create wallets"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// Save user to database with wallet keys
	user := User{
		FirebaseID:       firebaseID,
		TwitterID:        req.TwitterID,
		Username:         req.Username,
		DeviceIdentifier: req.DeviceIdentifier,
		EthPublicKey:     walletKeys.EthWallet.PublicAddress,
		EthPrivateKey:    walletKeys.EthWallet.PrivateKey,
		SolanaPublicKey:  walletKeys.SolanaWallet.PublicAddress,
		SolanaPrivateKey: walletKeys.SolanaWallet.PrivateKey,
	}

	success := CreateUserWithWallet(user)
	if !success {
		log.Printf("Failed to save user to database")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to register user"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	response := RegisterUserResponse{
		UID:       firebaseID,
		TwitterID: req.TwitterID,
		Username:  req.Username,
		Message:   "User registered successfully",
		Wallets: WalletKeys{
			EthWallet: WalletKeyPair{
				PublicAddress: walletKeys.EthWallet.PublicAddress,
				PrivateKey:    walletKeys.EthWallet.PrivateKey,
			},
			SolanaWallet: WalletKeyPair{
				PublicAddress: walletKeys.SolanaWallet.PublicAddress,
				PrivateKey:    walletKeys.SolanaWallet.PrivateKey,
			},
		},
		// Backward compatibility
		PrivateKey: walletKeys.EthWallet.PrivateKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("Successfully registered user - Firebase ID: %s, Twitter ID: %s, Username: %s", firebaseID, req.TwitterID, req.Username)
}
