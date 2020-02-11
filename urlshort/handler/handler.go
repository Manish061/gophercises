package handler

import (
	"net/http"

	"encoding/json"

	"gopkg.in/yaml.v2"
)

//MapHandler takes a pathToUrls a (key,value) pair and returns a http handler func,
//If path is not defined then a simple http serve is done, else route is redirected.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

//YAMLHandler takes URL mapping yaml configuration as bytes,
//builds a map of the passed paths URLs and returns a handler function
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

//JSONHandler takes URL mapping json config as bytes,
// builds a map and returns a handler function
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildJSONMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func buildMap(pathUrls []PathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func buildJSONMap(pathUrls []path) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte) ([]PathURL, error) {
	var pathUrls []PathURL
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func parseJSON(data []byte) ([]path, error) {
	var pathUrls []path
	err := json.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

//PathURL blueprint to hold the parsed yaml data
type PathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type path struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}
