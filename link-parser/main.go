package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gophercises/link-parser/htmlgo"
	"log"
	"os"
	"strings"
)

func main() {
	file := flag.String("file", "story.gohtml", "a html file")
	str := flag.String("str", "", "a html string")
	flag.Parse()

	if *str != "" {
		output, errParsing := htmlgo.Parse(strings.NewReader(*str))
		if errParsing != nil {
			log.Fatal(errParsing)
		}
		outputBytes, err := json.Marshal(output)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(outputBytes))
	} else {
		fileData, errOpen := os.Open(*file)
		if errOpen != nil {
			log.Fatal(errOpen)
		}
		output, errParsing := htmlgo.Parse(fileData)
		if errParsing != nil {
			log.Fatal(errParsing)
		}
		outputBytes, err := json.Marshal(output)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(outputBytes))
	}
	return
}
