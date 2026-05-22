package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/spf13/cobra"
)

// SearchOptions sono i parametri per un'operazione di ricerca.
type SearchOptions struct {
	Query string
	Tags  string
	Org   string
	Rows  int
	Agent bool
}

// SearchOutput è il risultato strutturato (usato in --agent e MCP).
type SearchOutput struct {
	Count   int              `json:"count"`
	Results []client.Dataset `json:"results"`
}

// RunSearch esegue la ricerca dataset tramite l'API CKAN.
func RunSearch(ctx context.Context, c *client.Client, opts SearchOptions) (*SearchOutput, error) {
	result, err := c.SearchDatasets(ctx, opts.Query, opts.Tags, opts.Org, opts.Rows)
	if err != nil {
		return nil, err
	}
	return &SearchOutput{Count: result.Count, Results: result.Results}, nil
}

// NewSearchCmd crea il comando Cobra `search`.
func NewSearchCmd(c *client.Client) *cobra.Command {
	var opts SearchOptions

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Cerca dataset su dati.gov.it",
		Long:  "Cerca dataset nel catalogo open data del governo italiano tramite l'API CKAN.",
		Example: `  dati-gov-it-cli search "istruzione"
  dati-gov-it-cli search "ambiente" --tags aria --rows 20
  dati-gov-it-cli search "comune" --org regione-lombardia --agent`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Query = args[0]
			out, err := RunSearch(cmd.Context(), c, opts)
			if err != nil {
				return err
			}
			return printSearchOutput(cmd.OutOrStdout(), out, opts.Agent)
		},
	}

	cmd.Flags().StringVar(&opts.Tags, "tags", "", "Filtra per tag (es. economia)")
	cmd.Flags().StringVar(&opts.Org, "org", "", "Filtra per organizzazione slug (es. regione-lombardia)")
	cmd.Flags().IntVar(&opts.Rows, "rows", 10, "Numero massimo di risultati")
	cmd.Flags().BoolVar(&opts.Agent, "agent", false, "Output JSON strutturato per agenti IA")
	return cmd
}

func printSearchOutput(w io.Writer, out *SearchOutput, agent bool) error {
	if agent {
		return json.NewEncoder(w).Encode(out)
	}
	fmt.Fprintf(w, "Trovati %d dataset (mostrati %d):\n\n", out.Count, len(out.Results))
	for _, d := range out.Results {
		tags := make([]string, 0, len(d.Tags))
		for _, t := range d.Tags {
			tags = append(tags, t.Name)
		}
		fmt.Fprintf(w, "  [%s]\n  Titolo: %s\n  Tag: %s\n\n",
			d.ID, d.Title, strings.Join(tags, ", "))
	}
	return nil
}
