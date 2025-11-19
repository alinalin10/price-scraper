package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"time"
	"os"
)

func ScrapePrice(targetURL string, selector string) (float64, error) {
	fmt.Printf("attempting to GET %s\n", targetURL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %s\n", err)
	}
	
	// set User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
        return 0, fmt.Errorf("error executing request: %w", err)
    }

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
	rawPriceText, exists := selection.First().Attr("data-value")
	if !exists {
		return 0, fmt.Errorf("data-value attribute not found in the selected element")
	}

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

func LogPrice(price float64) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	dataLine := fmt.Sprintf("%s,%.2f\n", timestamp, price)

	file, err := os.OpenFile("prices.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %s", err)
	}
	defer file.Close()

	if _, err := file.WriteString(dataLine); err != nil {
		return fmt.Errorf("error writing to log file: %s", err)
	}
	return nil

}

func main() {
	targetURL := "https://finance.yahoo.com/markets/crypto/all/"
	selector := `fin-streamer[data-symbol="BTC-USD"]`
	//scrape
	price, err := ScrapePrice(targetURL, selector)
	if err != nil {
		fmt.Printf("Error scraping price: %s\n", err)
		return
	}
	fmt.Printf("The scraped price is: $%.2f\n", price)

	//log
	if err := LogPrice(price); err != nil {
		fmt.Printf("Error logging price: %s\n", err)
		return
	}
	fmt.Println("Price logged successfully to prices.csv")
	

}