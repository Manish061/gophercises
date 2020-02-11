package main

import (
	"errors"
	"flag"
	"fmt"
	"gophercises/sitemap-worker-pool/href"
	"net/http"
	"strings"
	"sync"
)

var mapped = struct {
	m map[string]error
	sync.Mutex
}{m: make(map[string]error)}

var errLoading = errors.New("url load in progress") // sentinel value

// Map uses fetcher to recursively map
// pages starting with url, to a maximum of depth.
func Map(url string, depth int) {
	if depth <= 0 {
		fmt.Printf("<- Done with %v, depth 0.\n", url)
		return
	}

	mapped.Lock()
	if _, ok := mapped.m[url]; ok {
		mapped.Unlock()
		fmt.Printf("<- Done with %v, already fetched.\n", url)
		return
	}
	// We mark the url to be loading to avoid others reloading it at the same time.
	mapped.m[url] = errLoading
	mapped.Unlock()

	urls, errMapping := Fetch(url)

	mapped.Lock()
	mapped.m[url] = errMapping
	mapped.Unlock()

	if errMapping != nil {
		fmt.Printf("<- Error on %v: %v\n", url, errMapping)
		return
	}
	done := make(chan bool)
	for i, u := range urls {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		go func(url string) {
			Map(url, depth-1)
			done <- true
		}(u)
	}
	for i, u := range urls {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls), u)
		<-done
	}
	fmt.Printf("<- Done with %v\n", url)
}

func main() {
	url := flag.String("url", "", "a valid url")
	depth := flag.Int("depth", 2, "depth of site mapping")
	flag.Parse()
	if len(*url) < 1 {
		flag.Usage()
		return
	}
	Map(*url, *depth)
}

// Fetch returns the body of URL and
// a slice of URLs found on that page.
func Fetch(url string) (urls []string, err error) {

	response, errFetching := http.Get(url)
	if errFetching != nil {
		return nil, errFetching
	}
	bodyData, errParsing := href.Parse(response.Body)
	if errParsing != nil {
		return nil, errParsing
	}
	linkStr := ""
	for _, link := range bodyData {
		linkStr = link.Href
		if strings.Contains(linkStr, "#") {
			continue
		}
		if len(linkStr) > 1 && string(linkStr[0]) == "/" {
			linkStr = url + linkStr
		} else {
			if !strings.HasPrefix(linkStr, url) {
				continue
			}
		}
		urls = append(urls, linkStr)
	}
	return urls, nil
}
