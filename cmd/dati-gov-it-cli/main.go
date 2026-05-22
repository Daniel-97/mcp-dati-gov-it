package main

import (
	"os"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/daniel-97/mcp-dati-gov-it/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	c := client.New()

	root := &cobra.Command{
		Use:   "dati-gov-it-cli",
		Short: "CLI per il catalogo open data del governo italiano",
		Long: `dati-gov-it-cli permette di cercare, visualizzare e scaricare dataset
dal portale dati.gov.it (CKAN API).

Usa --agent su qualsiasi comando per output JSON strutturato.`,
		Version: "0.1.0",
	}

	root.AddCommand(
		commands.NewSearchCmd(c),
		commands.NewShowCmd(c),
		commands.NewOrgsCmd(c),
		commands.NewDownloadCmd(c),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
