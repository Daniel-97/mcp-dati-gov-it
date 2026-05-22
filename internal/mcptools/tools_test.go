package mcptools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func newMockClient(t *testing.T, payload any) *client.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"success": true, "result": payload})
	}))
	t.Cleanup(srv.Close)
	return client.NewWithBase(srv.URL)
}

func TestSearchTool_ReturnsJSON(t *testing.T) {
	c := newMockClient(t, map[string]any{
		"count":   1,
		"results": []map[string]any{{"id": "1", "title": "Test"}},
	})
	h := NewHandlers(c)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"query": "test"}

	result, err := h.SearchDatasets(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected MCP error: %v", result.Content)
	}
	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}
}

func TestGetDatasetTool_ReturnsDataset(t *testing.T) {
	c := newMockClient(t, map[string]any{"id": "abc", "title": "My Dataset"})
	h := NewHandlers(c)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"id": "abc"}

	result, err := h.GetDataset(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected MCP error")
	}
}

func TestGetDatasetTool_MissingID(t *testing.T) {
	c := newMockClient(t, nil)
	h := NewHandlers(c)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := h.GetDataset(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected MCP error for missing id")
	}
}

func TestListOrganizationsTool_ReturnsList(t *testing.T) {
	c := newMockClient(t, []string{"comune-di-roma", "regione-lombardia"})
	h := NewHandlers(c)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := h.ListOrganizations(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected MCP error")
	}
}
