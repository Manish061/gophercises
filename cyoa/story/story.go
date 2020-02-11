package story

import (
	"encoding/json"
	"os"
)

//CYOAStory defines the model to read from json.
//Also Validates the json file
type CYOAStory map[string]CYOAStoryData

//CYOAStoryData is the story data for a particular story arc
type CYOAStoryData struct {
	Title   string             `json:"title"`
	Story   []string           `json:"story"`
	Options []CYOAStoryOptions `json:"options"`
}

//CYOAStoryOptions links one or more story arc to another
type CYOAStoryOptions struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

//Story reads from input story json file and builds a map for all the story arc
func Story() (CYOAStory, error) {
	fileReader, err := os.Open("gopher.json")
	if err != nil {
		panic(err)
	}
	defer fileReader.Close()
	d := json.NewDecoder(fileReader)
	var story CYOAStory
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
