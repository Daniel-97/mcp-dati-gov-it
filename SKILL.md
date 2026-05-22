---
name: dati-gov-it
description: Cerca e scarica dataset open data del governo italiano da dati.gov.it
author: daniel-97
license: MIT
requires: install
bins: dati-gov-it-cli
module: go
---

# dati-gov-it CLI

CLI agent-native per il catalogo open data del governo italiano (dati.gov.it, CKAN API).

## Prerequisites

```bash
go install github.com/daniel-97/mcp-dati-gov-it/cmd/dati-gov-it-cli@latest
```

## When Not to Use

- Per scaricare file molto grandi (>1GB) considera un client HTTP diretto
- Per operazioni di scrittura sul catalogo (richiedono API key amministratore)

## Command Reference

| Comando | Descrizione |
|---|---|
| `search <query>` | Cerca dataset per keyword |
| `show <id>` | Metadati completi di un dataset |
| `orgs` | Lista PA che pubblicano dati |
| `download <dataset-id> <resource-id>` | Scarica file risorsa |

## Agent Mode

Ogni comando supporta `--agent` per output JSON strutturato:

```bash
dati-gov-it-cli search "istruzione" --agent | jq .
dati-gov-it-cli show <id> --agent
dati-gov-it-cli orgs --details --agent
```

## Exit Codes

| Codice | Significato |
|---|---|
| 0 | Successo |
| 1 | Errore generico (API, rete, file system) |

## Direct Use

```bash
dati-gov-it-cli search "qualità aria" --tags ambiente --rows 5
dati-gov-it-cli show <id-dataset>
dati-gov-it-cli download <dataset-id> <resource-id> --output dati.csv
```
