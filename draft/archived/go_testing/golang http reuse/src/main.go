package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

const (
	MaxIdleConnections int = 10
	RequestTimeout     int = 5
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

func main() {
	var endPoint string = "http://localhost:8087/"
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", endPoint, nil)
		response, _ := httpClient.Do(req)
		// MUST read all response's data
		ioutil.ReadAll(response.Body)
		// Close the connection to reuse it
		response.Body.Close()
	}
}
