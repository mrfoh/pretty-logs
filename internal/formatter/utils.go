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
	output := LogLine{
		Additional: make(map[string]interface{}),
	}

	output.Time = FormatTimestamp(ExtractValue(line, options.TimeKey), options.TimestampFormat)
	
	if level := ExtractValue(line, options.LevelKey); level != "" {
		if levelStr, ok := level.(string); ok {
			output.Level = levelStr
		}
	}
	
	if pid := ExtractValue(line, options.PidKey); pid != "" {
		if pidFloat, ok := pid.(float64); ok {
			output.Pid = pidFloat
		}
	}
	
	if name := ExtractValue(line, options.NameKey); name != "" {
		if nameStr, ok := name.(string); ok {
			output.Name = nameStr
		}
	}
	
	if context := ExtractValue(line, options.ContextKey); context != "" {
		if contextStr, ok := context.(string); ok {
			output.Context = contextStr
		}
	}
	
	if msg := ExtractValue(line, options.MsgKey); msg != "" {
		if msgStr, ok := msg.(string); ok {
			output.Msg = msgStr
		}
	}

	// Check for hostname
	if HasKey(line, options.HostnameKey) {
		if hostname := ExtractValue(line, options.HostnameKey); hostname != "" {
			if hostnameStr, ok := hostname.(string); ok {
				output.Hostname = hostnameStr
			}
		}
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

	// Extract additional keys
	for _, key := range options.AdditionalKeys {
		if HasKey(line, key) {
			output.Additional[key] = ExtractValue(line, key)
		}
	}

	return output
}
