package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockServer crea un httptest.Server che risponde sempre con success=true e il payload dato.
func mockServer(t *testing.T, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CKANResponse[any]{
			Success: true,
			Result:  payload,
		})
	}))
}

func TestSearchDatasets_Success(t *testing.T) {
	expected := SearchResult{
		Count:   1,
		Results: []Dataset{{ID: "abc", Title: "Test Dataset"}},
	}
	srv := mockServer(t, expected)
	defer srv.Close()

	c := NewWithBase(srv.URL)
	result, err := c.SearchDatasets(context.Background(), "test", "", "", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}
	if result.Results[0].Title != "Test Dataset" {
		t.Errorf("unexpected title: %s", result.Results[0].Title)
	}
}

func TestGetDataset_Success(t *testing.T) {
	srv := mockServer(t, Dataset{ID: "xyz", Title: "My Dataset", Notes: "Una descrizione"})
	defer srv.Close()

	c := NewWithBase(srv.URL)
	d, err := c.GetDataset(context.Background(), "xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.ID != "xyz" {
		t.Errorf("expected id xyz, got %s", d.ID)
	}
}

func TestCKANError_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error": map[string]any{
				"message": "dataset non trovato",
				"__type":  "Not Found Error",
			},
		})
	}))
	defer srv.Close()

	c := NewWithBase(srv.URL)
	_, err := c.GetDataset(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListOrganizations_Success(t *testing.T) {
	srv := mockServer(t, []string{"comune-di-roma", "regione-lombardia"})
	defer srv.Close()

	c := NewWithBase(srv.URL)
	orgs, err := c.ListOrganizations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 2 {
		t.Errorf("expected 2 orgs, got %d", len(orgs))
	}
}
