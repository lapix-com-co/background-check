package reply

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCharLength_Is(t *testing.T) {
	type input struct {
		name     string
		question string
		got      bool
	}
	cases := []*input{
		{
			"should match the given question",
			"¿Escriba la cantidad de letras del primer nombre de la persona a la cual esta expidiendo el certificado?",
			true,
		},
		{
			"should not match the given question",
			"¿Cúal es la capital de Amazonas?",
			false,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			s := &CharLength{}

			require.Equal(t, cs.got, s.Is(User{}, cs.question))
		})
	}
}

func TestCharLength_Answer(t *testing.T) {
	user := User{
		Name: Name{
			FirstName:     "John",
			MiddleName:    "Alexander",
			LastName:      "Doe",
			SecondSurname: "Bell",
		},
	}

	t.Run("should get the user user's first name", func(t *testing.T) {
		question := "¿Escriba la cantidad de letras del primer nombre de la persona a la cual esta expidiendo el certificado?"
		s := &CharLength{}

		answer, _, _ := s.Answer(user, question)
		require.Equal(t, "4", answer)
	})

	t.Run("should get the user user's middle name", func(t *testing.T) {
		question := "¿Escriba la cantidad de letras del segundo nombre de la persona a la cual esta expidiendo el certificado?"
		s := &CharLength{}

		answer, _, _ := s.Answer(user, question)
		require.Equal(t, "9", answer)
	})

	t.Run("should get the user user's last name", func(t *testing.T) {
		question := "¿Escriba la cantidad de letras del primer apellido de la persona a la cual esta expidiendo el certificado?"
		s := &CharLength{}

		answer, _, _ := s.Answer(user, question)
		require.Equal(t, "3", answer)
	})

	t.Run("should get the user user's second surname", func(t *testing.T) {
		question := "¿Escriba la cantidad de letras del segundo apellido de la persona a la cual esta expidiendo el certificado?"
		s := &CharLength{}

		answer, _, _ := s.Answer(user, question)
		require.Equal(t, "4", answer)
	})

	t.Run("should get the false if ask somthing unknown", func(t *testing.T) {
		question := "¿Escriba la cantidad de letras del pasaporte de la persona a la cual está expidiendo el certificado?"
		s := &CharLength{}

		_, ok, _ := s.Answer(user, question)
		require.Equal(t, false, ok)
	})

	t.Run("should get the error if the given question does not match", func(t *testing.T) {
		question := "¿Escriba la cantidad de numeros del pasaporte de la persona a la cual está expidiendo el certificado?"
		s := &CharLength{}

		_, _, err := s.Answer(user, question)
		require.Error(t, err)
	})
}
