package formatter

var DEFAULT_TIMESTAMP_FORMAT = "2006-01-02 15:04:05.000"

var DEFAULT_OUTPUT_FORMAT = "{time} {level} {pid} --- [{name}] {context} : {msg} {error}"

var ERROR_LIKE_KEYS = []string{"err", "error"}

var MESSAGE_KEY = "msg"

var LEVEL_KEY = "level"

var TIME_KEY = "time"

var PID_KEY = "pid"

var NAME_KEY = "name"

var CONTEXT_KEY = "context"

var HOSTNAME_KEY = "hostname"

var REQUEST_KEY = "req"

var RESPONSE_KEY = "res"
