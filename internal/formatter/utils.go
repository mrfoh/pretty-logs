package formatter

import (
	"encoding/json"
	"fmt"
	"time"
)

func IsErrorLine(entry *map[string]interface{}, keys []string) bool {
	for _, key := range keys {
		if _, ok := (*entry)[key]; ok {
			return true
		}
	}
	return false
}

func HasRequest(entry *map[string]interface{}, key string) bool {
	if _, ok := (*entry)[key]; ok {
		return true
	}
	return false
}

func HasError(entry *map[string]interface{}, keys []string) bool {
	for _, key := range keys {
		if _, ok := (*entry)[key]; ok {
			return true
		}
	}
	return false
}

func FormatTimestamp(timestamp interface{}, format string) string {
	switch t := timestamp.(type) {
	case string:
		return t
	case float64:
		return time.UnixMilli(int64(t)).Format(format)
	default:
		return ""
	}
}

func FormatError(err interface{}) string {
	if err == nil {
		return ""
	}

	switch e := err.(type) {
	case map[string]interface{}:
		// Check for common error fields
		// if msg, ok := e["message"].(string); ok {
		// 	return msg
		// }

		// if resp, ok := e["response"].(map[string]interface{}); ok {
		// 	if msg, ok := resp["message"].(string); ok {
		// 		return msg
		// 	}
		// }
		// Fallback to marshaling the entire error object
		if errBytes, err := json.Marshal(e); err == nil {
			return string(errBytes)
		}
	}

	return fmt.Sprintf("%v", err)
}

func ExtractValue(entry map[string]interface{}, key string) interface{} {
	value, ok := entry[key]
	if !ok {
		return ""
	}
	return value
}

func ExtractMapValue(entry map[string]interface{}, key string) string {
	value, ok := (entry)[key].(map[string]interface{})
	if !ok {
		return ""
	}

	if valueBytes, err := json.Marshal(value); err == nil {
		return string(valueBytes)
	}

	return fmt.Sprintf("%v", value)
}

func LogLineMapToStruct(line map[string]interface{}, options *FormatterOptions) LogLine {
	output := LogLine{}

	output.Time = FormatTimestamp(ExtractValue(line, options.TimeKey), options.TimestampFormat)
	output.Level = ExtractValue(line, options.LevelKey).(string)
	output.Pid = ExtractValue(line, options.PidKey).(float64)
	output.Name = ExtractValue(line, options.NameKey).(string)
	output.Context = ExtractValue(line, options.ContextKey).(string)
	output.Msg = ExtractValue(line, options.MsgKey).(string)

	if HasRequest(&line, options.RequestKey) {
		output.Req = ExtractMapValue(line, options.RequestKey)
	}

	if HasError(&line, options.ErrorObjectKeys) {
		output.Error = ExtractValue(line, "error").(string)
	}

	return output
}
