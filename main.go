package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

func ScrapePrice(targetURL string, selector string) (float64, error) {
	fmt.Printf("attempting to GET %s\n", targetURL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %s\n", err)
	}
	
	// set User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	resp, err := client.Do(req)


	// connection resources are released after the function finishes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status code: %d\n", resp.StatusCode)
	}

	// goquery document from http response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	
	if err != nil {
		return 0, fmt.Errorf("error loading HTTP response body into goquery: %s\n", err)
	}

	// use selector to find element
	selection := doc.Find(selector)
	
	if selection.Length() == 0 {
		return 0, fmt.Errorf("no elements found for selector: %s", selector)
	}

	// extract and print text content
	rawPriceText := selection.First().Text()

	// only num
	cleanedPriceText := strings.TrimSpace(rawPriceText)
	cleanedPriceText = strings.ReplaceAll(cleanedPriceText, "$", "")
	cleanedPriceText = strings.ReplaceAll(cleanedPriceText, ",", "")
	
	price, err := strconv.ParseFloat(cleanedPriceText, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing price text to float: %s", err)
	}
	return price, nil
}

func main() {
	targetURL := "https://www.newegg.com/seagate-expansion-26tb-black-usb-3-0/p/N82E16822185116?Item=N82E16822185116&cm_sp=Homepage_SS-_-P2_22-185-116-_-11192025"
	selector := "div.price-current"
	price, err := ScrapePrice(targetURL, selector)
	if err != nil {
		fmt.Printf("Error scraping price: %s\n", err)
		return
	}
	fmt.Printf("The scraped price is: $%.2f\n", price)
	

}