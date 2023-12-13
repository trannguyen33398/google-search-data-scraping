package scrapingService

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	
	"github.com/PuerkitoBio/goquery"
)
type TScrapingData struct {
	Keyword string `json:"keyword"`
	TotalAdvertised int `json:"totalAdvertised"`
	TotalLink int `json:"totalLink"`
	TotalSearch string `json:"totalSearch"`
	Html string `json:"html"`
}
func ScrapingData(keyword string) TScrapingData {
	textFormat := strings.Replace(keyword," ","+",-1)
	url := fmt.Sprintf("https://www.google.com/search?q=%s", textFormat)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Read the response body
	//body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		log.Fatal(err) 
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	totalLink := 0
	totalAdvertised := 0
	var html string
	var totalSearch string
	doc.Find(".LHJvCe").Each(func(i int, result *goquery.Selection) {
		totalSearch =  result.Find("div").First().Text()
	})

	doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
		totalLink++

	})
	doc.Find(".v7W49e").Each(func(i int, result *goquery.Selection) {
		html,_ =result.Find("div").Html()

	})

	doc.Find(".U3A9Ac").Each(func(i int, result *goquery.Selection) {
		totalAdvertised++
	})

	
	return TScrapingData{keyword,totalAdvertised, totalLink, totalSearch, html}
}
