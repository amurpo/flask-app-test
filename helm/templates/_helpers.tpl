# _helpers.tpl

{{- define "my-flask-app-chart.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name }}
{{- end }}

{{- define "my-flask-app-chart.labels" -}}
{{- dict "app" .Chart.Name "release" .Release.Name | toYaml | nindent 4 }}
{{- end }}

{{- define "my-flask-app-chart.selectorLabels" -}}
{{- dict "app" .Chart.Name "release" .Release.Name | toYaml | nindent 6 }}
{{- end }}

{{- define "my-flask-app-chart.containerSettings" -}}
# Define your container settings here
{{- end }}
