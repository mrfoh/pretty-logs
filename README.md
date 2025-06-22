# Pretty Logs

Pretty Logs is a powerful command-line tool designed to format and beautify JSON log output. It takes JSON-formatted log files as input and transforms them into human-readable output with customizable formatting options. This tool is particularly useful for developers and system administrators who need to analyze log files from applications that output JSON logs.

## Features

Pretty Logs offers extensive customization options to handle various JSON log formats:

- Flexible timestamp formatting with customizable patterns
- Support for custom field mapping through command-line flags
- **Additional keys extraction** for arbitrary JSON fields in templates
- Error object detection and formatting
- Request/response object handling
- Template-based output formatting
- Support for large log files through efficient buffered reading
- Cross-platform compatibility (Linux, macOS, Windows)

## Installation

To install Pretty Logs, you'll need Go 1.23.3 or later. You can build the binary using one of the following methods:

### Building from Source

Clone the repository and build using Make:

```bash
git clone https://github.com/mrfoh/pretty-logs.git
cd pretty-logs
make build
```

The binary will be available in the `bin` directory.

### Cross-Platform Builds

The project includes Make targets for various platforms:

- Linux (AMD64): `make build-linux`
- Linux (ARM64): `make build-linux-arm`
- macOS (AMD64): `make build-macos`
- macOS (ARM64): `make build-macos-arm`
- macOS (Universal): `make build-macos-universal`

## Usage

Pretty Logs reads JSON log input from stdin and outputs formatted logs to stdout. Here's the basic usage:

```bash
cat logfile.json | prettylogs [flags]
```

### Command Line Flags

The tool provides numerous flags to customize its behavior:

- `-F, --formatTemplateFile`: Specify a custom template file for log formatting
- `-f, --timestampFormat`: Set the timestamp format (default: "2006-01-02 15:04:05.000")
- `-L, --levelKey`: Specify the key for log level (default: "level")
- `-t, --timeKey`: Specify the key for timestamp (default: "time")
- `-p, --pidKey`: Specify the key for process ID (default: "pid")
- `-n, --nameKey`: Specify the key for application name (default: "name")
- `-c, --contextKey`: Specify the key for context (default: "context")
- `-m, --msgKey`: Specify the key for message (default: "msg")
- `-k, --errorLikeObjectKeys`: Specify keys for error objects (default: ["err", "error"])
- `-H, --hostnameKey`: Specify the key for hostname (default: "hostname")
- `-r, --requestKey`: Specify the key for request objects (default: "req")
- `-R, --responseKey`: Specify the key for response objects (default: "res")
- `-a, --additionalKeys`: Specify additional JSON keys to extract for template use (comma-separated)

### Example Usage

Basic usage with default settings:
```bash
tail -f application.log | prettylogs
```

Using custom field mappings:
```bash
cat logs.json | prettylogs --levelKey severity --timeKey timestamp --msgKey message
```

Using a custom format template:
```bash
cat logs.json | prettylogs -F my_template.tmpl
```

Extracting additional JSON fields for templates:
```bash
# Extract traceId and userId fields
cat logs.json | prettylogs --additionalKeys traceId,userId -F my_template.tmpl
```

Real-world example with additional keys:
```bash
echo '{"time":"2023-01-01T10:00:00Z","level":"error","msg":"Database connection failed","traceId":"abc123","userId":"user456","service":"api"}' | prettylogs -a traceId,userId,service
```

## Template Format

Pretty Logs supports custom templates for output formatting. Templates use Go's text/template syntax and have access to the following fields:

- `{{.Time}}`: Formatted timestamp
- `{{.Level}}`: Log level
- `{{.Pid}}`: Process ID
- `{{.Name}}`: Application name
- `{{.Context}}`: Log context
- `{{.Msg}}`: Log message
- `{{.Error}}`: Error message (if present)
- `{{.Req}}`: Request object (if present)
- `{{.Hostname}}`: Hostname
- `{{.Additional.keyName}}`: Additional extracted fields (when using `--additionalKeys`)

### Additional Keys in Templates

When using the `--additionalKeys` flag, you can access the extracted fields in templates using the `{{.Additional.keyName}}` syntax:

```go
{{.Time}} [{{.Level}}] {{.Msg}}
{{- if .Additional.traceId}} | Trace: {{.Additional.traceId}}{{end}}
{{- if .Additional.userId}} | User: {{.Additional.userId}}{{end}}
{{- if .Additional.service}} | Service: {{.Additional.service}}{{end}}
```

Example template output:
```
2023-01-01T10:00:00Z [error] Database connection failed | Trace: abc123 | User: user456 | Service: api
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests. Before submitting, ensure:

1. Your code passes all tests (`make test`)
2. Your code passes the linter (`make lint`)
3. You've added tests for new features
4. You've updated documentation as needed

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.