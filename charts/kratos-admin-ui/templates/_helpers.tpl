{{/*
Expand the name of the chart.
*/}}
{{- define "kratos-admin-ui.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kratos-admin-ui.fullname" -}}
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
{{- define "kratos-admin-ui.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kratos-admin-ui.labels" -}}
helm.sh/chart: {{ include "kratos-admin-ui.chart" . }}
{{ include "kratos-admin-ui.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kratos-admin-ui.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kratos-admin-ui.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend labels
*/}}
{{- define "kratos-admin-ui.backend.labels" -}}
{{ include "kratos-admin-ui.labels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "kratos-admin-ui.backend.selectorLabels" -}}
{{ include "kratos-admin-ui.selectorLabels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend labels
*/}}
{{- define "kratos-admin-ui.frontend.labels" -}}
{{ include "kratos-admin-ui.labels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "kratos-admin-ui.frontend.selectorLabels" -}}
{{ include "kratos-admin-ui.selectorLabels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kratos-admin-ui.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kratos-admin-ui.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Backend image
*/}}
{{- define "kratos-admin-ui.backend.image" -}}
{{- $tag := .Values.backend.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.backend.image.repository $tag }}
{{- end }}

{{/*
Frontend image
*/}}
{{- define "kratos-admin-ui.frontend.image" -}}
{{- $tag := .Values.frontend.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.frontend.image.repository $tag }}
{{- end }}

{{/*
Admin secret name
*/}}
{{- define "kratos-admin-ui.adminSecretName" -}}
{{- if .Values.backend.auth.existingSecret }}
{{- .Values.backend.auth.existingSecret }}
{{- else }}
{{- include "kratos-admin-ui.fullname" . }}-admin
{{- end }}
{{- end }}

{{/*
Admin secret key
*/}}
{{- define "kratos-admin-ui.adminSecretKey" -}}
{{- .Values.backend.auth.existingSecretKey | default "password" }}
{{- end }}

