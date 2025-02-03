{{- define "levelDisplay" -}}
    {{- if . -}}
        {{- . | printf "%-5s" -}}
    {{- end -}}
{{- end -}}

{{- define "metadata" -}}
    {{- if . -}}
        {{- if or .Name .Context -}}[{{- if .Name }}{{ .Name }}{{- end }}
            {{- if and .Name .Context }}/{{- end }}
            {{- if .Context }}{{ .Context }}{{- end }}]{{- end -}}
    {{- end -}}
{{- end -}}

{{- with . }}
{{ .Time }} {{ template "levelDisplay" .Level }}{{- if .Pid }} ({{ .Pid }}) {{- end -}}
{{- if or .Name .Context }} --- {{ template "metadata" . }} --- {{- end -}}
{{- if .Msg }} {{ .Msg }}{{- end -}}
{{- if .Req }}
  Request: {{ .Req }}
{{- end -}}
{{- if .Error }}
  Error: {{ .Error }}
{{- end -}}
{{ end }}