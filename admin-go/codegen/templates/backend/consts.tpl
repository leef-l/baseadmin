package consts
{{range .Fields}}{{if .IsEnum}}
// {{$.ModelName}}{{.NameCamel}} {{.Label}}
const (
{{- $fieldCamel := .NameCamel -}}
{{- range .EnumValues}}
	{{$.ModelName}}{{$fieldCamel}}{{if .NameIdent}}{{.NameIdent}}{{else}}V{{.Value}}{{end}} = {{if IsNumeric .Value}}{{.Value}}{{else}}"{{.Value}}"{{end}} // {{.Label}}
{{- end}}
)
{{end}}{{end}}
