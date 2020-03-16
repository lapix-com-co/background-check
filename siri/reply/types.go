package reply

import (
	"errors"
)

var (
	// ErrInvalidQuestion is returned when the strategy is evaluated with and invalid question.
	ErrInvalidQuestion = errors.New("question does not match")
)

type StrategyFinder interface {
	Strategy(User, string) Strategy
}

// User refers to the user who owns the verification.
type User struct {
	Name     Name
	Document Document
}

// Name refers to the user's name parts.
type Name struct {
	FirstName     string
	MiddleName    string
	LastName      string
	SecondSurname string
}

// Document refers to the user's document data.
type Document struct {
	Type         string
	Number       string
	CreationDate string
}

// Strategy will answer a question depending on the given input.
type Strategy interface {
	Is(User, string) bool
	Answer(User, string) (string, bool, error)
}
