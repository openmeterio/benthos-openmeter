{{/*
Expand the name of the chart.
*/}}
{{- define "benthos-openmeter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "benthos-openmeter.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "benthos-openmeter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "benthos-openmeter.labels" -}}
helm.sh/chart: {{ include "benthos-openmeter.chart" . }}
{{ include "benthos-openmeter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "benthos-openmeter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "benthos-openmeter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "benthos-openmeter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "benthos-openmeter.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified component name from the full app name and a component name.
We truncate the full name at 63 - 1 (last dash) - len(component name) chars because some Kubernetes name fields are limited to this (by the DNS naming spec)
and we want to make sure that the component is included in the name.

Usage: {{ include "benthos-openmeter.componentName" (list . "component") }}
*/}}
{{- define "benthos-openmeter.componentName" -}}
{{- $global := index . 0 -}}
{{- $component := index . 1 | trimPrefix "-" -}}
{{- printf "%s-%s" (include "benthos-openmeter.fullname" $global | trunc (sub 62 (len $component) | int) | trimSuffix "-" ) $component | trimSuffix "-" -}}
{{- end -}}

{{/*
Create args for the deployment
*/}}
{{- define "benthos-openmeter.args" -}}
{{- if .Values.config -}}
["benthos", "-c", "/etc/benthos/config.yaml"]
{{- else if .Values.useExistingConfigFile -}}
["benthos", "-c", "{{ .Values.useExistingConfigFile }}"]
{{- else if .Values.useExample }}
{{- if eq .Values.useExample "http-server" -}}
["benthos", "streams", "--no-api", "/etc/benthos/examples/http-server/input.yaml", "/etc/benthos/examples/http-server/output.yaml"]
{{- else if eq .Values.useExample "kubernetes-pod-exec-time" -}}
["benthos", "-c", "/etc/benthos/examples/kubernetes-pod-exec-time/config.yaml"]
{{- else }}
{{- fail (printf "Invalid example '%s" .Values.useExample) }}
{{- end }}
{{- else }}
{{- fail "One of 'config', 'useExistingConfigFile' or 'useExample' is required" }}
{{- end }}
{{- end }}
