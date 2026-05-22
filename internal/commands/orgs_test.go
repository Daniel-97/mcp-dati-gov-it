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

func TestRunOrgs_Basic(t *testing.T) {
	c := newMockClient(t, []string{"comune-di-roma", "regione-lombardia"})

	out, err := RunOrgs(context.Background(), c, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 orgs, got %d", len(out))
	}
	if out[0].Name != "comune-di-roma" {
		t.Errorf("unexpected first org: %s", out[0].Name)
	}
}

func TestRunOrgs_WithDetails(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			json.NewEncoder(w).Encode(map[string]any{"success": true, "result": []string{"regione-lombardia"}})
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"result": map[string]any{
					"id": "rl-id", "name": "regione-lombardia",
					"title": "Regione Lombardia", "description": "Dati aperti della Regione Lombardia",
				},
			})
		}
	}))
	defer srv.Close()
	c := client.NewWithBase(srv.URL)

	out, err := RunOrgs(context.Background(), c, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 org, got %d", len(out))
	}
	if out[0].Title != "Regione Lombardia" {
		t.Errorf("unexpected title: %s", out[0].Title)
	}
}

func TestNewOrgsCmd_AgentOutput(t *testing.T) {
	c := newMockClient(t, []string{"comune-di-roma"})

	cmd := NewOrgsCmd(c)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--agent"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd error: %v", err)
	}
	var orgs []client.Organization
	if err := json.NewDecoder(&buf).Decode(&orgs); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(orgs) != 1 {
		t.Errorf("expected 1 org, got %d", len(orgs))
	}
}
