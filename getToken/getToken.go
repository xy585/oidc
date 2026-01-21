// ...existing code...
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	keyPath = "oidc-server-key.pem"
)

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	// try PKCS1
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	// try PKCS8
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	rk, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}
	return rk, nil
}

func getToken(port string) string {
	pemBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("read key: %v", err)
	}
	privateKey, err := parsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key: %v", err)
	}
	result := ""
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": "https://34.59.255.46:30080",
		"aud": "testid",
		"exp": now.Add(time.Hour).Unix(),
		"_claim_names": map[string]interface{}{
			"testgroup": "src1",
		},
		"_claim_sources": map[string]interface{}{
			"src1": map[string]interface{}{
				"endpoint": "http://localhost:" + port + "/test?a=a",
			},
		},
		//"testgroup": "a",
		// "sub": fmt.Sprintf("user-%d", i),
		// "iat": now.Unix(),
		// "jti": fmt.Sprintf("%d-%d", now.UnixNano(), i),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		result = "err"
	}
	result = signed

	return result
}
