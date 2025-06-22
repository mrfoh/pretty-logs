package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

		tmpl, err := template.New(filepath.Base(options.FormatTemplateFile)).ParseFiles(options.FormatTemplateFile)
		if err != nil {
			return nil, fmt.Errorf("error parsing template file: %w", err)
		}

		jf.Template = tmpl
	}

	return jf, nil
}

// Print the log line using the provided template.
func (f *JsonFormatter) PrintLogLineTemplate(logLine LogLine) error {
	if f.Template == nil {
		return fmt.Errorf("no template available")
	}

	if err := f.Template.Execute(os.Stdout, logLine); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
}

// Print the log line to stdout using the provided template or raw JSON.
func (f *JsonFormatter) PrintLogLine(logLine LogLine) {
	var err error
	if f.Template != nil {
		err = f.PrintLogLineTemplate(logLine)
	} else {
		// Fallback: print the raw JSON log line.
		var out []byte
		out, err = json.Marshal(logLine)
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

		logLine := LogLineMapToStruct(entry, f.Options)
		f.PrintLogLine(logLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}
