package loader

{{range .}}import "funny/{{.}}"
{{end}}

func init() {
	mp = map[string]initFunc {
{{range .}}        "{{.}}": {{.}}.Init,
{{end}}
	}
}
