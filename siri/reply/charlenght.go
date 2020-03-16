package reply

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	charLengthPattern = regexp.MustCompile(`escriba la cantidad de letras del ([a-zA-Z\s]+) de la persona a la cual`)
)

var _ Strategy = &CharLength{}

// CharLength handles the users data char length: name, last name, middle name and second surname.
type CharLength struct{}

// Is should return true if the question ask the letters count.
func (c CharLength) Is(u User, i string) bool {
	return charLengthPattern.MatchString(strings.ToLower(i))
}

// Answer will return the user's data length.
func (c CharLength) Answer(u User, i string) (string, bool, error) {
	matches := charLengthPattern.FindAllStringSubmatch(strings.ToLower(i), 1)

	if len(matches) == 0 {
		return "", false, ErrInvalidQuestion
	}

	part := matches[0][1]

	switch part {
	case "primer nombre":
		return strconv.Itoa(len(u.Name.FirstName)), true, nil
	case "segundo nombre":
		return strconv.Itoa(len(u.Name.MiddleName)), true, nil
	case "primer apellido":
		return strconv.Itoa(len(u.Name.LastName)), true, nil
	case "segundo apellido":
		return strconv.Itoa(len(u.Name.SecondSurname)), true, nil
	}

	return "", false, nil
}
