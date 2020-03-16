package siri

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (s *handler) RequestPage() (*response, error) {
	request, _ := s.buildFormRequest()
	response, _ := http.DefaultClient.Do(request)
	return s.handlerFormResponse(response)
}

func (s *handler) buildFormRequest() (*http.Request, error) {
	r, err := http.NewRequest("GET", s.url, strings.NewReader(""))
	setFakeHeaders(r)
	return r, err
}

func (s *handler) handlerFormResponse(httpResponse *http.Response) (*response, error) {
	var err error
	var result = &response{
		View:  &view{},
		Event: &event{},
		KeyCode: &key{
			Name: "10E9410D",
		},
		Session: &session{
			Cookies: map[string]string{},
		},
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returns %d status code %w", httpResponse.StatusCode, ErrInvalidServerState)
	}

	defer httpResponse.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(httpResponse.Body)

	if result.Event.Target, err = value(doc, "#__EVENTTARGET"); err != nil {
		return nil, err
	}
	if result.Event.Argument, err = value(doc, "#__EVENTARGUMENT"); err != nil {
		return nil, err
	}
	if result.Event.Validation, err = value(doc, "#__EVENTVALIDATION"); err != nil {
		return nil, err
	}

	if result.KeyCode.Value, err = value(doc, `input[name="10E9410D"]`); err != nil {
		return nil, err
	}

	if result.View.State, err = value(doc, "#__VIEWSTATE"); err != nil {
		return nil, err
	}
	if result.View.Generator, err = value(doc, "#__VIEWSTATEGENERATOR"); err != nil {
		return nil, err
	}

	result.Session.Question = doc.Find("#lblPregunta").First().Text()

	for _, cookie := range httpResponse.Cookies() {
		result.Session.Cookies[cookie.Name] = cookie.Value
	}

	return result, nil
}

func value(doc *goquery.Document, selector string) (string, error) {
	val, exists := doc.Find(selector).First().Attr("value")
	if exists {
		return val, nil
	}

	return "", fmt.Errorf("the given selector '%s' does not have a valid value %w", selector, ErrInvalidHTMLTree)
}

func setFakeHeaders(r *http.Request) {
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Origin", "https://www.procuraduria.gov.co")
	r.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Referer", "https://www.procuraduria.gov.co/CertWEB/Certificado.aspx?tpo=1")
	r.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
}