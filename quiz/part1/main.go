package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// QuizFileFormatError populates error if file is not in valid format, having two comma separated values per line
func QuizFileFormatError(line []string, fileName string) {
	if len(line) != 2 {
		panic("expected a csv file in the format of 'question,answer'")
	}
	return
}

// Check panics if runtime error
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	userName := flag.String("name", "QuizMaker", "a string")
	fileName := flag.String("file", "problems.csv", "a csv file in the format of 'question,answer'")
	shuffle := flag.Bool("shuffle", false, "a boolean indicating whether to shuffle the questions")
	flag.Parse()
	quizFile, err := os.Open(*fileName)
	Check(err)
	defer quizFile.Close()
	reader := csv.NewReader(bufio.NewReader(quizFile))
	question := make([]string, 0)
	answer := make([]string, 0)
	for line, err := reader.Read(); err != io.EOF; {
		QuizFileFormatError(line, *fileName)
		question = append(question, line[0])
		answer = append(answer, line[1])
		line, err = reader.Read()
	}
	if *shuffle {
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Shuffle(len(question), func(i, j int) {
			question[i], question[j] = question[j], question[i]
			answer[i], answer[j] = answer[j], answer[i]
		})
	}
	quiz := Quiz{
		FileName:  *fileName,
		Question:  question,
		Answer:    answer,
		CreatedBy: *userName,
	}
	user := User{
		UserName:          *userName,
		UserQuizFileName:  *fileName,
		QuestionsAnswered: 0,
	}
	givenAnswer := ""
	score := 0
	for i := 0; i < len(quiz.Question); i++ {
		fmt.Printf("Problem #%d: %s? ", i, quiz.Question[i])
		fmt.Scanf("%s\n", &givenAnswer)
		if givenAnswer != "" {
			if strings.TrimSpace(givenAnswer) == strings.TrimSpace(quiz.Answer[i]) {
				score++
				user.Score = score
			}
			user.QuestionsAnswered++
		}
	}
	fmt.Printf("%s/%d", user, len(question))
}

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

//Quiz template for gophercise 1
type Quiz struct {
	FileName  string
	Question  []string
	Answer    []string
	CreatedBy string
}

func (quiz Quiz) String() string {
	return "Quiz created by " + quiz.CreatedBy
}
