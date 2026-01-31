package main

import (
	"go-oidc-server/internal/server"
	"go-oidc-server/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	utils.GetToken()
	r := mux.NewRouter()

	// Set up routes
	server.SetupRoutes(r)

	// Start the HTTPS server (expects cert.pem and key.pem in project root)
	log.Println("Starting OIDC server on https://localhost:30080")
	if err := http.ListenAndServeTLS(":30080", "cert.pem", "key.pem", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
