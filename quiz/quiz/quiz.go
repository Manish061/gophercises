package quiz

//Quiz template for gophercise 1
type Quiz struct {
	FileName 		string
	Question        []string
	Answer          []string
	CreatedBy 		string
}

func (quiz Quiz) String() string {
	return "Quiz created by " + quiz.CreatedBy
}
