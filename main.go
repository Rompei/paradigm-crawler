package main

import (
	"flag"
	"fmt"
	"github.com/Rompei/paradigm-crawler/pdmcrawler"
	"os"
)

func main() {

	// Taking name and url of root language from argument.

	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Println("not enough argument")
		os.Exit(1)
	}
	name := flag.Args()[0]
	url := flag.Args()[1]
	output := flag.Args()[2]

	// Start crawling.
	crawler := pdmcrawler.NewCrawler()
	_, err := crawler.Crawl(pdmcrawler.NewLanguage(name, url), pdmcrawler.NewCrawlChecker())
	if err != nil {
		fmt.Println(err)
	}

	// Store.
	err = crawler.Dump(output)
	if err != nil {
		panic(err)
	}
}
