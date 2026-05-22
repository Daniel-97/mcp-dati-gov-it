# dati-gov-it-cli + dati-gov-it-mcp

CLI e server MCP per il catalogo open data del governo italiano ([dati.gov.it](https://dati.gov.it)).

Compatibile con [Printing Press](https://printingpress.dev).

## Installazione

```bash
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-cli@latest
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-mcp@latest
```

## Uso rapido

```bash
dati-gov-it-cli search "qualità aria"
dati-gov-it-cli show <id>
dati-gov-it-cli download <dataset-id> <resource-id> --output dati.csv
dati-gov-it-cli orgs --details
```

## Configurazione MCP (Claude Desktop)

```json
{
  "mcpServers": {
    "dati-gov-it": {
      "command": "dati-gov-it-mcp"
    }
  }
}
```

## API

Basato su [CKAN API](https://docs.ckan.org/en/2.9/api/) — endpoint pubblici, nessuna autenticazione.

## Licenza

MIT
