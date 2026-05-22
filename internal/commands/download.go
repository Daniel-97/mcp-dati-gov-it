package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/spf13/cobra"
)

// RunDownload scarica la risorsa indicata e la salva su outputPath.
// Se outputPath è vuoto, usa il nome della risorsa nella directory corrente.
// Restituisce il path del file salvato.
func RunDownload(ctx context.Context, c *client.Client, datasetID, resourceID, outputPath string) (string, error) {
	// Recupera il dataset per trovare l'URL della risorsa
	d, err := c.GetDataset(ctx, datasetID)
	if err != nil {
		return "", fmt.Errorf("recupero dataset: %w", err)
	}

	var resURL, resName, resFormat string
	for _, r := range d.Resources {
		if r.ID == resourceID {
			resURL = r.URL
			resName = r.Name
			resFormat = r.Format
			break
		}
	}

	if resURL == "" {
		return "", fmt.Errorf("risorsa %s non trovata nel dataset %s", resourceID, datasetID)
	}

	// Determina il path di output
	if outputPath == "" {
		if resName != "" {
			outputPath = resName
		} else {
			outputPath = resourceID + "." + strings.ToLower(resFormat)
		}
	}

	// Scarica il file
	if err := downloadFile(ctx, resURL, outputPath); err != nil {
		os.Remove(outputPath) // rimuovi file parziale
		return "", fmt.Errorf("download: %w", err)
	}
	return outputPath, nil
}

func downloadFile(ctx context.Context, url, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// NewDownloadCmd crea il comando Cobra `download`.
func NewDownloadCmd(c *client.Client) *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download <dataset-id> <resource-id>",
		Short: "Scarica un file risorsa di un dataset",
		Example: `  dati-gov-it-cli download <dataset-id> <resource-id>
  dati-gov-it-cli download <dataset-id> <resource-id> --output ./dati.csv`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := RunDownload(cmd.Context(), c, args[0], args[1], outputPath)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "File salvato: %s\n", path)
			return nil
		},
	}
	cmd.Flags().StringVar(&outputPath, "output", "", "Path di destinazione (default: nome risorsa nella dir corrente)")
	return cmd
}
