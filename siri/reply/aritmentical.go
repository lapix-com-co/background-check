package reply

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	arithmeticalPattern = regexp.MustCompile(`cuanto es ([\d\s+x%-]+) ?\?$`)
)

type Arithmetical struct{}

func (a Arithmetical) Is(u User, i string) bool {
	return arithmeticalPattern.MatchString(strings.ToLower(i))
}

func (a Arithmetical) Answer(u User, i string) (string, bool, error) {
	matched := arithmeticalPattern.FindAllStringSubmatch(strings.ToLower(i), 1)

	if len(matched) == 0 {
		return "", false, ErrInvalidQuestion
	}

	parts := strings.Split(matched[0][1], " ")

	firstNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", false, fmt.Errorf("the numbers in the quiestion are not valid %w", ErrInvalidQuestion)
	}
	secondNumber, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", false, fmt.Errorf("the numbers in the quiestion are not valid %w", ErrInvalidQuestion)
	}

	switch parts[1] {
	case "+":
		return strconv.Itoa(firstNumber + secondNumber), true, nil
	case "x":
		return strconv.Itoa(firstNumber * secondNumber), true, nil
	case "-":
		return strconv.Itoa(firstNumber - secondNumber), true, nil
	}

	return "", false, nil
}
