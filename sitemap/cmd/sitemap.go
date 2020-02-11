package main

import (
	"flag"
	"fmt"
	"gophercises/sitemap"
	"io/ioutil"
	"log"
)

func main() {
	url := flag.String("url", "", "a valid url")
	depth := flag.Int("depth", 2, "depth of site mapping")
	flag.Parse()
	if len(*url) < 1 {
		flag.Usage()
		return
	}
	go func() {
		xmlFileStr, errSiteMapping := sitemap.SiteMapper(*url, *depth)
		if errSiteMapping != nil {
			log.Fatal(errSiteMapping)
		}
		xmlFileData := []byte(xmlFileStr)
		errWritingFile := ioutil.WriteFile("sitemap.xml", xmlFileData, 0664)
		if errWritingFile != nil {
			log.Fatal(errWritingFile)
		}
	}()

	go func() {
		xmlFileStr, errSiteMapping := sitemap.SiteMapperDFS(*url, *depth)
		if errSiteMapping != nil {
			log.Fatal(errSiteMapping)
		}
		xmlFileData := []byte(xmlFileStr)
		errWritingFile := ioutil.WriteFile("sitemap-1.xml", xmlFileData, 0664)
		if errWritingFile != nil {
			log.Fatal(errWritingFile)
		}
	}()

	fmt.Scanf("\n")
}
