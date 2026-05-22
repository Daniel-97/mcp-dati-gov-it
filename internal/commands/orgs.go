package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/spf13/cobra"
)

// RunOrgs restituisce la lista delle organizzazioni.
// Se details=true, chiama organization_show per ogni org per ottenere titolo e descrizione.
func RunOrgs(ctx context.Context, c *client.Client, details bool) ([]client.Organization, error) {
	names, err := c.ListOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	orgs := make([]client.Organization, 0, len(names))
	for _, name := range names {
		if details {
			o, err := c.GetOrganization(ctx, name)
			if err != nil {
				// Fallback al solo nome se organization_show fallisce
				orgs = append(orgs, client.Organization{Name: name})
				continue
			}
			orgs = append(orgs, *o)
		} else {
			orgs = append(orgs, client.Organization{Name: name})
		}
	}
	return orgs, nil
}

// NewOrgsCmd crea il comando Cobra `orgs`.
func NewOrgsCmd(c *client.Client) *cobra.Command {
	var details, agent bool

	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Elenca le PA che pubblicano dati su dati.gov.it",
		Example: `  dati-gov-it-cli orgs
  dati-gov-it-cli orgs --details
  dati-gov-it-cli orgs --agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			orgs, err := RunOrgs(cmd.Context(), c, details)
			if err != nil {
				return err
			}
			return printOrgsOutput(cmd.OutOrStdout(), orgs, agent)
		},
	}
	cmd.Flags().BoolVar(&details, "details", false, "Recupera titolo e descrizione di ogni organizzazione")
	cmd.Flags().BoolVar(&agent, "agent", false, "Output JSON strutturato per agenti IA")
	return cmd
}

func printOrgsOutput(w io.Writer, orgs []client.Organization, agent bool) error {
	if agent {
		return json.NewEncoder(w).Encode(orgs)
	}
	fmt.Fprintf(w, "Organizzazioni (%d):\n\n", len(orgs))
	for _, o := range orgs {
		if o.Title != "" {
			fmt.Fprintf(w, "  %s — %s\n", o.Name, o.Title)
		} else {
			fmt.Fprintf(w, "  %s\n", o.Name)
		}
	}
	return nil
}
