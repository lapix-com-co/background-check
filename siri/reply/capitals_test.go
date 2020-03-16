package reply

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCapitals_Is(t *testing.T) {
	type input struct {
		name     string
		question string
		got      bool
	}

	cases := []*input{
		{
			name:     "should match the capitals question",
			question: "¿ Cual es la Capital de Antioquia?",
			got:      true,
		},
		{
			name:     "should match the capitals question",
			question: "¿ Cual es la Capital del Vallle del Cauca?",
			got:      true,
		},
		{
			name:     "should not match the capitals question",
			question: "¿Escriba los tres primeros digitos del documento a consultar?",
			got:      false,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			s := &Capitals{}
			require.Equal(t, cs.got, s.Is(User{}, cs.question))
		})
	}
}

func TestCapitals_Answer(t *testing.T) {
	type input struct {
		name      string
		question  string
		answer    string
		hasAnswer bool
		err       error
	}

	cases := []*input{
		{
			name:      "should return a valid answer for the question",
			question:  "¿ Cual es la Capital de Antioquia?",
			answer:    "Medellín",
			hasAnswer: true,
		},
		{
			name:      "should return a valid answer for the question",
			question:  "¿ Cual es la Capital de Colombia (sin tilde)?",
			answer:    "Bogotá",
			hasAnswer: true,
		},
		{
			name:      "should return a valid answer for the question",
			question:  "¿Cual es la Capital del Vallle del Cauca?",
			answer:    "Cali",
			hasAnswer: true,
		},
		{
			name:      "should return false to an unknow location",
			question:  "¿ Cual es la Capital de Apartado?",
			hasAnswer: false,
		},
		{
			name:     "should return an error to an unknow question",
			question: "¿cómo es ser desconocidos para los conocidos de los deconocidos?",
			err:      ErrInvalidQuestion,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			s := NewCapitals()
			a, ok, err := s.Answer(User{}, cs.question)

			require.Equal(t, cs.answer, a)
			require.Equal(t, cs.hasAnswer, ok)

			if cs.err != nil && !errors.Is(err, cs.err) {
				t.Errorf("expect err %v but got %v", cs.err, err)
				return
			}

			if cs.err == nil && err != nil {
				t.Errorf("expect nil but got = %v", err)
			}
		})
	}
}

func TestNoAccentsCapitals_Answer(t *testing.T) {
	type input struct {
		name      string
		question  string
		answer    string
		hasAnswer bool
		err       error
	}

	cases := []*input{
		{
			name:      "should return a valid answer for the question",
			question:  "¿ Cual es la Capital de Antioquia?",
			answer:    "Medellin",
			hasAnswer: true,
		},
		{
			name:      "should return a valid answer for the question",
			question:  "¿Cual es la Capital del Valle del Cauca?",
			answer:    "Cali",
			hasAnswer: true,
		},
		{
			name:      "should return false to an unknow location",
			question:  "¿ Cual es la Capital de Apartado?",
			hasAnswer: false,
		},
		{
			name:     "should return an error to an unknow question",
			question: "¿cómo es ser desconocidos para los conocidos de los deconocidos?",
			err:      ErrInvalidQuestion,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			s := NewNoAccentsCapitals()
			a, ok, err := s.Answer(User{}, cs.question)

			require.Equal(t, cs.answer, a)
			require.Equal(t, cs.hasAnswer, ok)

			if cs.err != nil && !errors.Is(err, cs.err) {
				t.Errorf("expect err %v but got %v", cs.err, err)
				return
			}

			if cs.err == nil && err != nil {
				t.Errorf("expect nil but got = %v", err)
			}
		})
	}
}
