package formatter

import (
	"encoding/json"
	"fmt"
	"time"
)

// Check if a map contains any of the provided keys.
func HasAnyKey(entry map[string]interface{}, keys []string) bool {
	for _, key := range keys {
		if _, ok := entry[key]; ok {
			return true
		}
	}
	return false
}

// Check if a map contains a key.
func HasKey(entry map[string]interface{}, key string) bool {
	_, ok := (entry)[key]
	return ok
}

// Format a timestamp to a string.
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

// Extract an error from a log line.
func ExtractError(err interface{}) string {
	if err == nil {
		return ""
	}

	switch e := err.(type) {
	case map[string]interface{}:
		if errBytes, err := json.Marshal(e); err == nil {
			return string(errBytes)
		}
	case string:
		return e
	}

	return fmt.Sprintf("%v", err)
}

// Extract a value from a map.
func ExtractValue(entry map[string]interface{}, key string) interface{} {
	value, ok := entry[key]
	if !ok {
		return ""
	}
	return value
}

// Extract a map value from a map.
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

// Convert a map to a struct.
func LogLineMapToStruct(line map[string]interface{}, options *FormatterOptions) LogLine {
	output := LogLine{}

	output.Time = FormatTimestamp(ExtractValue(line, options.TimeKey), options.TimestampFormat)
	output.Level = ExtractValue(line, options.LevelKey).(string)
	output.Pid = ExtractValue(line, options.PidKey).(float64)
	output.Name = ExtractValue(line, options.NameKey).(string)
	output.Context = ExtractValue(line, options.ContextKey).(string)
	output.Msg = ExtractValue(line, options.MsgKey).(string)

	// Check for hostname
	if HasKey(line, options.HostnameKey) {
		output.Hostname = ExtractValue(line, options.HostnameKey).(string)
	}

	// Check for request object
	if HasKey(line, options.RequestKey) {
		output.Req = ExtractMapValue(line, options.RequestKey)
	}

	// Check for response object
	if HasKey(line, options.ResponseKey) {
		output.Res = ExtractMapValue(line, options.ResponseKey)
	}

	// Check for error object in log line
	if HasAnyKey(line, options.ErrorObjectKeys) {
		for _, key := range options.ErrorObjectKeys {
			if HasKey(line, key) {
				output.Error = ExtractValue(line, key)
				break
			}
		}
	}

	return output
}
