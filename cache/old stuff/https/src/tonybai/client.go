package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{Transport: tr}
	resp, err := http.Get("https://localhost:8081")
	// resp, err := http.Get("https://Frankie:8081")

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
