package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	codegenPrefix    = "apigen:api"
	structCodegenTag = "`apivalidator:"

	RuleRequired  = "required"
	RuleMin       = "min"
	RuleMax       = "max"
	RuleParamName = "paramname"
	RuleEnum      = "enum"
	RuleDefault   = "default"
)

type methodCodegenParams struct {
	URL              string
	ValidationStruct string
	MethodName       string
	Params           []queryParams
	HTTPMethod       string
	Auth             bool
}

type queryParams struct {
	Name  string
	Rules []string
}

type validationStructFields struct {
	Name  string
	Type  string
	Rules map[string]interface{}
}

func main() {
	tplVars := make(map[string]interface{})
	structs := make(map[string][]methodCodegenParams)

	file, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}

	commonPart := `
		// Code generated automatically. DO NOT EDIT.
		package {{.Package}}
		
		import (
			"log"
			"net/http"
			"unicode/utf8"
			"strconv"
		)
		
		{{range $name, $actions := .Structs}}
			func (h *{{$name}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
					{{range $action := $actions}}
						case "{{$action.URL}}":
							h.handler{{$action.MethodName}}(w, r)
					{{end}}
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}
		{{end}}
	`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	validationStructures := make(map[string][]validationStructFields)
	for _, f := range node.Decls {
		switch f.(type) {
		case *ast.FuncDecl:
			function := f.(*ast.FuncDecl)

			comment := function.Doc
			if comment == nil {
				continue
			}

			if !strings.HasPrefix(comment.Text(), codegenPrefix) {
				continue
			}

			params := comment.Text()[len(codegenPrefix):len(comment.Text())]
			methodCodegenParams := &methodCodegenParams{
				MethodName:       function.Name.Name,
				HTTPMethod:       http.MethodGet,
				ValidationStruct: fmt.Sprint(function.Type.Params.List[1].Type),
			}
			err := json.Unmarshal([]byte(params), methodCodegenParams)
			if err != nil {
				log.Fatalln("Unmarshal err: ", err)
			}

			structName := function.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
			structs[structName] = append(structs[structName], *methodCodegenParams)
		case *ast.GenDecl:
			g := f.(*ast.GenDecl)
			for _, spec := range g.Specs {
				currType, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				currStruct, ok := currType.Type.(*ast.StructType)
				if !ok {
					continue
				}

				var fields []validationStructFields
				for _, field := range currStruct.Fields.List {
					if field.Tag != nil && strings.HasPrefix(field.Tag.Value, structCodegenTag) {
						rexp := regexp.MustCompile(structCodegenTag + `"(.*)"`)
						matches := rexp.FindStringSubmatch(field.Tag.Value)

						fieldRules := strings.Split(matches[1], ",")

						fields = append(fields, validationStructFields{
							Name:  fmt.Sprint(field.Names[0].Name),
							Type:  fmt.Sprint(field.Type),
							Rules: getRules(fieldRules),
						})
					}
				}

				if len(fields) > 0 {
					validationStructures[currType.Name.Name] = fields
				}
			}
		}
	}

	tplVars["Package"] = node.Name.Name
	tplVars["Structs"] = structs

	tpl := template.Must(template.New("").Parse(commonPart))
	err = tpl.Execute(file, tplVars)
	if err != nil {
		log.Fatalln(err)
	}

	for name, actions := range structs {
		out := ""

		for _, action := range actions {
			out += fmt.Sprintf("func (s *%s) handler%s(w http.ResponseWriter, r *http.Request) {\n", name, action.MethodName)

			out += fmt.Sprintf(`
				err := r.ParseForm()
				if err != nil {
					log.Fatalln("Error parse query: ", err)
				}
			`)

			out += fmt.Sprintf("params := %s{} \n", action.ValidationStruct)
			out += "var rawVal string"

			for _, field := range validationStructures[action.ValidationStruct] {
				param := strings.ToLower(field.Name)

				out += fmt.Sprintf(`
					rawVal = ""

					if len(r.Form["%s"]) != 0 {
						rawVal = r.Form["%s"][0]
					}
				`, param, param)

				switch field.Type {
				case "string":
					out += fmt.Sprintf(`
						%s := rawVal
					`, param)
				case "int":
					out += fmt.Sprintf(`
						%s, err := strconv.Atoi(rawVal)
						if err != nil {
						
						}
					`, param)
				}

				for rule, val := range field.Rules {
					switch rule {
					case RuleRequired:
						out += fmt.Sprintf(`
							if rawVal == "" {}
						`)
					case RuleMin:
						if field.Type == "string" {
							out += fmt.Sprintf(`
								if utf8.RuneCountInString(%s) < %s {}
							`, param, val)
						} else {
							out += fmt.Sprintf(`
								if %s < %s {}
							`, param, val)
						}
					case RuleMax:
						if field.Type == "string" {
							out += fmt.Sprintf(`
								if utf8.RuneCountInString(%s) > %s {}
							`, param, val)
						} else {
							out += fmt.Sprintf(`
								if %s > %s {}
							`, param, val)
						}
					}
				}

				out += fmt.Sprintf("params.%s = %s \n", field.Name, param)
			}

			out += fmt.Sprintf(`
				res, err := s.%s(r.Context(), params)
				if err != nil {
					 w.WriteHeader(http.StatusInternalServerError)
				}
			`, action.MethodName)

			out += "\n}\n"
		}

		fmt.Fprintln(file, out)
	}
}

func getRules(rawRules []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, rule := range rawRules {
		var val interface{}
		data := strings.Split(rule, "=")
		ruleName := data[0]

		if len(data) > 1 {
			val = data[1]
		}

		if strings.HasPrefix(rule, RuleEnum) {
			val = strings.Split(val.(string), "|")
		}

		result[ruleName] = val
	}

	return result
}
