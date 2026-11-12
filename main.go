package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	rtclisten "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder"
)

type TokenRequest struct {
	ChannelName     string `json:"channelName"`
	UID             uint32 `json:"uid"`
	Role            uint32 `json:"role"`
	ExpirationInSec uint32 `json:"expirationInSec"`
}

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appID := os.Getenv("AGORA_APP_ID")
	appCertificate := os.Getenv("AGORA_APP_CERTIFICATE")

	if appID == "" || appCertificate == "" {
		http.Error(w, "Environment variables not set", http.StatusInternalServerError)
		return
	}

	expireTimestamp := uint32(3600)
	if req.ExpirationInSec > 0 {
		expireTimestamp = req.ExpirationInSec
	}

	result, err := rtclisten.BuildTokenWithUid(
		appID,
		appCertificate,
		req.ChannelName,
		req.UID,
		rtclisten.RoleSubscriber,
		expireTimestamp,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w.Encode(map[string]string{"token": result}))
}

func main() {
	http.HandleFunc("/token", generateTokenHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Agora Token Server is running!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
