// ...existing code...
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getTokens() map[string]string {
	ports := []string{"78", "79", "80", "81", "82", "22", "789"}
	results := map[string]string{}
	for _, port := range ports {
		token := getToken(port)
		results[port] = token
	}
	return results
}

func main() {
	t := flag.String("t", "token", "type to get: token/scan")
	port := flag.String("port", "22", "port number of the OIDC server")
	flag.Parse()
	if *t == "token" {
		fmt.Println(getToken(*port))
	} else {
		apiServerAddr := os.Getenv("APISERVER")
		if apiServerAddr == "" {
			apiServerAddr = "https://192.168.204.131:6443"
		}
		userToken := os.Getenv("TOKEN")
		if userToken == "" {
			userToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkNFRTdfd0pEY2F2Qkpvb2Nxal8wdG5yVEUtcWJSRzNsVVVZTWhTczFsWk0ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzY2MDU2MjQ4LCJpYXQiOjE3NjYwNTI2NDgsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiMTU5NDJlMzMtZGYzMS00MGQ4LWIxYWQtMWJkODk4NmI0MTQ5Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImRlZmF1bHQiLCJ1aWQiOiJjYjYzOTVkNi1kODQ5LTQyZGItODY3MS04YjJjOWU5NmQyZjUifX0sIm5iZiI6MTc2NjA1MjY0OCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVmYXVsdCJ9.PlKTvKFLE7qTFAuVhq3jCeXhcJLFDVetWowfVdpbFkMOamuGe1WCMcC2mz14100zDpYn0mf1XxkzkcYpk0A2MnKRxJNjE65dDEg0NP6waPQR3_d7qXKCLhqJcoChiCFw4s4ilpqRUyv19v0BsWkz7V6VrS64thp4NJJgK3S9xcpq3-mzKFiVYhROPP9XulklSxs0bGL9mxThVM65wDV04HPj-luyiR0QBA8l3SOLLHXWKX-pQ89ZEgqN1MGoLh52MTutedbzwDBm2NLMlwhjVx4rwnzk7Dl_Uv0wz3L-QNVLtlLN7DNgcexC13sNQvErTeMO78efoWk22o4tgAQY6A"
		}
		config := &rest.Config{
			Host:        apiServerAddr,
			BearerToken: userToken,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}
		clientset, _ := kubernetes.NewForConfig(config)
		tokens := getTokens()
		for port, token := range tokens {
			tr := &authenticationv1.TokenReview{
				Spec: authenticationv1.TokenReviewSpec{
					Token: token,
				},
			}
			result, err := clientset.AuthenticationV1().TokenReviews().Create(
				context.TODO(),
				tr,
				metav1.CreateOptions{},
			)
			if err != nil {
				panic(err)
			}
			if !strings.Contains(result.Status.Error, "connection refused") {
				fmt.Println(port + "[√]")
				fmt.Println("port msg:", result.Status.Error)
				fmt.Println("====================")
			}
		}
	}

}
