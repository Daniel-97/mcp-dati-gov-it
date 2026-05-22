package commands

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunDownload_Success(t *testing.T) {
	// Server che serve il contenuto del file
	fileSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("col1,col2\n1,2\n"))
	}))
	defer fileSrv.Close()

	// Client che risponde con il dataset contenente la risorsa
	c := newMockClient(t, map[string]any{
		"id": "ds-1", "title": "Test",
		"resources": []map[string]any{
			{"id": "res-1", "name": "data.csv", "format": "CSV", "url": fileSrv.URL + "/data.csv"},
		},
	})

	outPath := filepath.Join(t.TempDir(), "output.csv")
	path, err := RunDownload(context.Background(), c, "ds-1", "res-1", outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != outPath {
		t.Errorf("expected path %s, got %s", outPath, path)
	}
	content, _ := os.ReadFile(outPath)
	if string(content) != "col1,col2\n1,2\n" {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestRunDownload_ResourceNotFound(t *testing.T) {
	c := newMockClient(t, map[string]any{
		"id": "ds-3", "title": "Test",
		"resources": []map[string]any{},
	})

	_, err := RunDownload(context.Background(), c, "ds-3", "nonexistent-res", "")
	if err == nil {
		t.Fatal("expected error for missing resource")
	}
}

func TestRunDownload_DefaultOutputPath(t *testing.T) {
	fileSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("data"))
	}))
	defer fileSrv.Close()

	c := newMockClient(t, map[string]any{
		"id": "ds-2", "title": "Test",
		"resources": []map[string]any{
			{"id": "res-2", "name": "report.json", "format": "JSON", "url": fileSrv.URL + "/report.json"},
		},
	})

	// Output path vuoto → usa il nome della risorsa nella dir corrente
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	path, err := RunDownload(context.Background(), c, "ds-2", "res-2", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(path) != "report.json" {
		t.Errorf("expected filename report.json, got %s", filepath.Base(path))
	}
	os.Remove(path)
}
