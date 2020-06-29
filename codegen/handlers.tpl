// Code generated. DO NOT EDIT.
package {{.Package}}

import "net/http"

{{range $name, $actions := .Structs}}
    func (h *{{$name}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        {{range $action := $actions}}
        case "{{$action.URL}}":
            h.handler{{$action.Name}}(w, r)
        {{end}}
        default:
            w.WriteHeader(http.StatusNotFound)
        }
    }
{{end}}

{{range $name, $actions := .Structs}}
    {{range $action := $actions}}
        func (s *{{$name}}) handler{{$action.Name}}(w http.ResponseWriter, r *http.Request) {
            params :=

            res, err := s.{{$action.Name}}(ctx, params)
            if err != nil {
                 w.WriteHeader(http.StatusInternalServerError)
            }
        }
    {{end}}
{{end}}