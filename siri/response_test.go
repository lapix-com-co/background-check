package siri

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"repoer/siri/reply"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var annotationsFromFile = []*Annotation{
	&Annotation{
		Type:     "ANTECEDENTES PENALES",
		Category: "Sanciones",
		Fields: []*Field{
			&Field{
				Field: "Sanción",
				Value: "PRISION",
			},
			&Field{
				Field: "Término",
				Value: "2 AÑOS",
			},
			&Field{
				Field: "Clase",
				Value: "PRINCIPAL",
			},
			&Field{
				Field: "Suspendida",
				Value: "SI",
			},
		},
	},
	&Annotation{
		Type:     "ANTECEDENTES PENALES",
		Category: "Sanciones",
		Fields: []*Field{
			&Field{
				Field: "Sanción",
				Value: "INHABILIDAD PARA EL EJERCICIO DE DERECHOS Y FUNCIONES PUBLICAS",
			},
			&Field{
				Field: "Término",
				Value: "2 AÑOS",
			},
			&Field{
				Field: "Clase",
				Value: "ACCESORIA",
			},
			&Field{
				Field: "Suspendida",
			},
		},
	},
	&Annotation{
		Type:     "ANTECEDENTES PENALES",
		Category: "Delitos",
		Fields: []*Field{
			&Field{
				Field: "Descripción del Delito",
				Value: "PORTE ILEGAL DE ARMAS DE FUEGO (LEY 599 DE 2000)",
			},
		},
	},
	&Annotation{
		Type:     "ANTECEDENTES PENALES",
		Category: "Instancias",
		Fields: []*Field{
			&Field{
				Field: "Nombre",
				Value: "PRIMERA",
			},
			&Field{
				Field: "Autoridad",
				Value: "JUZGADO 6 PENAL DEL CIRCUITO DE CONOCIMIENTO - BUCARAMANGA (SANTANDER)",
			},
			&Field{
				Field: "Fecha providencia",
				Value: "26/04/2018",
			},
			&Field{
				Field: "fecha efecto Juridicos",
				Value: "26/04/2018",
			},
		},
	},
	&Annotation{
		Type: "INHABILIDADES",
		Fields: []*Field{
			&Field{
				Field: "SIRI",
				Value: "201171721",
			},
			&Field{
				Field: "Módulo",
				Value: "PENAL",
			},
			&Field{
				Field: "Inhabilidad legal",
				Value: "INHABILIDAD PARA CONTRATAR CON EL ESTADO LEY 80 ART 8 LIT. D",
			},
			&Field{
				Field: "Fecha de inicio",
				Value: "26/04/2018",
			},
			&Field{
				Field: "Fecha fin",
				Value: "25/04/2023",
			},
		},
	},
}

type stubFinder struct {
	strategy reply.Strategy
}

func (s stubFinder) Strategy(reply.User, string) reply.Strategy { return s.strategy }

type stubEmptyFinder struct{}

func (s stubEmptyFinder) Strategy(reply.User, string) reply.Strategy { return nil }

type stubStrategy struct{}

func (s stubStrategy) Is(reply.User, string) bool                      { return true }
func (s stubStrategy) Answer(reply.User, string) (string, bool, error) { return "my-answer", true, nil }

type stubUnresolvedStrategy struct{}

func (s stubUnresolvedStrategy) Is(reply.User, string) bool { return true }
func (s stubUnresolvedStrategy) Answer(reply.User, string) (string, bool, error) {
	return "", false, nil
}

var htmlFileResponse = func(filename string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}

var doesNotHaveAnnotations = func(w http.ResponseWriter, r *http.Request) {
	htmlFileResponse("testdata/annotations.html")(w, r)
}

var haveAnnotations = func(w http.ResponseWriter, r *http.Request) {
	htmlFileResponse("testdata/negative-annotations.html")(w, r)
}

var wrongAnswer = func(w http.ResponseWriter, r *http.Request) {
	htmlFileResponse("testdata/wrong-answer.html")(w, r)
}

var dataAssert = func(t *testing.T, i *AnnotationsInput) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		// Cookies
		cookies := r.Cookies()
		if len(cookies) == 0 {
			t.Error("expect at least one cookie found cero")
		} else {
			assert.Equal(t, "ASP.NET_SessionId", cookies[0].Name)
			assert.Equal(t, "bgnyi0asjvigvh55t2gfrk3x", cookies[0].Value)
		}

		// Headers
		assert.Equal(t, "no-cache", r.Header["Cache-Control"][0])
		require.Equal(t, "https://www.procuraduria.gov.co", r.Header["Origin"][0])
		require.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36", r.Header["User-Agent"][0])
		require.Equal(t, "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9", r.Header["Accept"][0])
		require.Equal(t, "document", r.Header["Sec-Fetch-Dest"][0])
		require.Equal(t, "navigate", r.Header["Sec-Fetch-Mode"][0])
		require.Equal(t, "?1", r.Header["Sec-Fetch-User"][0])
		require.Equal(t, "https://www.procuraduria.gov.co/CertWEB/Certificado.aspx?tpo=1", r.Header["Referer"][0])
		require.Equal(t, "es-ES,es;q=0.9,en;q=0.8", r.Header["Accept-Language"][0])
		require.Equalf(t, "1", r.Header["Dnt"][0], "expect dnt value = 1, but got %v", r.Header["Dnt"])

		// Form Content
		c, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			panic("could not read the request body")
		}

		v, err := url.ParseQuery(string(c))
		if err != nil {
			panic("could not parse the request body")
		}

		require.Equal(t, "", v.Get("__EVENTTARGET"))
		require.Equal(t, "", v.Get("__EVENTARGUMENT"))
		require.Equal(t, "/wEPDwUJLTU5NTU5MDcyDxYCHgpJZFByZWd1bnRhBQIxMxYCAgMPZBYMAgEPDxYCHgRUZXh0BRhDb25zdWx0YSBkZSBhbnRlY2VkZW50ZXNkZAINDxYCHgdWaXNpYmxlaBYEAgEPZBYCAgEPEGRkFgFmZAIDD2QWAgIBDxBkZBYAZAIPDw8WAh8BBSTCvyBDdWFsIGVzIGxhIENhcGl0YWwgZGVsIEF0bGFudGljbz9kZAIYDw8WAh8CaGRkAiAPDxYCHwEFR0ZlY2hhIGRlIGNvbnN1bHRhOiBzw6FiYWRvLCBtYXJ6byAxNCwgMjAyMCAtIEhvcmEgZGUgY29uc3VsdGE6IDEyOjQwOjE5ZGQCJA8PFgIfAQUHVi4wLjAuNGRkGAEFHl9fQ29udHJvbHNSZXF1aXJlUG9zdEJhY2tLZXlfXxYBBQxJbWFnZUJ1dHRvbjEGteAxMJen2KgR0FzO7YGpaAYG/g==", v.Get("__VIEWSTATE"))
		require.Equal(t, "D8335CE7", v.Get("__VIEWSTATEGENERATOR"))
		require.Equal(t, "/wEWCgK3hK+PCwL8kK+TAQLwkOOQAQLvkOOQAQLxkOOQAQL0kOOQAQK8zP8SAtLCmdMIAsimk6ECApWrsq8ISuNT/p3oBQjSn0uui6yv9yX6YJg=", v.Get("__EVENTVALIDATION"))
		require.Equal(t, "1", v.Get("ddlTipoID"))
		require.Equal(t, "1029788261", v.Get("txtNumID"))
		require.Equal(t, "my-answer", v.Get("txtRespuestaPregunta"))
		require.Equal(t, "Consultar", v.Get("btnConsultar"))
		require.Equal(t, "063515BAFE28EF92C9E1E74EDAB80FF5", v.Get("10E9410D"))

		doesNotHaveAnnotations(w, r)
	}
}

func TestResponse_Annotations(t *testing.T) {
	t.Run("should return the customer annotations", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(serverResponse))
		defer server.Close()
		handler := &handler{url: server.URL, factory: &stubFinder{&stubStrategy{}}}
		want := []*Annotation{}

		got, _ := handler.Annotations(&AnnotationsInput{
			User: &reply.User{reply.Name{}, reply.Document{}},
			Response: &response{
				View:    &view{},
				Event:   &event{},
				KeyCode: &key{},
				Session: &session{Question: "any-question"},
			},
		})

		require.Equal(t, want, got)
	})

	t.Run("should return error if the question is unknown", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(serverResponse))
		defer server.Close()
		handler := &handler{
			url:     server.URL,
			factory: &stubEmptyFinder{},
		}

		_, err := handler.Annotations(&AnnotationsInput{
			User: &reply.User{},
			Response: &response{
				Session: &session{
					Question: "any-question",
				},
			},
		})

		require.Error(t, err)
	})

	t.Run("should return error if the question does not have an answer", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(serverResponse))
		defer server.Close()
		handler := &handler{
			url:     server.URL,
			factory: &stubFinder{&stubUnresolvedStrategy{}},
		}

		_, err := handler.Annotations(&AnnotationsInput{
			User: &reply.User{},
			Response: &response{
				Session: &session{
					Question: "any-question",
				},
			},
		})

		require.Error(t, err)
	})

	t.Run("should return a list of annotations", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(haveAnnotations))
		defer server.Close()
		handler := &handler{
			url:     server.URL,
			factory: &stubFinder{&stubStrategy{}},
		}

		got, _ := handler.Annotations(&AnnotationsInput{
			User: &reply.User{
				Document: reply.Document{
					Type:   "1",
					Number: "1029788261",
				},
			},
			Response: &response{
				View:    &view{},
				Event:   &event{},
				KeyCode: &key{},
				Session: &session{Question: "any-question"},
			},
		})

		want := annotationsFromFile

		require.Equal(t, want, got)
	})

	t.Run("should return an error is the given answer is not valid", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(wrongAnswer))
		defer server.Close()
		handler := &handler{
			url:     server.URL,
			factory: &stubFinder{&stubStrategy{}},
		}

		_, err := handler.Annotations(&AnnotationsInput{
			User: &reply.User{
				Document: reply.Document{
					Type:   "1",
					Number: "1029788261",
				},
			},
			Response: &response{
				View:    &view{},
				Event:   &event{},
				KeyCode: &key{},
				Session: &session{Question: "any-question"},
			},
		})

		require.Error(t, err)
	})

	t.Run("should send the proper header, cookies and form content", func(t *testing.T) {
		input := &AnnotationsInput{
			User: &reply.User{
				Name: reply.Name{},
				Document: reply.Document{
					Type:   "1",
					Number: "1029788261",
				},
			},
			Response: &response{
				View: &view{
					State:     "/wEPDwUJLTU5NTU5MDcyDxYCHgpJZFByZWd1bnRhBQIxMxYCAgMPZBYMAgEPDxYCHgRUZXh0BRhDb25zdWx0YSBkZSBhbnRlY2VkZW50ZXNkZAINDxYCHgdWaXNpYmxlaBYEAgEPZBYCAgEPEGRkFgFmZAIDD2QWAgIBDxBkZBYAZAIPDw8WAh8BBSTCvyBDdWFsIGVzIGxhIENhcGl0YWwgZGVsIEF0bGFudGljbz9kZAIYDw8WAh8CaGRkAiAPDxYCHwEFR0ZlY2hhIGRlIGNvbnN1bHRhOiBzw6FiYWRvLCBtYXJ6byAxNCwgMjAyMCAtIEhvcmEgZGUgY29uc3VsdGE6IDEyOjQwOjE5ZGQCJA8PFgIfAQUHVi4wLjAuNGRkGAEFHl9fQ29udHJvbHNSZXF1aXJlUG9zdEJhY2tLZXlfXxYBBQxJbWFnZUJ1dHRvbjEGteAxMJen2KgR0FzO7YGpaAYG/g==",
					Generator: "D8335CE7",
				},
				Event: &event{
					Argument:   "",
					Target:     "",
					Validation: "/wEWCgK3hK+PCwL8kK+TAQLwkOOQAQLvkOOQAQLxkOOQAQL0kOOQAQK8zP8SAtLCmdMIAsimk6ECApWrsq8ISuNT/p3oBQjSn0uui6yv9yX6YJg=",
				},
				KeyCode: &key{
					Name:  "10E9410D",
					Value: "063515BAFE28EF92C9E1E74EDAB80FF5",
				},
				Session: &session{
					Question: "¿ Cual es la capital de Atlantico?",
					Cookies: map[string]string{
						"ASP.NET_SessionId": "bgnyi0asjvigvh55t2gfrk3x",
					},
				},
			},
		}
		server := httptest.NewServer(http.HandlerFunc(dataAssert(t, input)))
		defer server.Close()
		handler := &handler{url: server.URL, factory: &stubFinder{&stubStrategy{}}}

		got, err := handler.Annotations(input)

		require.NoError(t, err)
		require.Equal(t, []*Annotation{}, got)
	})
}
