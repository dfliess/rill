// Package odootest is a fake Odoo JSON-2 server for tests of the odoo driver.
//
// It serves only what the driver calls: fields_get, search_read and context_get. search_read
// honours the keyset term ["id", ">", n] and the limit, and ignores every other domain term;
// tests assert on the recorded domain instead. Every request is recorded, including rejected ones.
package odootest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
)

// Model is a fake Odoo model: field name to Odoo type, and records sorted by id.
type Model struct {
	Fields  map[string]string
	Records []map[string]any
}

// Call is one recorded request. Numbers in Body are json.Number.
type Call struct {
	Model  string
	Method string
	Header http.Header
	Body   map[string]any
}

// Server is a running fake Odoo.
type Server struct {
	*httptest.Server
	APIKey string

	mu     sync.Mutex
	models map[string]Model
	calls  []Call
}

// NewServer starts a fake Odoo that accepts apiKey and serves models. It stops with the test.
func NewServer(t testing.TB, apiKey string, models map[string]Model) *Server {
	s := &Server{APIKey: apiKey, models: models}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

// Calls returns every request received so far.
func (s *Server) Calls() []Call {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.calls)
}

// CallsTo returns the requests received for one method.
func (s *Server) CallsTo(method string) []Call {
	var out []Call
	for _, c := range s.Calls() {
		if c.Method == method {
			out = append(out, c)
		}
	}
	return out
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/json/2/"), "/")
	if r.Method != http.MethodPost || len(path) != 2 {
		writeError(w, http.StatusNotFound, "werkzeug.exceptions.NotFound", "404 Not Found: Did you mean POST /json/2/<model>/<method>?")
		return
	}
	var body map[string]any
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "werkzeug.exceptions.BadRequest", err.Error())
		return
	}
	s.mu.Lock()
	s.calls = append(s.calls, Call{Model: path[0], Method: path[1], Header: r.Header.Clone(), Body: body})
	model, ok := s.models[path[0]]
	s.mu.Unlock()

	if r.Header.Get("Authorization") != "bearer "+s.APIKey {
		writeError(w, http.StatusUnauthorized, "werkzeug.exceptions.Unauthorized", "401 Unauthorized: Invalid apikey")
		return
	}
	if path[1] == "context_get" {
		writeJSON(w, map[string]any{"lang": "en_US", "tz": "Europe/Madrid"})
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "werkzeug.exceptions.NotFound", fmt.Sprintf("404 Not Found: the model %q does not exist", path[0]))
		return
	}

	switch path[1] {
	case "fields_get":
		out := map[string]map[string]any{}
		for _, f := range strs(body["allfields"]) {
			if typ, ok := model.Fields[f]; ok {
				out[f] = map[string]any{"type": typ}
			}
		}
		writeJSON(w, out)
	case "search_read":
		after, limit := keysetFloor(body["domain"]), integer(body["limit"])
		out := []map[string]any{}
		for _, rec := range model.Records {
			if integer(rec["id"]) <= after {
				continue
			}
			row := map[string]any{}
			for _, f := range strs(body["fields"]) {
				row[f] = rec[f]
			}
			out = append(out, row)
			if limit > 0 && int64(len(out)) == limit {
				break
			}
		}
		writeJSON(w, out)
	default:
		writeError(w, http.StatusNotFound, "werkzeug.exceptions.NotFound", "404 Not Found: method not served by the fake")
	}
}

// keysetFloor returns n from the last ["id", ">", n] term of a domain, or 0.
func keysetFloor(domain any) int64 {
	terms, _ := domain.([]any)
	for i := len(terms) - 1; i >= 0; i-- {
		if t, ok := terms[i].([]any); ok && len(t) == 3 && t[0] == "id" && t[1] == ">" {
			return integer(t[2])
		}
	}
	return 0
}

func integer(v any) int64 {
	switch n := v.(type) {
	case json.Number:
		i, _ := n.Int64()
		return i
	case int:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	}
	return 0
}

func strs(v any) []string {
	items, _ := v.([]any)
	out := make([]string, 0, len(items))
	for _, it := range items {
		if s, ok := it.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func writeJSON(w http.ResponseWriter, v any) {
	writeJSONStatus(w, http.StatusOK, v)
}

// writeError answers like Odoo does, traceback included, so tests can check it never leaks.
func writeError(w http.ResponseWriter, status int, name, message string) {
	writeJSONStatus(w, status, map[string]any{
		"name": name, "message": message, "arguments": []any{message}, "context": map[string]any{},
		"debug": "Traceback (most recent call last):\n  File \"/opt/odoo/odoo/http.py\" ...",
	})
}

// writeJSONStatus encodes v before sending the status, so a value that cannot be encoded is answered as a 500.
func writeJSONStatus(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("odootest: cannot encode response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
