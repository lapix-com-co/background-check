package siri

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/atom"
)

// Annotations send a form to siri and returns the annotations
func (s *handler) Annotations(i *AnnotationsInput) ([]*Annotation, error) {
	answer, err := s.findAnswer(i)
	if err != nil {
		return nil, err
	}

	request, _ := s.buildRequest(i, answer)
	response, _ := http.DefaultClient.Do(request)

	return s.handleResponse(response, answer, i.Response.Session.Question)
}

func (s *handler) findAnswer(i *AnnotationsInput) (string, error) {
	var strategy = s.factory.Strategy(*i.User, i.Response.Session.Question)

	if strategy == nil {
		return "", fmt.Errorf("could not found a strategy for this question: %v. %w", i.Response.Session.Question, ErrUnknownQuestion)
	}

	answer, ok, err := strategy.Answer(*i.User, i.Response.Session.Question)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("could not found an answer for this question: %v. %w", i.Response.Session.Question, ErrUnknownQuestion)
	}

	return answer, nil
}

func (s *handler) buildRequest(i *AnnotationsInput, answer string) (*http.Request, error) {
	requestContent := buildFormContent(i, answer)

	r, err := http.NewRequest("POST", s.url, strings.NewReader(requestContent))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Content-Length", strconv.Itoa(len(requestContent)))
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("dnt", "1")
	setFakeHeaders(r)

	r.AddCookie(&http.Cookie{Name: "ASP.NET_SessionId", Value: i.Response.Session.Cookies["ASP.NET_SessionId"]})
	return r, nil
}

func (s *handler) handleResponse(response *http.Response, answer, question string) ([]*Annotation, error) {
	var result = make([]*Annotation, 0)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returns %d status code %w", response.StatusCode, ErrInvalidServerState)
	}

	doc, _ := goquery.NewDocumentFromReader(response.Body)

	v := doc.Find("#ValidationSummary1").First().Text()

	if strings.TrimSpace(v) == "El valor ingresado para la respuesta no responde a la pregunta." {
		return nil, fmt.Errorf(`the given answer "%s" for the question "%s" is not valid: %w`, answer, question, ErrInvalidAnswer)
	}

	v = doc.Find("#divSec h2").First().Text()

	if v != "El ciudadano no presenta antecedentes" {
		result = extractAnnotations(doc)
	}

	return result, nil
}

func buildFormContent(i *AnnotationsInput, answer string) string {
	var form = url.Values{}
	form.Set("__EVENTTARGET", i.Response.Event.Target)
	form.Set("__EVENTARGUMENT", i.Response.Event.Argument)
	form.Set("__EVENTVALIDATION", i.Response.Event.Validation)
	form.Set("__VIEWSTATE", i.Response.View.State)
	form.Set("__VIEWSTATEGENERATOR", i.Response.View.Generator)
	form.Set("ddlTipoID", i.User.Document.Type)
	form.Set("txtNumID", i.User.Document.Number)
	form.Set("txtRespuestaPregunta", answer)
	form.Set("btnConsultar", "Consultar")
	form.Set("10E9410D", "063515BAFE28EF92C9E1E74EDAB80FF5")

	return form.Encode()
}

// extractAnnotations created the annotations elements
func extractAnnotations(doc *goquery.Document) []*Annotation {
	var o = make([]*Annotation, 0)

	doc.Find(".SeccionAnt").Each(func(i int, selection *goquery.Selection) {
		annotationType := selection.Find("h2").First().Text()
		session := selection.Find(".SessionNumSiri").First()

		if session.Length() == 0 {
			session = selection
		}

		session.Find("table").Each(func(k int, table *goquery.Selection) {
			category := ""
			prev := table.Prev()

			if prev.Get(0).DataAtom == atom.H3 {
				category = normalizeText(prev.Text())
			}

			var headers []string

			table.Find("th").Each(func(v int, head *goquery.Selection) {
				headers = append(headers, normalizeText(head.Text()))
			})

			table.Find("tr").Each(func(v int, row *goquery.Selection) {
				dataElements := row.Find("td")
				if dataElements.Length() == 0 {
					return
				}

				an := &Annotation{
					Type:     annotationType,
					Category: category,
					Fields:   []*Field{},
				}

				o = append(o, an)

				dataElements.Each(func(z int, data *goquery.Selection) {
					an.Fields = append(an.Fields, &Field{
						Field: headers[z],
						Value: normalizeText(data.Text()),
					})
				})
			})
		})
	})

	return o
}

func normalizeText(i string) string {
	return strings.TrimSpace(i)
}
