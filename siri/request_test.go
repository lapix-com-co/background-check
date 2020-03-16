package siri

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var serverResponse = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "private")
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Set-Cookie", "ASP.NET_SessionId=s2w31o55ha20p355d2mf04ut; path=/; HttpOnly")
	w.Header().Add("Set-Cookie", "awaf-sid=7c264e320116d59f; path=/; max-age=3600")

	file, err := os.Open("testdata/form-page.html")
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

var notFoundResponse = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<html>Does not matter</html>"))
}

var invalidResponse = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html>Does not matter</html>"))
}

func TestPullValues(t *testing.T) {
	t.Run("retrieve the data from the given URL", func(t *testing.T) {
		want := &response{
			KeyCode: &key{
				Name:  "10E9410D",
				Value: "52E1353EFA68D09E7661B409857CE1F0",
			},
			Event: &event{
				Argument:   "",
				Target:     "",
				Validation: "/wEWCgLLt9MRAvyQr5MBAvCQ45ABAu+Q45ABAvGQ45ABAvSQ45ABArzM/xIC0sKZ0wgCyKaToQIClauyrwiGJpyEdTLpX1p4nOZKqr9HdpV16Q==",
			},
			View: &view{
				State:     "/wEPDwUJLTU5NTU5MDcyDxYCHgpJZFByZWd1bnRhBQIxMhYCAgMPZBYMAgEPDxYCHgRUZXh0BRhDb25zdWx0YSBkZSBhbnRlY2VkZW50ZXNkZAINDxYCHgdWaXNpYmxlaBYEAgEPZBYCAgEPEGRkFgFmZAIDD2QWAgIBDxBkZBYAZAIPDw8WAh8BBS/CvyBDdWFsIGVzIGxhIENhcGl0YWwgZGUgQW50aW9xdWlhIChzaW4gdGlsZGUpP2RkAhgPDxYCHwJoZGQCIA8PFgIfAQVGRmVjaGEgZGUgY29uc3VsdGE6IG1hcnRlcywgbWFyem8gMTAsIDIwMjAgLSBIb3JhIGRlIGNvbnN1bHRhOiAyMzoxMjoyOGRkAiQPDxYCHwEFB1YuMC4wLjRkZBgBBR5fX0NvbnRyb2xzUmVxdWlyZVBvc3RCYWNrS2V5X18WAQUMSW1hZ2VCdXR0b24xZaGi0bBJ6VWq9y2p/xlIIKaLITE=",
				Generator: "D8335CE7",
			},
			Session: &session{
				Question: "¿ Cual es la Capital de Antioquia (sin tilde)?",
				Cookies: map[string]string{
					"ASP.NET_SessionId": "s2w31o55ha20p355d2mf04ut",
					"awaf-sid":          "7c264e320116d59f",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(serverResponse))
		defer server.Close()
		handler := &handler{url: server.URL}

		got, err := handler.RequestPage()

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("should returns an error if the server response with an invalid status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(notFoundResponse))
		defer server.Close()
		handler := &handler{url: server.URL}
		_, err := handler.RequestPage()

		if !errors.Is(err, ErrInvalidServerState) {
			t.Errorf("expect invalid server state but got = %v", err)
		}
	})

	t.Run("should return an error if the response does not have the valid fields", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(invalidResponse))
		defer server.Close()
		handler := &handler{url: server.URL}
		_, err := handler.RequestPage()

		if !errors.Is(err, ErrInvalidHTMLTree) {
			t.Errorf("expect invalid html tree but got = %v", err)
		}
	})
}
