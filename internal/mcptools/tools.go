package mcptools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/daniel-97/mcp-dati-gov-it/internal/commands"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Handlers raccoglie i handler MCP condividendo il client.
type Handlers struct {
	client *client.Client
}

// NewHandlers crea Handlers con il client fornito.
func NewHandlers(c *client.Client) *Handlers {
	return &Handlers{client: c}
}

// RegisterTools registra tutti i tool MCP sul server.
func RegisterTools(s *server.MCPServer, h *Handlers) {
	s.AddTool(
		mcp.NewTool("search_datasets",
			mcp.WithDescription("Cerca dataset nel catalogo open data dati.gov.it"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Testo da cercare")),
			mcp.WithString("tags", mcp.Description("Filtra per tag (es. economia)")),
			mcp.WithString("org", mcp.Description("Filtra per organizzazione slug")),
			mcp.WithNumber("rows", mcp.Description("Numero massimo di risultati (default 10)")),
		),
		h.SearchDatasets,
	)
	s.AddTool(
		mcp.NewTool("get_dataset",
			mcp.WithDescription("Ottieni i metadati completi di un dataset per ID"),
			mcp.WithString("id", mcp.Required(), mcp.Description("ID del dataset CKAN")),
		),
		h.GetDataset,
	)
	s.AddTool(
		mcp.NewTool("list_organizations",
			mcp.WithDescription("Elenca le PA che pubblicano dati su dati.gov.it"),
			mcp.WithBoolean("details", mcp.Description("Se true, recupera titolo e descrizione di ogni PA")),
		),
		h.ListOrganizations,
	)
	s.AddTool(
		mcp.NewTool("download_resource",
			mcp.WithDescription("Scarica un file risorsa associato a un dataset"),
			mcp.WithString("dataset_id", mcp.Required(), mcp.Description("ID del dataset")),
			mcp.WithString("resource_id", mcp.Required(), mcp.Description("ID della risorsa")),
			mcp.WithString("output_path", mcp.Description("Path di destinazione (opzionale)")),
		),
		h.DownloadResource,
	)
}

// argsMap estrae Arguments come map[string]any da una CallToolRequest.
func argsMap(req mcp.CallToolRequest) map[string]any {
	if m, ok := req.Params.Arguments.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

// SearchDatasets è il handler MCP per search_datasets.
func (h *Handlers) SearchDatasets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := argsMap(req)
	query, _ := args["query"].(string)
	tags, _ := args["tags"].(string)
	org, _ := args["org"].(string)
	rows := 10
	if r, ok := args["rows"].(float64); ok {
		rows = int(r)
	}

	out, err := commands.RunSearch(ctx, h.client, commands.SearchOptions{
		Query: query, Tags: tags, Org: org, Rows: rows,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

// GetDataset è il handler MCP per get_dataset.
func (h *Handlers) GetDataset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := argsMap(req)
	id, _ := args["id"].(string)
	if id == "" {
		return mcp.NewToolResultError("id è obbligatorio"), nil
	}
	d, err := commands.RunShow(ctx, h.client, id)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(d)
}

// ListOrganizations è il handler MCP per list_organizations.
func (h *Handlers) ListOrganizations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := argsMap(req)
	details, _ := args["details"].(bool)
	orgs, err := commands.RunOrgs(ctx, h.client, details)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(orgs)
}

// DownloadResource è il handler MCP per download_resource.
func (h *Handlers) DownloadResource(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := argsMap(req)
	datasetID, _ := args["dataset_id"].(string)
	resourceID, _ := args["resource_id"].(string)
	outputPath, _ := args["output_path"].(string)

	if datasetID == "" || resourceID == "" {
		return mcp.NewToolResultError("dataset_id e resource_id sono obbligatori"), nil
	}
	path, err := commands.RunDownload(ctx, h.client, datasetID, resourceID, outputPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("File salvato: %s", path)), nil
}

// jsonResult serializza v come JSON e lo incapsula in un CallToolResult.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
