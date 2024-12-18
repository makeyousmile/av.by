package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Data []struct {
	ID      int `json:"id"`
	Country struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
		Emoji string `json:"emoji"`
		Code  string `json:"code"`
	} `json:"country"`
	Number string `json:"number"`
}
type Phone struct {
	number string
	date   time.Time
}

func getPhones(pages int) []Phone {
	phones := []Phone{}
	for i := 0; i < pages; i++ {
		for _, id := range GetIds(strconv.FormatInt(int64(i), 10)) {
			phone := Phone{}
			phone.number = GetNumber(id)
			phone.date = time.Now()
			if phone.number != "" {
				phones = append(phones, phone)
			}
			//writeTofile(phone.number + "\n")
		}
	}
	return phones
}

func getPhone(pages int, phones chan Phone) {

	for i := 0; i < pages; i++ {
		for _, id := range GetIds(strconv.FormatInt(int64(i), 10)) {
			phone := Phone{}
			phone.number = GetNumber(id)
			phone.date = time.Now()
			if phone.number != "" {
				log.Print(phone.number)
			}
			//writeTofile(phone.number + "\n")
		}
	}

}

func ScrapPage(pageNumber string) []string {

	links := []string{}
	c := colly.NewCollector()

	c.OnResponse(func(response *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			log.Print(err)
		} else {
			doc.Find(".listing-item__link").Each(func(i int, selection *goquery.Selection) {
				href, exist := selection.Attr("href")
				if exist {
					links = append(links, href)
				}
			})
		}

	})
	url := getFilterForScraper() + "&page=" + pageNumber + "&sort=4"
	err := c.Visit(url)
	if err != nil {
		log.Print(err)
	}
	log.Print("visit page")
	log.Print(url)
	c.Wait()
	//log.Print(links)
	return links
}
func GetIds(pageNumber string) []string {
	Ids := []string{}
	for _, link := range ScrapPage(pageNumber) {
		i := len(link) - 9
		id := link[i:]
		//log.Print(id)
		Ids = append(Ids, id)
	}
	return Ids
}
func GetNumber(id string) string {
	phone := ""
	data := Data{}

	link := "https://api.av.by/offers/" + id + "/phones"
	//log.Print(link)
	err := getJson(link, &data)
	if err != nil {
		log.Print("Error get phones")
	} else {
		phone = data[0].Number
	}
	return phone
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
func getFilterForScraper() string {
	link := "https://cars.av.by/filter?seller_type[0]=1"

	return link
}
