# dati-gov-it-cli + dati-gov-it-mcp

CLI and MCP server for the Italian government open data catalog ([dati.gov.it](https://dati.gov.it)).

Compatible with [Printing Press](https://printingpress.dev).

## Installation

```bash
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-cli@latest
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-mcp@latest
```

## Quick Start

```bash
dati-gov-it-cli search "air quality"
dati-gov-it-cli show <id>
dati-gov-it-cli download <dataset-id> <resource-id> --output data.csv
dati-gov-it-cli orgs --details
```

## MCP Configuration (Claude Desktop)

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

Based on the [CKAN API](https://docs.ckan.org/en/2.9/api/) — public endpoints, no authentication required.

## License

MIT
