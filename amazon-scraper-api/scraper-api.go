package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type dataExtractedMap struct {
	Name         string `json:"name,omitempty"`
	ImageURL     string `json:"imageURL,omitempty"`
	Desc         string `json:"description,omitempty"`
	Price        string `json:"price,omitempty"`
	TotalReviews string `json:"totalReviews,omitempty"`
}

type respMapObj struct {
	URL     string           `json:"url,omitempty"`
	Product dataExtractedMap `json:"product,omitempty"`
}

type StatusObject struct {
	InsertedID    string `json:"InsertedID,omitempty"`
	MatchedCount  int    `json:"MatchedCount,omitempty"`
	ModifiedCount int    `json:"ModifiedCount,omitempty"`
}

func scraper(url string) respMapObj {

	respMap := respMapObj{}

	respMap.URL = url

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "www.amazon.com/*",
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnHTML("#ppd", func(e *colly.HTMLElement) {
		var productName, imgUrl, stars, desc, price string

		productName = e.ChildText("#title")

		desc = "Amazon Product"

		stars = e.ChildText("span.a-icon-alt")
		FormatStars(&stars)

		price = e.ChildText("span.a-size-medium.a-color-price.priceBlockBuyingPriceString")
		FormatPrice(&price)

		imgUrl = e.ChildAttr("img.a-dynamic-image", "src")

		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			if i != 0 {
				desc += strings.TrimSpace(e.ChildText("span.a-list-item")) + ". "
			}
		})

		dataExtracted := dataExtractedMap{
			Name:         productName,
			ImageURL:     imgUrl,
			Desc:         desc,
			Price:        price,
			TotalReviews: stars,
		}

		respMap.Product = dataExtracted
	})

	c.Visit(url)

	return respMap
}

func getFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is REST-API-in-Go!"+
		"\nPlease do POST Request and pass product url in body to the API for scrapping Amazon Product Details.")
}

func postFunc(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	form_data := respMapObj{}
	err := decoder.Decode(&form_data)
	if err != nil {
		panic(err)
	}

	form_data = scraper(form_data.URL)

	product_details, err := json.Marshal(form_data)
	if err != nil {
		log.Fatal("json.Marshal failed due to the error:", err)
	}

	collector_url := "http://backend:3031/collector"
	requestObject, err := http.NewRequest("POST", collector_url, bytes.NewBuffer(product_details))
	requestObject.Header.Set("content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(requestObject)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var status StatusObject
	_ = json.NewDecoder(response.Body).Decode(&status)

	if status.MatchedCount == 0 {
		fmt.Fprintf(writer, "For URL: %s\nProduct details scraped and stored in database with ID: %s\n", form_data.URL, status.InsertedID)
	} else {
		if status.ModifiedCount == 0 {
			fmt.Fprintf(writer, "For URL: %s\nProduct details already exists in Database, and they match.\n", form_data.URL)
		} else {
			fmt.Fprintf(writer, "For URL: %s\nProduct details already exists in Database, and are updated with latest.\n", form_data.URL)
		}
	}

}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/scraper", getFunc).Methods("GET")
	router.HandleFunc("/scraper", postFunc).Methods("POST")
	log.Fatal(http.ListenAndServe(":3030", router))
}
