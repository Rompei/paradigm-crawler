package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
)

// BaseURL is url for wikipedia.
const BaseURL = "https://en.wikipedia.org/wiki/"

type Language struct{
	Name string `json:"name"`
	URL string `json:"url"`
	Ansestors []Language`json:"ansestors"`
}

func crawl(l *Language, errCh chan error) {
	doc, err := goquery.NewDocument(l.URL)
	if err != nil {
		errCh <- err
		return
	}
	var prevText string
	doc.Find(".vevent").Find("tr").Each(func(i int, s *goquery.Selection) {
		if prevText == "Influenced" {
			s.Find("td").Find("a").Each(func(i int, s *goquery.Selection){
				fmt.Println(s.Text())
			})
		}
		prevText = s.Find("th").Text()
	})

	errCh <- nil
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Error")
		os.Exit(1)
	}
	url := flag.Args()[0]

	errCh := make(chan error)
	go crawl(&Language{
		URL:url,
		Name:"C",
	}, errCh)

	err := <-errCh
	if err != nil {
		fmt.Printf("%v", err)
	}
}
