package gui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A custom provider on OpenAI's Responses API, saved by the editor with
// "Searches the web by itself" ticked and its models typed in because it
// has no model list, is offered in Settings › Web search's searcher
// picker with those models (axelburks on Discord): it was named among the
// relays to pick there, with nothing to pick.
func TestSearcherPickerOffersResponsesRelay(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	srv := httptest.NewServer(http.NotFoundHandler()) // no /models
	t.Cleanup(srv.Close)
	mux := http.NewServeMux()
	providerRoutes(mux, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("POST", "/api/provider/save", strings.NewReader(
		`{"new":true,"name":"Relay","responses":"`+srv.URL+`/v1","baseAPI":"responses","key":"k","searches":true,"icon":"generic","models":["gpt-5.5"]}`)))
	if w.Code != 200 {
		t.Fatalf("save: %d %s", w.Code, w.Body)
	}
	var s settingsJSON
	searchState(&s)
	if len(s.SearchChoices) != 1 || s.SearchChoices[0].Name != "Relay" {
		t.Fatalf("choices = %+v", s.SearchChoices)
	}
	c := s.SearchChoices[0]
	if len(c.Models) != 1 || c.Models[0].ID != c.ID+"/gpt-5.5" || c.Small != "gpt-5.5" {
		t.Errorf("choice = %+v", c)
	}
	if len(s.SearchRelays) != 1 || len(s.SearchLeftOut) != 0 {
		t.Errorf("relays %v, left out %v", s.SearchRelays, s.SearchLeftOut)
	}
}
