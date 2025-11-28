package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	// rsaPriv is generated at startup for demo; only public part is exposed via JWKS.
	rsaPriv *rsa.PrivateKey
	// jwksResp holds the computed JWKS response
	jwksResp map[string]interface{}
)

func init() {
	// Try to load a PEM private key file named "key.pem" from working dir.
	// If not present or fails to parse, fall back to generating an ephemeral RSA key.
	var err error
	rsaPriv, err = loadRSAPrivateKeyFromPEM("key.pem")
	if err != nil {
		// generate ephemeral key
		rsaPriv, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			// If key generation fails, fall back to an empty JWKS to avoid panics.
			jwksResp = map[string]interface{}{"keys": []interface{}{}}
			return
		}
	}

	pub := &rsaPriv.PublicKey

	// modulus (n) and exponent (e) must be base64url-encoded without padding.
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	eBytes := big.NewInt(int64(pub.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	// Stable key id (kid) derived from public key (SHA-256 of modulus).
	kidBytes := sha256.Sum256(pub.N.Bytes())
	kid := base64.RawURLEncoding.EncodeToString(kidBytes[:])

	jwk := map[string]interface{}{
		"kty": "RSA",
		"kid": kid,
		"use": "sig",
		"alg": "RS256",
		"n":   n,
		"e":   e,
	}

	jwksResp = map[string]interface{}{
		"keys": []interface{}{jwk},
	}
}

func loadRSAPrivateKeyFromPEM(path string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var block *pem.Block
	block, _ = pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}
	// try PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	// try PKCS#8
	if parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if key, ok := parsed.(*rsa.PrivateKey); ok {
			return key, nil
		}
		return nil, fmt.Errorf("pkcs8 key is not RSA")
	}
	return nil, fmt.Errorf("failed to parse private key")
}

func SetupRoutes(r *mux.Router) {
	// OIDC discovery
	r.HandleFunc("/.well-known/openid-configuration", discoveryHandler).Methods("GET")

	// JWKS
	r.HandleFunc("/.well-known/jwks.json", jwksHandler).Methods("GET")
	r.HandleFunc("/jwks", jwksHandler).Methods("GET") // alternate path

	// Authorization endpoint (very simple demo)
	r.HandleFunc("/authorize", authorizeHandler).Methods("GET")

	// Token endpoint (very simple demo)
	r.HandleFunc("/token", tokenHandler).Methods("POST", "OPTIONS")

	// Userinfo
	r.HandleFunc("/userinfo", userinfoHandler).Methods("GET", "POST")

	// Health
	r.HandleFunc("/health", healthHandler).Methods("GET", "HEAD")

	// Simple CORS for browser testing
	r.Use(corsMiddleware)
}

func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	// derive scheme/host from request (supports X-Forwarded-Proto)
	scheme := "https"
	if r.TLS == nil {
		if p := r.Header.Get("X-Forwarded-Proto"); p != "" {
			scheme = p
		} else {
			scheme = "http"
		}
	}
	host := r.Host
	if host == "" {
		host = "localhost:8443"
	}
	base := scheme + "://" + host

	resp := map[string]interface{}{
		"issuer":                                base,
		"authorization_endpoint":                base + "/authorize",
		"token_endpoint":                        base + "/token",
		"userinfo_endpoint":                     base + "/userinfo",
		"jwks_uri":                              base + "/.well-known/jwks.json",
		"response_types_supported":              []string{"code", "token", "id_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}
	writeJSON(w, http.StatusOK, resp)
}

//	func jwksHandler(w http.ResponseWriter, r *http.Request) {
//		// Minimal placeholder JWKS (no real keys). Replace with real keys in production.
//		resp := map[string]interface{}{
//			"keys": nil,
//		}
//		writeJSON(w, http.StatusOK, resp)
//	}
func jwksHandler(w http.ResponseWriter, r *http.Request) {
	// Return the precomputed JWKS (public key only).
	if jwksResp == nil {
		// Fallback empty keys array
		writeJSON(w, http.StatusOK, map[string]interface{}{"keys": []interface{}{}})
		return
	}
	writeJSON(w, http.StatusOK, jwksResp)
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	// Very simple demo: build a redirect with a fake code and preserve state if provided.
	q := r.URL.Query()
	redirect := q.Get("redirect_uri")
	state := q.Get("state")
	if redirect == "" {
		http.Error(w, "missing redirect_uri", http.StatusBadRequest)
		return
	}
	// NOTE: In real server validate client_id, redirect_uri, scopes, user consent, etc.
	code := "fake-code-" + time.Now().Format("20060102150405")
	u := redirect + "?code=" + code
	if state != "" {
		u += "&state=" + state
	}
	http.Redirect(w, r, u, http.StatusFound)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// Very simple demo token response. Accepts form params grant_type, code, client_id, client_secret, etc.
	// NOTE: In production validate client auth, grant, and issue real tokens.
	resp := map[string]interface{}{
		"access_token": "fake-access-token",
		"token_type":   "Bearer",
		"expires_in":   3600,
		"id_token":     "fake-id-token",
	}
	writeJSON(w, http.StatusOK, resp)
}

func userinfoHandler(w http.ResponseWriter, r *http.Request) {
	// Demo user info. In production validate access token and return real claims.
	resp := map[string]interface{}{
		"sub":   "user-1234",
		"name":  "Demo User",
		"email": "demo@example.com",
	}
	writeJSON(w, http.StatusOK, resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// permissive for local testing
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
