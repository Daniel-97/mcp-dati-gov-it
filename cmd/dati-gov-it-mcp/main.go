package main

import (
	"fmt"
	"os"

	"github.com/daniel-97/mcp-dati-gov-it/internal/client"
	"github.com/daniel-97/mcp-dati-gov-it/internal/mcptools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	c := client.New()

	srv := server.NewMCPServer(
		"dati-gov-it",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	handlers := mcptools.NewHandlers(c)
	mcptools.RegisterTools(srv, handlers)

	if err := server.ServeStdio(srv); err != nil {
		fmt.Fprintf(os.Stderr, "Errore server MCP: %v\n", err)
		os.Exit(1)
	}
}
