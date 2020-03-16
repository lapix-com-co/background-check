package reply

import (
	"regexp"
	"strings"
)

var (
	namesPattern = regexp.MustCompile("cual es el primer nombre de la persona a la cual esta expidiendo el certificado")
)

type Names struct {}

func (n Names) Is(u User, i string) bool {
	return namesPattern.MatchString(strings.ToLower(i))
}

func (n Names) Answer(u User, i string) (string, bool, error) {
	return u.Name.FirstName, true, nil
}

