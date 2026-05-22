package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
)

// newMockClient crea un client che risponde sempre con il payload dato.
func newMockClient(t *testing.T, payload any) *client.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"success": true, "result": payload})
	}))
	t.Cleanup(srv.Close)
	return client.NewWithBase(srv.URL)
}

// newErrorMockClient crea un client che risponde sempre con success=false e il messaggio dato.
func newErrorMockClient(t *testing.T, message string) *client.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   map[string]any{"message": message, "__type": "Not Found Error"},
		})
	}))
	t.Cleanup(srv.Close)
	return client.NewWithBase(srv.URL)
}

func TestRunSearch_ReturnsResults(t *testing.T) {
	c := newMockClient(t, map[string]any{
		"count":   2,
		"results": []map[string]any{{"id": "1", "title": "A"}, {"id": "2", "title": "B"}},
	})

	out, err := RunSearch(context.Background(), c, SearchOptions{Query: "test", Rows: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Count != 2 {
		t.Errorf("expected count 2, got %d", out.Count)
	}
	if len(out.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(out.Results))
	}
}

func TestRunSearch_PropagatesError(t *testing.T) {
	c := newErrorMockClient(t, "server error")

	_, err := RunSearch(context.Background(), c, SearchOptions{Query: "test", Rows: 10})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewSearchCmd_AgentOutput(t *testing.T) {
	c := newMockClient(t, map[string]any{
		"count":   1,
		"results": []map[string]any{{"id": "1", "title": "Test"}},
	})

	cmd := NewSearchCmd(c)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"test", "--agent"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd error: %v", err)
	}

	var out SearchOutput
	if err := json.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatalf("decode agent output: %v", err)
	}
	if out.Count != 1 {
		t.Errorf("expected count 1, got %d", out.Count)
	}
}

func TestNewSearchCmd_HumanOutput(t *testing.T) {
	c := newMockClient(t, map[string]any{
		"count":   1,
		"results": []map[string]any{{"id": "ds-1", "title": "Dataset Ambiente"}},
	})

	cmd := NewSearchCmd(c)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"ambiente"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd error: %v", err)
	}
	if buf.String() == "" {
		t.Error("expected non-empty human output")
	}
}
