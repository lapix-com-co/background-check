package reply

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument_Is(t *testing.T) {
	t.Run("should return the question as valid", func(t *testing.T) {
		d := DocumentCharts{}
		result := d.Is(User{}, "¿Escriba los dos ultimos digitos del documento a consultar?")
		require.True(t, result)
	})

	t.Run("should return the question as valid", func(t *testing.T) {
		d := DocumentCharts{}
		result := d.Is(User{}, "¿Escriba las dos primeras letras del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.True(t, result)
	})

	t.Run("should return the question as valid with a single value", func(t *testing.T) {
		d := DocumentCharts{}
		result := d.Is(User{}, "¿Escriba el ultimo digito del documento a consultar?")
		require.True(t, result)
	})

	t.Run("should return the question as valid with a the first digits", func(t *testing.T) {
		d := DocumentCharts{}
		result := d.Is(User{}, "¿Escriba los tres primeros digitos del documento a consultar?")
		require.True(t, result, "asking the first three characteres should be valid")
	})

	t.Run("should return the question as invalid", func(t *testing.T) {
		d := DocumentCharts{}
		result := d.Is(User{}, "¿Escriba la cantidad de letras del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.False(t, result)
	})

}

func TestDocument_Answer(t *testing.T) {
	user := User{
		Name: Name{
			FirstName: "John",
		},
		Document: Document{
			Number: "123411",
		},
	}

	t.Run("should return the last two numbers from the document number", func(t *testing.T) {

		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba los dos ultimos digitos del documento a consultar?")
		require.Equal(t, "11", a)
	})

	t.Run("should return the last three numbers from the document number", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba los tres ultimos digitos del documento a consultar?")
		require.Equal(t, "411", a)
	})

	t.Run("should return error with an unknown question", func(t *testing.T) {
		d := DocumentCharts{}
		_, _, err := d.Answer(user, "¿Escriba los nueve primeros números del documento a consultar?")
		require.Error(t, err)
	})

	t.Run("should return the first two characteres of the first name", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba las dos primeras letras del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.Equal(t, "Jo", a)
	})
	t.Run("should return the first three characteres of the first name", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba las tres primeras letras del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.Equal(t, "Joh", a)
	})
	t.Run("should return the first characteres of the first name", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba la primer letra del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.Equal(t, "J", a)
	})

	t.Run("should return the last two characteres of the first name", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba las dos ultimas letras del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.Equal(t, "hn", a)
	})
	t.Run("should return the last three characteres of the first name", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba las tres ultimas letras del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.Equal(t, "ohn", a)
	})
	t.Run("should return the last characteres of the first name", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba la ultima letra del primer nombre de la persona a la cual esta expidiendo el certificado?")
		require.Equal(t, "n", a)
	})
	t.Run("should return the first three digts from the document number", func(t *testing.T) {
		d := DocumentCharts{}
		a, _, _ := d.Answer(user, "¿Escriba los tres primeros digitos del documento a consultar?")
		require.Equal(t, "123", a)
	})
}
