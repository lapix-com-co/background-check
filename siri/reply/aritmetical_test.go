package reply

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArithmetical_Is(t *testing.T) {
	t.Run("should match the given question with product", func(t *testing.T) {
		question := "¿ Cuanto es 3 X 3 ?"
		a := Arithmetical{}
		result := a.Is(User{}, question)
		require.True(t, result)
	})
	t.Run("should not match the given question", func(t *testing.T) {
		question := "¿ Cual es la capital de Noruega?"
		a := Arithmetical{}
		result := a.Is(User{}, question)
		require.False(t, result)
	})
	t.Run("should match the given question with addition", func(t *testing.T) {
		question := "¿ Cuanto es 3 + 3 ?"
		a := Arithmetical{}
		result := a.Is(User{}, question)
		require.True(t, result)
	})
	t.Run("should match the given question with substraction", func(t *testing.T) {
		question := "¿ Cuanto es 3 - 3 ?"
		a := Arithmetical{}
		result := a.Is(User{}, question)
		require.True(t, result)
	})
}

func TestArithmetical_Answer(t *testing.T) {
	t.Run("should return the product", func(t *testing.T) {
		question := "¿ Cuanto es 3 X 3 ?"
		a := Arithmetical{}
		result, _, _ := a.Answer(User{}, question)
		require.Equal(t, "9", result)
	})
	t.Run("should return the addition", func(t *testing.T) {
		question := "¿ Cuanto es 3 + 3 ?"
		a := Arithmetical{}
		result, _, _ := a.Answer(User{}, question)
		require.Equal(t, "6", result)
	})
	t.Run("should return the subsctraction", func(t *testing.T) {
		question := "¿ Cuanto es 4 - 3 ?"
		a := Arithmetical{}
		result, _, _ := a.Answer(User{}, question)
		require.Equal(t, "1", result)
	})
	t.Run("should return false with an unknown operation", func(t *testing.T) {
		question := "¿ Cuanto es 4 % 3 ?"
		a := Arithmetical{}
		_, ok, _ := a.Answer(User{}, question)
		require.False(t, ok)
	})
	t.Run("should return error with an unknown question", func(t *testing.T) {
		question := "¿ Cual es la capital de España ?"
		a := Arithmetical{}
		_, _, err := a.Answer(User{}, question)
		require.Error(t, err)
	})
}
