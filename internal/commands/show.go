package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/spf13/cobra"
)

// RunShow restituisce i metadati completi di un dataset.
func RunShow(ctx context.Context, c *client.Client, id string) (*client.Dataset, error) {
	return c.GetDataset(ctx, id)
}

// NewShowCmd crea il comando Cobra `show`.
func NewShowCmd(c *client.Client) *cobra.Command {
	var agent bool

	cmd := &cobra.Command{
		Use:   "show <dataset-id>",
		Short: "Mostra i metadati completi di un dataset",
		Example: `  dati-gov-it-cli show d3e3a7c2-1234-5678-abcd-ef0123456789
  dati-gov-it-cli show <id> --agent`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := RunShow(cmd.Context(), c, args[0])
			if err != nil {
				return err
			}
			return printShowOutput(cmd.OutOrStdout(), d, agent)
		},
	}
	cmd.Flags().BoolVar(&agent, "agent", false, "Output JSON strutturato per agenti IA")
	return cmd
}

func printShowOutput(w io.Writer, d *client.Dataset, agent bool) error {
	if agent {
		return json.NewEncoder(w).Encode(d)
	}
	fmt.Fprintf(w, "ID:       %s\n", d.ID)
	fmt.Fprintf(w, "Titolo:   %s\n", d.Title)
	fmt.Fprintf(w, "Licenza:  %s\n", d.LicenseTitle)
	fmt.Fprintf(w, "Modifica: %s\n", d.MetadataModified)
	if d.Organization != nil {
		fmt.Fprintf(w, "Ente:     %s\n", d.Organization.Title)
	}
	fmt.Fprintf(w, "\nDescrizione:\n%s\n", d.Notes)
	if len(d.Resources) > 0 {
		fmt.Fprintf(w, "\nRisorse (%d):\n", len(d.Resources))
		for _, r := range d.Resources {
			fmt.Fprintf(w, "  [%s] %s (%s) — %s\n", r.ID, r.Name, r.Format, r.URL)
		}
	}
	return nil
}
