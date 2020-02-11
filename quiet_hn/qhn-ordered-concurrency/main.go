package main

import (
	"flag"
	"fmt"
	"gophercises/quiet_hn/hn"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	var port int
	var numStories int
	flag.IntVar(&port, "port", 6060, "port to run the server on")
	flag.IntVar(&numStories, "stories", 30, "number of stories to fetch from the Hacker News API")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("../views/index.gohtml"))

	mux := http.NewServeMux()
	mux.Handle("/", handler(numStories, tpl))
	fmt.Printf("server running on port :%d", port)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		topNumStories, errFetchingTopNumStories := getTopNumStories(numStories)
		if errFetchingTopNumStories != nil {
			http.Error(w, errFetchingTopNumStories.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: topNumStories,
			Time:    time.Now().Sub(start),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	}
}

func getTopNumStories(numStories int) ([]item, error) {
	var hnClient hn.Client
	storyIds, errFetchingTopStories := hnClient.TopItems()
	if errFetchingTopStories != nil {
		return nil, errFetchingTopStories
	}
	type result struct {
		index int
		item  item
		err   error
	}
	resultCh := make(chan result)
	var stories []item
	for i := 0; i < numStories; i++ {
		go func(index, storyID int) {
			hnItem, errGettingStoryItem := hnClient.GetItem(storyID)
			if errGettingStoryItem != nil {
				resultCh <- result{index: index, err: errGettingStoryItem}
			} else {
				resultCh <- result{index: index, item: parseHNItem(hnItem)}
			}
		}(i, storyIds[i])
	}
	var results []result
	for i := 0; i < numStories; i++ {
		results = append(results, <-resultCh)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories, nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}
