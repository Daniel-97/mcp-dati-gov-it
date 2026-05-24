# dati-gov-it-cli + dati-gov-it-mcp

CLI and MCP server for the Italian government open data catalog ([dati.gov.it](https://dati.gov.it)).

Built with [Printing Press](https://printingpress.dev).

## Installation

```bash
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-cli@latest
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-mcp@latest
```

## Build from Source

```bash
git clone https://github.com/daniel-97/mcp-dati-gov-it.git
cd mcp-dati-gov-it

go build -o dati-gov-it-cli ./cmd/dati-gov-it-cli/
go build -o dati-gov-it-mcp ./cmd/dati-gov-it-mcp/
```

The binaries will be available in the current directory as `./dati-gov-it-cli` and `./dati-gov-it-mcp`.

## Quick Start

```bash
dati-gov-it-cli search "air quality"
dati-gov-it-cli show <id>
dati-gov-it-cli download <dataset-id> <resource-id> --output data.csv
dati-gov-it-cli orgs --details
```

## Examples

### Search datasets

Search across the full catalog by keyword:

```bash
./dati-gov-it-cli search "qualità aria"
```

Filter by tag or organization, and increase the number of results:

```bash
./dati-gov-it-cli search "ambiente" --tags aria --rows 20
./dati-gov-it-cli search "bilancio" --org comune-di-roma
```

### Show dataset details

Retrieve full metadata for a dataset (title, description, license, resources):

```bash
./dati-gov-it-cli show <dataset-id>
```

### Download a resource

Download a specific file attached to a dataset. Use `show` first to get the resource ID:

```bash
./dati-gov-it-cli show <dataset-id>
./dati-gov-it-cli download <dataset-id> <resource-id> --output data.csv
```

If `--output` is omitted, the file is saved in the current directory using the resource name.

### List organizations

List all public administrations that publish data on dati.gov.it:

```bash
./dati-gov-it-cli orgs
```

Include title and description for each organization:

```bash
./dati-gov-it-cli orgs --details
```

### Agent / JSON output

Every command supports `--agent` to output structured JSON, useful for scripting or AI agents:

```bash
./dati-gov-it-cli search "istruzione" --agent
./dati-gov-it-cli show <dataset-id> --agent
./dati-gov-it-cli orgs --agent
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
