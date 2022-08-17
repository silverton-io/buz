package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	protocol := flag.String("protocol", "http", "The protocol to poll for healthcheck")
	host := flag.String("host", "localhost", "The host to poll for healthcheck")
	port := flag.String("port", "8080", "The port to poll for healthcheck")
	path := flag.String("path", "/health", "The path to poll for healthcheck")
	flag.Parse()
	healthcheckPath := *protocol + "://" + *host + ":" + *port + *path
	fmt.Println("Checking service health via: " + healthcheckPath)
	resp, err := http.Get(healthcheckPath)
	if err != nil {
		fmt.Println("error when getting response: " + err.Error())
		os.Exit(1)
	}
	if resp == nil {
		fmt.Println("empty response")
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("non-200 status code: " + fmt.Sprint(resp.StatusCode))
		os.Exit(1)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("could not read response body")
	}
	fmt.Println(string(b))
	os.Exit(0)
}
