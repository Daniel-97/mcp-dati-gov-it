package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
)

func TestRunShow_ReturnsDataset(t *testing.T) {
	c := newMockClient(t, map[string]any{
		"id": "abc", "title": "Dataset Ambiente",
		"notes": "Dati ambientali nazionali",
		"resources": []map[string]any{
			{"id": "r1", "name": "dati.csv", "format": "CSV", "url": "https://example.com/dati.csv"},
		},
	})

	d, err := RunShow(context.Background(), c, "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Title != "Dataset Ambiente" {
		t.Errorf("unexpected title: %s", d.Title)
	}
	if len(d.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(d.Resources))
	}
}

func TestRunShow_PropagatesError(t *testing.T) {
	c := newErrorMockClient(t, "dataset non trovato")
	_, err := RunShow(context.Background(), c, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewShowCmd_AgentOutput(t *testing.T) {
	c := newMockClient(t, map[string]any{"id": "xyz", "title": "Test Show"})

	cmd := NewShowCmd(c)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"xyz", "--agent"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd error: %v", err)
	}
	var d client.Dataset
	if err := json.NewDecoder(&buf).Decode(&d); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if d.ID != "xyz" {
		t.Errorf("expected id xyz, got %s", d.ID)
	}
}
