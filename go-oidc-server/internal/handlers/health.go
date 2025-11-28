package handlers

import (
    "net/http"
)

// HealthCheck responds with a simple message indicating the server is healthy.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}