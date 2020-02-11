package main

import (
	"flag"
	"fmt"
	urlShort "gophercises/urlshort/handler"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	yaml := flag.String("yaml", "maps.yaml", "a valid yaml file of the format\n- path: /some-path\n  url: https://www.some-url.com/demo")
	json := flag.String("json", "maps.json", "a valid json file of the format- {\"path\":\"/some-path\",\"url\":\"https://www.some-url.com/demo\"}")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlShort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	handler := mapHandler

	// if *yaml != "" {
	yamlFile, err := os.Open(*yaml)
	if err != nil {
		panic(err)
	}
	defer yamlFile.Close()

	yamlData, err := ioutil.ReadFile(*yaml)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := urlShort.YAMLHandler(yamlData, mapHandler)
	handler = yamlHandler
	if err != nil {
		panic(err)
	}
	// } else if *json != "" {
	jsonFile, err := os.Open(*json)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadFile(*json)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlShort.JSONHandler(jsonData, handler)
	handler = jsonHandler
	if err != nil {
		panic(err)
	}
	// }
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
