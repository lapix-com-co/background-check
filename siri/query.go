package siri

import (
	"repoer/siri/reply"
	"time"
)

func (q query) Pull(i *PullInput) ([]*Annotation, error) {
	response, err := q.interactor.RequestPage()
	if err != nil {
		return nil, err
	}

	time.Sleep(q.delay)

	return q.interactor.Annotations(&AnnotationsInput{
		Response: response,
		User: &reply.User{
			Name: reply.Name{
				FirstName:     i.FirstName,
				MiddleName:    i.MiddleName,
				LastName:      i.LastName,
				SecondSurname: i.SecondSurname,
			},
			Document: reply.Document{
				Type:   string(i.DocumentType),
				Number: i.DocumentNumber,
			},
		},
	})
}
