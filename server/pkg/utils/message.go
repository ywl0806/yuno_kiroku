package utils

import (
	"bytes"
	"text/template"
)

func GetMessage(message string, fields map[string]string) string {
	t := template.Must(template.New("").Parse(message))
	var buf bytes.Buffer
	t.Execute(&buf, fields)
	return buf.String()
}
