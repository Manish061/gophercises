package sitemap

import (
	"errors"
	"gophercises/sitemap/href"
	"net/http"
	"strings"
	"sync"
)

var visitDFS = struct {
	m map[string]error
	sync.Mutex
}{m: make(map[string]error)}

const maxDepth = 2

var errLoading = errors.New("url load in progress") // sentinel value

var baseURL string

//SiteMapper makes a sitemap of the given URL until depth 'depth'
//and returns a standard sitemap.xml string
func SiteMapper(url string, depth int) (string, error) {

	var visit = make(map[string]bool)

	type siteNode struct {
		URL   string
		Depth int
	}
	xmlData := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"

	xmlURLData := "\n\t<url>\n\t\t<loc>"

	var queue = make([]siteNode, 0)
	if depth > maxDepth || depth <= 0 {
		depth = maxDepth
	}
	queue = append(queue, siteNode{url, 0})
	data := ""
	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]
		currURL := top.URL
		currDepth := top.Depth
		if currDepth >= depth {
			continue
		}
		data = xmlURLData + currURL + "</loc>\n\t</url>"
		xmlData += data

		bodyData, err := fetch(url)
		if err != nil {
			return "", err
		}
		for _, link := range bodyData {
			_, found := visit[link]
			if !found {
				queue = append(queue, siteNode{link, currDepth + 1})
				visit[link] = true
			}
		}
	}
	xmlData += "\n</urlset>"
	return xmlData, nil
}

// SiteMapperDFS makes a sitemap of the given URL until depth 'depth'
// and returns a standard sitemap.xml string
func SiteMapperDFS(url string, depth int) (string, error) {

	baseURL = url

	xmlData := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"
	xmlURLData := "\n\t<url>\n\t\t<loc>"
	if depth > maxDepth || depth <= 0 {
		depth = maxDepth
	}
	siteMapperDFSUtil(url, depth)
	for key := range visitDFS.m {
		xmlData += xmlURLData + key + "</loc>\n\t</url>"
	}
	return xmlData + "\n</urlset>", nil
}

func siteMapperDFSUtil(url string, depth int) {
	if depth <= 0 {
		return
	}
	visitDFS.Lock()
	if _, found := visitDFS.m[url]; found {
		visitDFS.Unlock()
		return
	}
	visitDFS.m[url] = errLoading
	visitDFS.Unlock()

	bodyData, err := fetch(url)

	visitDFS.Lock()
	visitDFS.m[url] = err
	visitDFS.Unlock()
	done := make(chan bool)
	for _, link := range bodyData {
		go func(link string) {
			siteMapperDFSUtil(link, depth-1)
			done <- true
		}(link)
	}
	for range bodyData {
		<-done
	}
}

func fetch(url string) ([]string, error) {
	response, errGetting := http.Get(url)
	if errGetting != nil {
		return nil, errGetting
	}
	bodyData, errParsing := href.Parse(response.Body)
	if errParsing != nil {
		return nil, errParsing
	}
	linkStr := ""
	var parsedURL = make([]string, 0)
	for _, link := range bodyData {
		linkStr = link.Href
		if strings.Contains(linkStr, "#") {
			continue
		}
		if len(linkStr) > 1 && string(linkStr[0]) == "/" {
			linkStr = url + linkStr
		} else {
			if !strings.HasPrefix(linkStr, baseURL) {
				continue
			}
		}
		parsedURL = append(parsedURL, linkStr)
	}
	return parsedURL, nil
}
