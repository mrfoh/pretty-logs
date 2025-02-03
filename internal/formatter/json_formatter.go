package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type JsonFormatter struct {
	Options *FormatterOptions
}

func NewJsonFormatter(options *FormatterOptions) Formatter {
	// Set default error keys if not provided
	if options.ErrorObjectKeys == nil {
		options.ErrorObjectKeys = []string{"err", "error"}
	}
	// Set default timestamp format if not provided
	if options.TimestampFormat == "" {
		options.TimestampFormat = "2006-01-02 15:04:05.000"
	}
	return &JsonFormatter{
		Options: options,
	}
}

func (f *JsonFormatter) PrintLogLineTemplate(line map[string]interface{}) error {
	tmpl, err := template.New(f.Options.FormatTemplateFile).ParseFiles(f.Options.FormatTemplateFile)
	if err != nil {
		return fmt.Errorf("error parsing template file: %v", err)
	}

	input := LogLineMapToStruct(line, f.Options)

	err = tmpl.Execute(os.Stdout, input)
	if err != nil {
		panic(err)
	}

	return nil
}

func (f *JsonFormatter) PrintLogLine(line map[string]interface{}) {
	// Check if the template file exists
	if _, err := os.Stat(f.Options.FormatTemplateFile); err == nil {
		// Use the template file
		if err := f.PrintLogLineTemplate(line); err != nil {
			fmt.Fprintf(os.Stderr, "Error printing log line: %v\n", err)
		}
	}
}

func (f *JsonFormatter) Process(input *os.File) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // Increase buffer size for large lines

	for scanner.Scan() {
		var entry map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing log entry: %v\n", err)
			continue
		}

		logLine := make(map[string]interface{})

		// Extract standard fields
		logLine[f.Options.TimeKey] = ExtractValue(entry, f.Options.TimeKey)
		logLine[f.Options.LevelKey] = ExtractValue(entry, f.Options.LevelKey)
		logLine[f.Options.PidKey] = ExtractValue(entry, f.Options.PidKey)
		logLine[f.Options.NameKey] = ExtractValue(entry, f.Options.NameKey)
		logLine[f.Options.ContextKey] = ExtractValue(entry, f.Options.ContextKey)
		logLine[f.Options.HostnameKey] = ExtractValue(entry, f.Options.HostnameKey)
		logLine[f.Options.MsgKey] = ExtractValue(entry, f.Options.MsgKey)
		logLine[f.Options.RequestKey] = ExtractValue(entry, f.Options.RequestKey)

		// Handle error if present
		if HasKeys(&entry, f.Options.ErrorObjectKeys) {
			for _, key := range f.Options.ErrorObjectKeys {
				if err, ok := entry[key]; ok {
					logLine["error"] = FormatError(err)
					break
				}
			}
		}

		f.PrintLogLine(logLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}
