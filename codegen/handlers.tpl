// Code generated. DO NOT EDIT.

package {{.Package}}

import "net/http"

{{range $name, $action := .Structs}}
	{{$name}}
	{{$action}}
{{end}}