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
			apiServerAddr = "https://5A4BA0B59FA3A5FAA15C483FC3A5D6B5.gr7.us-east-1.eks.amazonaws.com"
		}
		userToken := os.Getenv("TOKEN")
		if userToken == "" {
			userToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImI1MmRhYzcwOWY3NDFkZjBjMjczZjZmMzg1N2QxM2I0YmZjNGEyMGEifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjIl0sImV4cCI6MTc2OTU3NDc4NiwiaWF0IjoxNzY5NTcxMTg2LCJpc3MiOiJodHRwczovL29pZGMuZWtzLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL2lkLzVBNEJBMEI1OUZBM0E1RkFBMTVDNDgzRkMzQTVENkI1IiwianRpIjoiYzdmZWExOTMtNzAzYS00YTE1LTkzZjEtNDUzNjVmOWJhMGQ3Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImFkbWluIiwidWlkIjoiMDEwZDFkNDktZjhmZi00MGU3LTk3OGMtNmU1MjY1NWZiOTBkIn19LCJuYmYiOjE3Njk1NzExODYsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmFkbWluIn0.ee0o0GOIKC8LQZkIEFzjKAYaXFbSZthsR7oq1DR_H0u_STzoALLqRAS4oaYg5BCNmGcpaPyy5nmxPMj2f3YDNlcrZwR2GmG7FqNGovxXMkaHa2j8avJfyNzQ5glntw3uBROscBNtsl00pu2f9jPqEpsi_ONBqvCNJi7mEjgPQpEfbXVZVMOXsCEwekUrawJcd-buv8pa3T6d_F8dNr68dSUtLoygMhqRXG_EGdB3fR_WT9L7KiibOzGF-Ac-gN4uhPwpcuGgQNak5hl_v-uaXcz-WliWfJ7mDdKJ2OVpc-J1epJz_fuN85RdX3BW1axNJkGMvy32QCwcc2uh_KQmEA"
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
			fmt.Println(result)
			if !strings.Contains(result.Status.Error, "connection refused") {
				fmt.Println(port + "[√]")
				fmt.Println("port msg:", result.Status.Error)
				fmt.Println("====================")
			}
		}
	}

}
