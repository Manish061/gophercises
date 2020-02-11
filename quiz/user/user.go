package user

 import "strconv"

//User template for gophercise 1
type User struct {
	UserName          string
	UserQuizFileName  string
	QuestionsAnswered int
	Score             int
}

func (user User) String() string {
	// return fmt.Sprintf("UserName: %v, Score: %v", user.UserName, user.Score)
	return "UserName: " + user.UserName + ", Score: " + strconv.Itoa(user.Score)
}
