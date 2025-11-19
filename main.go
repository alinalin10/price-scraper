package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

func main() {
	targetURL := "https://www.newegg.com/seagate-expansion-26tb-black-usb-3-0/p/N82E16822185116?Item=N82E16822185116&cm_sp=Homepage_SS-_-P2_22-185-116-_-11192025"
	selector := "div.price-current"
	fmt.Printf("attempting to GET %s\n", targetURL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Printf("error creating HTTP request: %s\n", err)
		return
	}
	
	// set User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	resp, err := client.Do(req)

	// resp, err := http.Get(targetURL)

	// if err != nil {
	// 	fmt.Printf("error making GET request: %s\n", err)
	// 	return
	// }

	// connection resources are released after the function finishes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status code: %d\n", resp.StatusCode)
	}

	// goquery document from http response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	
	if err != nil {
		fmt.Printf("error loading HTTP response body into goquery: %s\n", err)
		return
	}

	// use selector to find element
	selection := doc.Find(selector)
	
	if selection.Length() == 0 {
		fmt.Printf("no elements found for selector: %s\n", selector)
		return
	}

	// extract and print text content
	rawPriceText := selection.First().Text()

	// only num
	cleanedPriceText := strings.TrimSpace(rawPriceText)
	cleanedPriceText = strings.ReplaceAll(cleanedPriceText, "$", "")
	cleanedPriceText = strings.ReplaceAll(cleanedPriceText, ",", "")
	
	price, err := strconv.ParseFloat(cleanedPriceText, 64)
	if err != nil {
		fmt.Printf("error parsing price text to float: %s\n", err)
		return
	}
	
	fmt.Printf("extracted price: $%.2f\n", price)


	// // read response body (html content)
	// htmlBytes, err := io.ReadAll(resp.Body)

	// if err != nil {
	// 	fmt.Printf("error reading response body: %s\n", err)
	// 	return
	// }

	// htmlContent := string(htmlBytes)
	// fmt.Println("successfully fetched HTML content:")

	// snippetLength := 500
	// if len(htmlContent) > snippetLength {
	// 	fmt.Println(htmlContent[:snippetLength] + "...")
	// } else {
	// 	fmt.Println(htmlContent)
	// }
}