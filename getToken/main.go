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
			apiServerAddr = "https://136.119.37.168:6443"
		}
		userToken := os.Getenv("TOKEN")
		if userToken == "" {
			userToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkZHT3NLTXF4Z0EwMmNqNE1wSmtKck03RVU1Smo1V2hrc2VPV1VpMmxIR2MifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzY4OTk3ODMyLCJpYXQiOjE3Njg5OTQyMzIsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNjk5ZjFjYjYtNmM0My00YTM2LTgwYjMtZDU4ZGIwODgyOWY3Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImRlZmF1bHQiLCJ1aWQiOiI3YWJlMWFkOC1mNWNkLTQ1NmEtODFiZS1jYTZkYmYwNjRkMTEifX0sIm5iZiI6MTc2ODk5NDIzMiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVmYXVsdCJ9.Y0T2rsW4rZt9mqYaB_J5ZxYo2vHqDYr1r_je24C1JZ5ZrWFB_u6UPQ0MZzmrPzi4L4SZwMbdTI0qmgyKkJdEqtfqhXYIID9suCl0uy2bdak5tlQ1XkC3S2VOl4NnaugOSPH9H8DX3KUVMfHDmuukuiAbtxE2oBHPqjOF1O0eBPkhsIJB74OTHbAsa-Do8D7B-qSFB406-Xu2ObvQvs864N8KbaG6i5dgoEi0XpFbfhFtr5gTji4SmqDLa9GRifVPo0JKPjmQYTr6KMqeOxikWcPUB47AnYfbzEk6f1x3D233oOEx48KKrCkE-S-J7jNkvK5CVr-5GvhihwJEJWTB1A"
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
