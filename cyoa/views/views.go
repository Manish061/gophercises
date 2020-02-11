package views

import (
	Story "gophercises/cyoa/story"
	"html/template"
	"net/http"
	"strings"
)

//ViewData is the data interface that is used for passing values from
//the json data to the view template
type viewData struct {
	PageTitle string
	Title     string
	Story     string
	Options   []Story.CYOAStoryOptions
}

//Views return the templated html response to the user
func Views() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story, err := Story.Story()
		if err != nil {
			panic(err)
		}
		query := r.URL.Query()
		page := query.Get("page")
		returnData := Story.CYOAStoryData{}
		lookup := true
		if page != "" {
			returnData, lookup = story[page]
			if lookup == false {
				tmpl := template.Must(template.ParseFiles("sorry.gohtml"))
				executingError := tmpl.Execute(w, nil)
				if executingError != nil {
					panic(executingError)
				}
				return
			}
		} else {
			http.Redirect(w, r, "/view?page=intro", http.StatusSeeOther)
			return
		}
		template := template.Must(template.ParseFiles("story.gohtml"))
		data := viewData{
			PageTitle: "Choose Your Own Adventure",
			Title:     returnData.Title,
			Story:     strings.Join(returnData.Story, " "),
			Options:   returnData.Options,
		}
		ExecutingError := template.Execute(w, data)
		if ExecutingError != nil {
			panic(ExecutingError)
		}
	}
}
