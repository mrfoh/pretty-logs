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

// Extract a typed value from a map with a default value.
func ExtractTypedValue[T any](entry map[string]interface{}, key string, defaultValue T) T {
	if value := ExtractValue(entry, key); value != nil && value != "" {
		if typed, ok := value.(T); ok {
			return typed
		}
	}
	return defaultValue
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
	output := LogLine{
		Additional: make(map[string]interface{}),
	}

	output.Time = FormatTimestamp(ExtractValue(line, options.TimeKey), options.TimestampFormat)
	output.Level = ExtractTypedValue(line, options.LevelKey, "")
	output.Pid = ExtractTypedValue(line, options.PidKey, float64(0))
	output.Name = ExtractTypedValue(line, options.NameKey, "")
	output.Context = ExtractTypedValue(line, options.ContextKey, "")
	output.Msg = ExtractTypedValue(line, options.MsgKey, "")
	output.Hostname = ExtractTypedValue(line, options.HostnameKey, "")

	// Extract request and response objects
	output.Req = ExtractMapValue(line, options.RequestKey)
	output.Res = ExtractMapValue(line, options.ResponseKey)

	// Check for error object in log line
	if HasAnyKey(line, options.ErrorObjectKeys) {
		for _, key := range options.ErrorObjectKeys {
			if HasKey(line, key) {
				output.Error = ExtractValue(line, key)
				break
			}
		}
	}

	// Extract additional keys
	for _, key := range options.AdditionalKeys {
		if value := ExtractValue(line, key); value != "" {
			output.Additional[key] = value
		}
	}

	return output
}
