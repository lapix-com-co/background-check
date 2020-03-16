package siri

import (
	"errors"
	"repoer/siri/reply"
	"time"
)

var (
	ErrInvalidServerState = errors.New("invalid server response")
	ErrInvalidHTMLTree    = errors.New("the html tree is not valid")
	ErrUnknownQuestion    = errors.New("unknown question")
	ErrInvalidAnswer      = errors.New("invalid answer")
)

type AnnotationsAware interface {
	Pull(*PullInput) ([]*Annotation, error)
}

type DocumentType string

const (
	DNI         DocumentType = "1"
	PEP         DocumentType = "0"
	NIT         DocumentType = "2"
	ExternalDNI DocumentType = "3"
)

type Option func(*query)

func Delay(t time.Duration) Option {
	return func(q *query) {
		q.delay = t
	}
}

// NewQuery returns a new interactor with a default delay time of four seconds between requests.
func NewQuery(i ...Option) *query {
	q := &query{
		delay:      time.Second * 4,
		interactor: newHandler(),
	}

	for _, op := range i {
		op(q)
	}

	return q
}

type PullInput struct {
	FirstName      string
	MiddleName     string
	LastName       string
	SecondSurname  string
	DocumentType   DocumentType
	DocumentNumber string
}

type Annotation struct {
	Type     string
	Category string
	Fields   []*Field
}

type Field struct {
	Field string
	Value string
}

type AnnotationsInput struct {
	Response *response
	User     *reply.User
}

type interactor interface {
	RequestPage() (*response, error)
	Annotations(*AnnotationsInput) ([]*Annotation, error)
}

type view struct {
	State     string
	Generator string
}

type event struct {
	Argument   string
	Target     string
	Validation string
}

type key struct {
	Name  string
	Value string
}

type response struct {
	View    *view
	Event   *event
	KeyCode *key
	Session *session
}

type session struct {
	Question string
	Cookies  map[string]string
}

type handler struct {
	url     string
	factory reply.StrategyFinder
}

type query struct {
	delay      time.Duration
	interactor interactor
}

func newHandler() *handler {
	return &handler{
		url:     "https://www.procuraduria.gov.co/CertWEB/Certificado.aspx",
		factory: reply.NewFactory(),
	}
}
