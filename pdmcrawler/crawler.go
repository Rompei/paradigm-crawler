package pdmcrawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"regexp"
	"strings"
	"time"
)

// BaseURL is url for wikipedia.
const BaseURL = "https://en.wikipedia.org"

// Crawler is crawler module
type Crawler struct {
	Time         time.Time `json:"time"`
	LanguageTree *Language `json:"languageTree"`
}

// NewCrawler is constructor of Crawler.
func NewCrawler() *Crawler {
	return &Crawler{}
}

// Language is object of programing language.
type Language struct {
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	Descendents []Language `json:"descendents"`
}

// NewLanguage is constructor of Language.
func NewLanguage(name, url string) *Language {
	return &Language{
		Name: name,
		URL:  url,
	}
}

// ShowLanguages shows languages with descendents.
func (l *Language) ShowLanguages() {
	fmt.Println(l.Name)
	for _, v := range l.Descendents {
		v.ShowLanguages()
	}
}

// CrawlChecker is wrapper of crawled languages.
type CrawlChecker struct {
	crawled []string `json:"crawled"`
}

// NewCrawlChecker is constructot of CrawlChecker.
func NewCrawlChecker() *CrawlChecker {
	return &CrawlChecker{}
}

// AddCrawled adds crawled language to crawled.
func (c *CrawlChecker) AddCrawled(name string) {
	c.crawled = append(c.crawled, name)
}

// ShowCrawled shows crawled languages.
func (c *CrawlChecker) ShowCrawled() {
	fmt.Println(strings.Join(c.crawled, ","))
}

// GetCrawled returns crawled languages.
func (c *CrawlChecker) GetCrawled() []string {
	return c.crawled
}

// Len returns length of crawled languages.
func (c *CrawlChecker) Len() int {
	return len(c.crawled)
}

// Crawl crawls wikipedia.
func (cl *Crawler) Crawl(l *Language, c *CrawlChecker) (*Language, error) {
	var err error
	cl.LanguageTree, err = cl.crawlCore(l, c)

	return cl.LanguageTree, err
}
func (cl *Crawler) crawlCore(l *Language, c *CrawlChecker) (*Language, error) {

	// Start crawling to the page.
	doc, err := goquery.NewDocument(l.URL)
	if err != nil {
		return l, err
	}
	var prevText string
	doc.Find(".vevent").Find("tr").Each(func(i int, s *goquery.Selection) {
		if prevText == "Influenced" {
			s.Find("td").Find("a").Each(func(i int, s *goquery.Selection) {

				// Getting href for next page.
				href, ok := s.Attr("href")
				if !ok {
					fmt.Println("This page was wrong.")
					return
				}

				// Getting name of the language.
				name, ok := s.Attr("title")
				if !ok {
					fmt.Println("This page was wrong.")
					return
				}

				// Triming name.
				cl.trimName(&name)
				if c.Len() == 0 {

					// First time.

					url := BaseURL + href
					fmt.Printf("Start searching: %s, URL: %s\n", name, url)
					desLang := NewLanguage(name, url)

					// Add language to crawled languages.
					c.AddCrawled(name)
					desLang, err = cl.crawlCore(desLang, c)
					if err != nil {
						fmt.Println(err)
					}
					l.Descendents = append(l.Descendents, *desLang)
				} else {
					isExist := false
					for _, v := range c.GetCrawled() {
						if v == name {

							// If the language was already crawled, it ignore that.

							isExist = true
						}
					}
					if !isExist {
						url := BaseURL + href
						fmt.Printf("Start searching: %s, URL: %s\n", name, url)
						desLang := NewLanguage(name, url)

						// Add language to crawled languages.
						c.AddCrawled(name)
						desLang, err = cl.crawlCore(desLang, c)
						if err != nil {
							fmt.Println(err)
						}
						l.Descendents = append(l.Descendents, *desLang)
					}
				}
			})
		}

		// Take th text for searching row of Inflenced.
		prevText = s.Find("th").Text()
	})

	return l, nil
}

func (cl *Crawler) trimName(src *string) {
	re, _ := regexp.Compile(`\(programming language\)|\(page does not exist\)`)
	*src = re.ReplaceAllString(*src, "")
}

// Dump stores language tree as JSON.
func (cl *Crawler) Dump(fname string) error {

	// Storing language tree as JSON format.

	cl.Time = time.Now()
	b, err := json.Marshal(cl)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()
	file.Truncate(0)
	out.WriteTo(file)
	return nil
}
