package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
)

func crawl(url string, errCh chan error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		errCh <- err
		return
	}
	var prevText string
	doc.Find(".vevent").Find("tr").Each(func(i int, s *goquery.Selection) {
		if prevText == "Influenced" {
			fmt.Println(s.Text())
			return
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
	go crawl(url, errCh)

	err := <-errCh
	if err != nil {
		fmt.Printf("%v", err)
	}
}
