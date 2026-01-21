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
			userToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImJuTkRVcmRsdzJ6QnpsVG9zV0RvOVgybFJVd3NWUkRyaWJ5Tkp1clFvc3cifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzY5MDAxNzUwLCJpYXQiOjE3Njg5OTgxNTAsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNTA0MjliYTctNzcyOC00NjE5LTk1YTgtMGU3MWMyYzQyNzI2Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImRlZmF1bHQiLCJ1aWQiOiI1YjQwMjRjNS0yNTgyLTQwM2UtYTIyZi0wNDNlOTJkYzJlNzQifX0sIm5iZiI6MTc2ODk5ODE1MCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVmYXVsdCJ9.bxMh-hjqMVhWSzCMelgobAZwxoxKUuUbQUt6wFwSJTVL2enXm8yHEuAw4P79iV4NvMJPCFS-kys0uH6gXefABTuCS-9vYG_1vHvVMXvT3t4QK5gzISIjN38VA3mcKN28-UjfoAVe_I0FWa4JBckJoAZXO3wj77CvnXQs_j7VHTE-GT4Vc2yg7T0YzNlEk2LCf_bqzrmt1XvnaNjsWAvwvJF9XvWGFOysYKvYz9H5DgcSyAjar_KVU59KFr1enA_tGXGDdVfPMv2mkFFwRnq5LShkL1e2ZfDUJ6g2KZ2KfHEMyQKhysifArm8kmSgpRX_7LmvpC9jXteGD0jxUVbPYQ"
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
