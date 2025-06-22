{{ .Time }} {{ .Level | printf "%-5s" }}{{- if .Pid }} ({{ .Pid }}) {{- end -}}
{{- if or .Name .Context }} --- [{{- if .Name }}{{ .Name }}{{- end }}
    {{- if and .Name .Context }}/{{- end }}
    {{- if .Context }}{{ .Context }}{{- end }}] --- {{- end -}}
{{- if .Msg }} {{ .Msg }}{{- end -}}
{{- if .Additional -}}
    {{- range $key, $value := .Additional -}}
        {{- printf " %s=%v" $key $value -}}
    {{- end -}}
{{- end -}}
{{- if .Req }}
  Request: {{ .Req }}
{{- end -}}
{{- if .Error }}
  Error: {{ .Error }}
{{- end }}