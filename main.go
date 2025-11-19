package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	targetURL := "http://google.com"

	fmt.Printf("attempting to GET %s\n", targetURL)
	resp, err := http.Get(targetURL)

	if err != nil {
		fmt.Printf("error making GET request: %s\n", err)
		return
	}

	// connection resources are released after the function finishes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status code: %d\n", resp.StatusCode)
	}

	// read response body (html content)
	htmlBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
		return
	}

	htmlContent := string(htmlBytes)
	fmt.Println("successfully fetched HTML content:")

	snippetLength := 500
	if len(htmlContent) > snippetLength {
		fmt.Println(htmlContent[:snippetLength] + "...")
	} else {
		fmt.Println(htmlContent)
	}
}