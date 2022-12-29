// Code generated by codegen; DO NOT EDIT.

package plugin

import (
	"github.com/cloudquery/plugin-sdk/schema"
{{- range . }}
	"github.com/cloudquery/cloudquery/plugins/source/snyk/resources/services/{{ .Service }}"
{{- end }}
)

func tables() []*schema.Table {
	return []*schema.Table{
    {{- range . }}
        {{ .Service }}.{{ .SubService | ToCamel }}(),
    {{- end }}
	}
}
