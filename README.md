# logslice

Fast log filtering utility with time-range and field-based queries for structured JSON logs.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

## Usage

Filter logs by time range:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log
```

Filter by field value:

```bash
logslice --field level=error --field service=api app.log
```

Combine time range and field filters:

```bash
logslice --from "2024-01-15T08:00:00Z" --field level=error --field user_id=42 app.log
```

Pipe from stdin:

```bash
cat app.log | logslice --from "2024-01-15T08:00:00Z" --field status=500
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--field` | Field filter as `key=value` (repeatable) |
| `--timestamp` | JSON field to use as timestamp (default: `time`) |
| `--pretty` | Pretty-print matching log lines |

## Requirements

- Go 1.21+

## License

MIT © 2024 yourusername