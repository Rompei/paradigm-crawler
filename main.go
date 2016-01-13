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
	l := pdmcrawler.NewLanguage(name, url)
	cc := pdmcrawler.NewCrawlChecker()
	_, err := crawler.Crawl(l, cc)
	if err != nil {
		fmt.Println(err)
	}

	// Store.
	err = crawler.Dump(output)
	if err != nil {
		panic(err)
	}
}
