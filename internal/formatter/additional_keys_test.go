package formatter

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdditionalKeysExtraction(t *testing.T) {
	tests := []struct {
		name           string
		jsonInput      string
		additionalKeys []string
		expectedKeys   map[string]interface{}
	}{
		{
			name:           "Single additional key",
			jsonInput:      `{"msg":"test message","traceId":"abc123","level":"info"}`,
			additionalKeys: []string{"traceId"},
			expectedKeys:   map[string]interface{}{"traceId": "abc123"},
		},
		{
			name:           "Multiple additional keys",
			jsonInput:      `{"msg":"test","userId":"user123","sessionId":"sess456","traceId":"trace789","level":"info"}`,
			additionalKeys: []string{"userId", "sessionId", "traceId"},
			expectedKeys: map[string]interface{}{
				"userId":    "user123",
				"sessionId": "sess456",
				"traceId":   "trace789",
			},
		},
		{
			name:           "Additional key with different data types",
			jsonInput:      `{"msg":"test","count":42,"active":true,"score":3.14,"level":"info"}`,
			additionalKeys: []string{"count", "active", "score"},
			expectedKeys: map[string]interface{}{
				"count":  float64(42),
				"active": true,
				"score":  3.14,
			},
		},
		{
			name:           "Non-existent additional key",
			jsonInput:      `{"msg":"test message","level":"info"}`,
			additionalKeys: []string{"nonExistent"},
			expectedKeys:   map[string]interface{}{},
		},
		{
			name:           "Empty additional keys",
			jsonInput:      `{"msg":"test message","customField":"value","level":"info"}`,
			additionalKeys: []string{},
			expectedKeys:   map[string]interface{}{},
		},
		{
			name:           "Additional key with nested object",
			jsonInput:      `{"msg":"test","metadata":{"service":"api","version":"1.0"},"level":"info"}`,
			additionalKeys: []string{"metadata"},
			expectedKeys: map[string]interface{}{
				"metadata": map[string]interface{}{
					"service": "api",
					"version": "1.0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &FormatterOptions{
				TimeKey:         "time",
				LevelKey:        "level",
				MsgKey:          "msg",
				AdditionalKeys:  tt.additionalKeys,
				TimestampFormat: DEFAULT_TIMESTAMP_FORMAT,
			}

			var entry map[string]interface{}
			err := json.Unmarshal([]byte(tt.jsonInput), &entry)
			require.NoError(t, err)

			logLine := LogLineMapToStruct(entry, options)

			// Verify that additional keys are extracted correctly
			assert.Equal(t, len(tt.expectedKeys), len(logLine.Additional))
			for key, expectedValue := range tt.expectedKeys {
				assert.Equal(t, expectedValue, logLine.Additional[key])
			}
		})
	}
}

func TestAdditionalKeysInTemplate(t *testing.T) {
	templateContent := `Time: {{.Time}}
Level: {{.Level}}
Message: {{.Msg}}
{{- if .Additional.traceId}}
TraceId: {{.Additional.traceId}}
{{- end}}
{{- if .Additional.userId}}
UserId: {{.Additional.userId}}
{{- end}}
{{- if .Additional.metadata}}
Metadata: {{.Additional.metadata}}
{{- end}}`

	// Create a temporary template file
	tmpFile, err := os.CreateTemp("", "test_template_*.tpl")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, err = tmpFile.WriteString(templateContent)
	require.NoError(t, err)
	_ = tmpFile.Close()

	tests := []struct {
		name           string
		jsonInput      string
		additionalKeys []string
		expectedOutput string
	}{
		{
			name:           "Template with additional keys present",
			jsonInput:      `{"time":"2023-01-01T10:00:00Z","level":"info","msg":"test message","traceId":"abc123","userId":"user456"}`,
			additionalKeys: []string{"traceId", "userId"},
			expectedOutput: "Time: 2023-01-01T10:00:00Z\nLevel: info\nMessage: test message\nTraceId: abc123\nUserId: user456",
		},
		{
			name:           "Template with some additional keys missing",
			jsonInput:      `{"time":"2023-01-01T10:00:00Z","level":"info","msg":"test message","traceId":"abc123"}`,
			additionalKeys: []string{"traceId", "userId"},
			expectedOutput: "Time: 2023-01-01T10:00:00Z\nLevel: info\nMessage: test message\nTraceId: abc123",
		},
		{
			name:           "Template with nested object in additional key",
			jsonInput:      `{"time":"2023-01-01T10:00:00Z","level":"info","msg":"test message","metadata":{"service":"api","version":"1.0"}}`,
			additionalKeys: []string{"metadata"},
			expectedOutput: "Time: 2023-01-01T10:00:00Z\nLevel: info\nMessage: test message\nMetadata: map[service:api version:1.0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &FormatterOptions{
				FormatTemplateFile: tmpFile.Name(),
				TimeKey:            "time",
				LevelKey:           "level",
				MsgKey:             "msg",
				AdditionalKeys:     tt.additionalKeys,
				TimestampFormat:    DEFAULT_TIMESTAMP_FORMAT,
			}

			_, err = NewJsonFormatter(options)
			require.NoError(t, err)

			var entry map[string]interface{}
			err = json.Unmarshal([]byte(tt.jsonInput), &entry)
			require.NoError(t, err)

			logLine := LogLineMapToStruct(entry, options)

			// Execute template
			tmpl, err := template.New("test").Parse(templateContent)
			require.NoError(t, err)

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, logLine)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedOutput, strings.TrimSpace(buf.String()))
		})
	}
}

func TestJsonFormatterWithAdditionalKeys(t *testing.T) {
	tests := []struct {
		name           string
		jsonInput      string
		additionalKeys []string
		template       string
		expectedOutput string
	}{
		{
			name:           "Simple additional key extraction",
			jsonInput:      `{"time":"2023-01-01T10:00:00Z","level":"info","msg":"hello","traceId":"xyz789"}`,
			additionalKeys: []string{"traceId"},
			template:       `{{.Time}} [{{.Level}}] {{.Msg}} (trace: {{.Additional.traceId}})`,
			expectedOutput: "2023-01-01T10:00:00Z [info] hello (trace: xyz789)",
		},
		{
			name:           "Multiple additional keys",
			jsonInput:      `{"time":"2023-01-01T10:00:00Z","level":"error","msg":"failed","userId":"u123","sessionId":"s456","errorCode":500}`,
			additionalKeys: []string{"userId", "sessionId", "errorCode"},
			template:       `{{.Time}} [{{.Level}}] {{.Msg}} | User: {{.Additional.userId}} | Session: {{.Additional.sessionId}} | Code: {{.Additional.errorCode}}`,
			expectedOutput: "2023-01-01T10:00:00Z [error] failed | User: u123 | Session: s456 | Code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary template file
			tmpFile, err := os.CreateTemp("", "test_*.tpl")
			require.NoError(t, err)
			defer func() { _ = os.Remove(tmpFile.Name()) }()

			_, err = tmpFile.WriteString(tt.template)
			require.NoError(t, err)
			_ = tmpFile.Close()

			options := &FormatterOptions{
				FormatTemplateFile: tmpFile.Name(),
				TimeKey:            "time",
				LevelKey:           "level",
				MsgKey:             "msg",
				AdditionalKeys:     tt.additionalKeys,
				TimestampFormat:    DEFAULT_TIMESTAMP_FORMAT,
			}

			_, err = NewJsonFormatter(options)
			require.NoError(t, err)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// We need to test the individual components since Process reads from a file
			var entry map[string]interface{}
			err = json.Unmarshal([]byte(tt.jsonInput), &entry)
			require.NoError(t, err)

			logLine := make(map[string]interface{})
			logLine[options.TimeKey] = ExtractValue(entry, options.TimeKey)
			logLine[options.LevelKey] = ExtractValue(entry, options.LevelKey)
			logLine[options.MsgKey] = ExtractValue(entry, options.MsgKey)

			// Extract additional keys
			for _, key := range options.AdditionalKeys {
				if HasKey(entry, key) {
					logLine[key] = ExtractValue(entry, key)
				}
			}

			// Use template directly for testing
			tmpl, err := template.ParseFiles(tmpFile.Name())
			require.NoError(t, err)
			
			logLineStruct := LogLineMapToStruct(logLine, options)
			err = tmpl.Execute(w, logLineStruct)
			require.NoError(t, err)

			// Restore stdout and read output
			_ = w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)

			assert.Equal(t, tt.expectedOutput, strings.TrimSpace(buf.String()))
		})
	}
}

func TestAdditionalKeysConflictWithStandardKeys(t *testing.T) {
	// Test case where additional key conflicts with a standard key
	jsonInput := `{"time":"2023-01-01T10:00:00Z","level":"info","msg":"test","time":"custom_time_value"}`
	additionalKeys := []string{"time"} // This conflicts with the standard time key

	options := &FormatterOptions{
		TimeKey:         "time",
		LevelKey:        "level",
		MsgKey:          "msg",
		AdditionalKeys:  additionalKeys,
		TimestampFormat: DEFAULT_TIMESTAMP_FORMAT,
	}

	var entry map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &entry)
	require.NoError(t, err)

	logLine := LogLineMapToStruct(entry, options)

	// The standard time processing should take precedence
	assert.NotEmpty(t, logLine.Time)
	// But the additional key should also be available
	assert.Contains(t, logLine.Additional, "time")
}

func TestEmptyAdditionalKeys(t *testing.T) {
	jsonInput := `{"time":"2023-01-01T10:00:00Z","level":"info","msg":"test message"}`

	options := &FormatterOptions{
		TimeKey:         "time",
		LevelKey:        "level",
		MsgKey:          "msg",
		AdditionalKeys:  nil, // No additional keys specified
		TimestampFormat: DEFAULT_TIMESTAMP_FORMAT,
	}

	var entry map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &entry)
	require.NoError(t, err)

	logLine := LogLineMapToStruct(entry, options)

	// Additional map should be empty but not nil
	assert.NotNil(t, logLine.Additional)
	assert.Empty(t, logLine.Additional)
}