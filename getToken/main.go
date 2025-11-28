// ...existing code...
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func main() {
	keyPath := flag.String("key", "oidc-server-key.pem", "path to RSA private key (PEM)")
	num := flag.Int("n", 1, "number of tokens to generate")
	concurrency := flag.Int("c", 10, "concurrency level")
	flag.Parse()

	pemBytes, err := ioutil.ReadFile(*keyPath)
	if err != nil {
		log.Fatalf("read key: %v", err)
	}
	privateKey, err := parsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key: %v", err)
	}

	var wg sync.WaitGroup
	tokens := make(chan string, *num)
	sem := make(chan struct{}, *concurrency)

	for i := 0; i < *num; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			now := time.Now()
			claims := jwt.MapClaims{
				"iss": "https://192.168.204.1:8080",
				"aud": "testid",
				"exp": now.Add(time.Hour).Unix(),
				"_claim_names": map[string]interface{}{
					"testgroup": "src1",
				},
				"_claim_sources": map[string]interface{}{
					"src1": map[string]interface{}{
						"endpoint": "http://localhost:2334/claims?user=user1",
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
				log.Printf("sign error (i=%d): %v", i, err)
				return
			}
			tokens <- signed
		}(i)
	}

	go func() {
		wg.Wait()
		close(tokens)
	}()

	for t := range tokens {
		fmt.Println(t)
	}
}
