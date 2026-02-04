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

	results := map[string]string{}
	ports := []string{"10278", "11278", "12278", "13278", "14278", "15278", "16278", "17278", "18278", "19278", "20000"}
	for _, port := range ports {
		results[port] = getToken(port)
	}
	// startPort := 10278
	// endPort := 20000
	// for port := startPort; port <= endPort; port++ {
	// 	portStr := strconv.Itoa(port)
	// 	token := getToken(portStr)
	// 	results[portStr] = token
	// }
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
			apiServerAddr = ""
		}
		userToken := os.Getenv("TOKEN")
		if userToken == "" {
			userToken = "eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjIl0sImV4cCI6MTc3MDIxNDQ1MCwiaWF0IjoxNzcwMjEwODUwLCJpc3MiOiJodHRwczovL29pZGMuZWtzLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL2lkLzBGODNCMkVGMUI3QjFBNzcxNTdGRTI0RjMzODNFOUQ2IiwianRpIjoiMmY5MGZjNTctZDVlYi00M2M0LTk4NzMtM2Y2MzRiMzQ1YTJmIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImRlZmF1bHQiLCJ1aWQiOiI4MTQwYmYwNC0zYTFlLTQ1ZjUtODJjZi0xYWQ5OGE4ZWYyNzcifX0sIm5iZiI6MTc3MDIxMDg1MCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVmYXVsdCJ9.Tw-sZbe42rP_rYXvqWeToZGuMvaGvFp-h9sRvqs8SGA9GBP5UrFtSdAF-E1HjEizKknZKCBa5rDVb0XJrd286kCkfMGzwRE_r_auvMdcCLPotDrDNUVr0NUfVo_akwC9EHBeRipXyabCIRlBHtUbtEtBTGEYeeZaouioyunFQmx48ZjK0FiCh63evJyEkB6RxGHGCRD3ldaSc1tfY3ek6mEHs7D5IUfTk73Y7JxSIfx6QszzM6-rd_VfRzzmcVBNfljaSu0wu9k3YyX7KO3swuLsoofkPE4GJbzzNKP-KPcV8A-9eBqXN5Nrgb5qhrZ89lebEhBpjGaOy2tlNbQJ0A"
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
			//fmt.Println(result)
			if !strings.Contains(result.Status.Error, "connection refused") {
				fmt.Println(port + "[√]")
				fmt.Println("port msg:", result.Status.Error)
				fmt.Println("====================")
			}
		}
	}

}
