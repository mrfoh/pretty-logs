package formatter

import "os"

type FormatterOptions struct {
	FormatTemplateFile string
	OutputFormat       string
	TimestampFormat    string
	LevelKey           string
	TimeKey            string
	PidKey             string
	NameKey            string
	ContextKey         string
	MsgKey             string
	ErrorObjectKeys    []string
	HostnameKey        string
	RequestKey         string
	ResponseKey        string
}

type LogLine struct {
	Time     string
	Level    string
	Pid      float64
	Name     string
	Context  string
	Msg      string
	Error    string
	Req      string
	Hostname string
}

type Formatter interface {
	Process(input *os.File) error
}
