package handlers

import (
	"cg-mentions-bot/internal/services"
	"cg-mentions-bot/internal/utils/wallet"
	"encoding/json"
	"log"
	"net/http"
)

// NewSessionHandler handles device changes and returns user wallet data
//
//	@Summary		Create new session with updated device
//	@Description	Update device identifier for existing user and return wallet information
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		NewSessionRequest	true	"New Session Request"
//	@Success		200		{object}	RegisterUserResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		Bearer
//	@Router			/user/new-session [post]
func NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req NewSessionRequest
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

	// Get existing user
	user, userExists := GetUserByFirebaseID(firebaseID)
	if !userExists {
		log.Printf("User not found for Firebase ID: %s", firebaseID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// Update device identifier
	success := UpdateUserDeviceIdentifier(firebaseID, req.DeviceIdentifier)
	if !success {
		log.Printf("Failed to update device identifier for user: %s", firebaseID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update device"}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// Get or create wallet keys if they don't exist
	var walletKeys *wallet.WalletKeys
	if user.EthPublicKey != "" && user.SolanaPublicKey != "" {
		// User already has wallet keys, construct them from database
		walletKeys = &wallet.WalletKeys{
			EthWallet: wallet.WalletKeyPair{
				PublicAddress: user.EthPublicKey,
				PrivateKey:    user.EthPrivateKey,
			},
			SolanaWallet: wallet.WalletKeyPair{
				PublicAddress: user.SolanaPublicKey,
				PrivateKey:    user.SolanaPrivateKey,
			},
		}
	} else {
		// Create new wallet keys if they don't exist
		walletService := services.NewWalletService(DBClient)
		var err error
		walletKeys, err = walletService.CreateOrGetWallet(user.TwitterID)
		if err != nil {
			log.Printf("Failed to create wallets: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create wallets"}); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
			return
		}

		// Update user with wallet keys
		user.EthPublicKey = walletKeys.EthWallet.PublicAddress
		user.EthPrivateKey = walletKeys.EthWallet.PrivateKey
		user.SolanaPublicKey = walletKeys.SolanaWallet.PublicAddress
		user.SolanaPrivateKey = walletKeys.SolanaWallet.PrivateKey
		user.DeviceIdentifier = req.DeviceIdentifier

		success = CreateUserWithWallet(*user)
		if !success {
			log.Printf("Failed to update user with wallet keys")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user"}); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
			return
		}
	}

	// Return the same response format as register endpoint
	response := RegisterUserResponse{
		UID:       firebaseID,
		TwitterID: user.TwitterID,
		Username:  user.Username,
		Message:   "Session updated successfully",
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
	log.Printf("Successfully updated session for user - Firebase ID: %s, Twitter ID: %s, New Device: %s", firebaseID, user.TwitterID, req.DeviceIdentifier)
}
