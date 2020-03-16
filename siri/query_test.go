package siri

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stubInteractor struct{}

func (s stubInteractor) RequestPage() (*response, error) { return &response{}, nil }
func (s stubInteractor) Annotations(*AnnotationsInput) ([]*Annotation, error) {
	return []*Annotation{}, nil
}

func TestQuery_Pull(t *testing.T) {
	t.Run("should ask for the user annotations", func(t *testing.T) {
		q := query{
			delay:      time.Second * 0,
			interactor: &stubInteractor{},
		}

		r, err := q.Pull(&PullInput{
			FirstName:      "John",
			MiddleName:     "Alexander",
			LastName:       "Doe",
			SecondSurname:  "Bell",
			DocumentType:   DNI,
			DocumentNumber: "1010778260",
		})

		require.NoError(t, err)
		require.Empty(t, r)
	})
}
