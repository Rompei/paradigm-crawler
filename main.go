package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
)

// BaseURL is url for wikipedia.
const BaseURL = "https://en.wikipedia.org"

type Language struct {
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	Descendents []Language `json:"descendents"`
}

func (l Language) ShowLanguages() {
	fmt.Println(l.Name)
	for _, v := range l.Descendents {
		v.ShowLanguages()
	}
}

type CrawlChecker struct {
	crawled []string `json:"crawled"`
}

func (c CrawlChecker) AddCrawled(name string) {
	c.crawled = append(c.crawled, name)
}

func (c CrawlChecker) GetCrawled() []string {
	return c.crawled
}

func (c CrawlChecker) Len() int {
	return len(c.crawled)
}

func crawl(l *Language, c *CrawlChecker) (*Language, error) {
	doc, err := goquery.NewDocument(l.URL)
	if err != nil {
		return l, err
	}
	var prevText string
	doc.Find(".vevent").Find("tr").Each(func(i int, s *goquery.Selection) {
		if prevText == "Influenced" {
			s.Find("td").Find("a").Each(func(i int, s *goquery.Selection) {
				href, ok := s.Attr("href")
				if !ok {
					fmt.Println("This page was wrong.")
					return
				}
				name, ok := s.Attr("title")
				if !ok {
					fmt.Println("This page was wrong.")
					return
				}
				if c.Len() == 0 {
					url := BaseURL + href
					fmt.Printf("Start searching: %s, URL: %s\n", name, url)
					desLang := &Language{
						Name: name,
						URL:  url,
					}
					c.AddCrawled(name)
					desLang, err = crawl(desLang, c)
					if err != nil {
						fmt.Println(err)
					}
					l.Descendents = append(l.Descendents, *desLang)
				} else {
					isExist := false
					for _, v := range c.GetCrawled() {
						if v == name {
							isExist = true
						}
					}
					if isExist {
						fmt.Printf("Stop searching: %s\n", name)
						l.Descendents = append(l.Descendents, Language{Name: name})
					} else {
						url := BaseURL + href
						fmt.Printf("Start searching: %s, URL: %s\n", name, url)
						desLang := &Language{
							Name: name,
							URL:  url,
						}
						c.AddCrawled(name)
						desLang, err = crawl(desLang, c)
						if err != nil {
							fmt.Println(err)
						}
						l.Descendents = append(l.Descendents, *desLang)
					}
				}
			})
		}
		prevText = s.Find("th").Text()
	})

	return l, nil
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Error")
		os.Exit(1)
	}
	url := flag.Args()[0]

	l := &Language{
		Name: "C",
		URL:  url,
	}

	lc := &CrawlChecker{}

	language, err := crawl(l, lc)
	if err != nil {
		fmt.Println(err)
	}

	language.ShowLanguages()
}
