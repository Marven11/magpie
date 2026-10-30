package mcpauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// discover reads each well-known document into a structure of its own: one
// that fails to decode must not leave the fields it did read behind for the
// next document to complete. Both authorization server documents here name
// one endpoint each, and the first also a field of the wrong type (the
// challenge methods as a string), so it doesn't decode: the sign-in must not
// be made out of the two of them.
func TestDiscoverDoesNotMixDocuments(t *testing.T) {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	defer srv.Close()
	issuer := srv.URL + "/as"
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+srv.URL+`/.well-known/oauth-protected-resource/mcp"`)
		w.WriteHeader(http.StatusUnauthorized)
	})
	mux.HandleFunc("/.well-known/oauth-protected-resource/mcp", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"authorization_servers": []string{issuer}})
	})
	mux.HandleFunc("/.well-known/oauth-authorization-server/as", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"issuer":"`+issuer+`","token_endpoint":"`+issuer+`/first-token","code_challenge_methods_supported":"S256"}`)
	})
	mux.HandleFunc("/.well-known/openid-configuration/as", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"authorization_endpoint":"`+issuer+`/second-authorize"}`)
	})
	if m, err := discover(context.Background(), srv.URL+"/mcp"); err == nil {
		t.Fatalf("discover made a sign-in out of two documents: %+v", m)
	}
}
