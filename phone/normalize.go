package phone

import (
	"regexp"
)

// Normalize a phone numbber as input and returns
//the normalized version of the phone
//Eg:- input: (123) 456 7892, output: 1234567892
func Normalize(phoneNumber string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phoneNumber, "")
}
