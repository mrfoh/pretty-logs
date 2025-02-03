package formatter

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func OpenTestLogFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

func CaptureOutput(f func()) string {
	// Save the original os.Stdout.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function whose output we wish to capture.
	f()

	// Close writer and restore os.Stdout.
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()
	return buf.String()
}

func TestJsonFormatter(t *testing.T) {
	t.Run("throws an error if template file is not found", func(t *testing.T) {
		formatter, err := NewJsonFormatter(&FormatterOptions{
			FormatTemplateFile: "fixtures/some_file.tpl",
			ErrorObjectKeys:    []string{"err", "error"},
			TimeKey:            "time",
			LevelKey:           "level",
			PidKey:             "pid",
			NameKey:            "name",
			ContextKey:         "context",
			MsgKey:             "msg",
			TimestampFormat:    "2006-01-02 15:04:05.000",
			RequestKey:         "req",
			ResponseKey:        "res",
		})

		assert.NotNil(t, err)
		assert.Error(t, err, "template file does not exist: fixtures/some_file.tpl")
		assert.Nil(t, formatter)
	})
}
