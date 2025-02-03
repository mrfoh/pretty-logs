package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type JsonFormatter struct {
	Options  *FormatterOptions
	Template *template.Template
}

// Create a new JsonFormatter with the provided options.
func NewJsonFormatter(options *FormatterOptions) (Formatter, error) {
	// Set default error keys if not provided.
	if options.ErrorObjectKeys == nil {
		options.ErrorObjectKeys = []string{"err", "error"}
	}
	// Set default timestamp format if not provided.
	if options.TimestampFormat == "" {
		options.TimestampFormat = "2006-01-02 15:04:05.000"
	}

	jf := &JsonFormatter{
		Options: options,
	}

	// Parse the template file once if provided.
	if options.FormatTemplateFile != "" {
		// Check if the template file exists
		if _, err := os.Stat(options.FormatTemplateFile); err != nil {
			return nil, fmt.Errorf("template file does not exist: %s", options.FormatTemplateFile)
		}

		tmpl, err := template.New(options.FormatTemplateFile).ParseFiles(options.FormatTemplateFile)
		if err != nil {
			return nil, fmt.Errorf("error parsing template file: %w", err)
		}

		jf.Template = tmpl
	}

	return jf, nil
}

// Print the log line using the provided template.
func (f *JsonFormatter) PrintLogLineTemplate(line map[string]interface{}) error {
	input := LogLineMapToStruct(line, f.Options)

	if f.Template == nil {
		return fmt.Errorf("no template available")
	}

	if err := f.Template.Execute(os.Stdout, input); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
}

// Print the log line to stdout using the provided template or raw JSON.
func (f *JsonFormatter) PrintLogLine(line map[string]interface{}) {
	var err error
	if f.Template != nil {
		err = f.PrintLogLineTemplate(line)
	} else {
		// Fallback: print the raw JSON log line.
		var out []byte
		out, err = json.Marshal(line)
		if err == nil {
			fmt.Println(string(out))
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error printing log line: %v\n", err)
	}
}

// Process the input file and print the formatted log lines.
func (f *JsonFormatter) Process(input *os.File) error {
	scanner := bufio.NewScanner(input)
	// Increase buffer size for large lines.
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		var entry map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing log entry: %v\n", err)
			continue
		}

		logLine := make(map[string]interface{})

		// Extract standard fields.
		logLine[f.Options.TimeKey] = ExtractValue(entry, f.Options.TimeKey)
		logLine[f.Options.LevelKey] = ExtractValue(entry, f.Options.LevelKey)
		logLine[f.Options.PidKey] = ExtractValue(entry, f.Options.PidKey)
		logLine[f.Options.NameKey] = ExtractValue(entry, f.Options.NameKey)
		logLine[f.Options.ContextKey] = ExtractValue(entry, f.Options.ContextKey)
		logLine[f.Options.HostnameKey] = ExtractValue(entry, f.Options.HostnameKey)
		logLine[f.Options.MsgKey] = ExtractValue(entry, f.Options.MsgKey)
		logLine[f.Options.RequestKey] = ExtractValue(entry, f.Options.RequestKey)

		// Handle error if present.
		if HasAnyKey(entry, f.Options.ErrorObjectKeys) {
			// Find the first error key and extract the error message.
			for _, key := range f.Options.ErrorObjectKeys {
				// Check if the key exists in the log entry.
				if HasKey(entry, key) {
					logLine[key] = ExtractError(entry[key])
					break
				}
			}
		}

		f.PrintLogLine(logLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}
