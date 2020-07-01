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

type action struct {
	URL    string
	Name   string
	Params []queryParams
	Method string
	Auth   bool
}

type queryParams struct {
	Name  string
	Rules []string
}

type validationStruct struct {
	Name   string
	Fields []validationStructFields
}

type validationStructFields struct {
	Type  string
	Rules map[string]interface{}
}

func main() {
	tplVars := make(map[string]interface{})
	structs := make(map[string][]action)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	var validationStructures []validationStruct
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
			action := &action{
				Name:   function.Name.Name,
				Method: http.MethodGet,
			}
			err := json.Unmarshal([]byte(params), action)
			if err != nil {
				log.Fatalln("Unmarshal err: ", err)
			}

			structName := function.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
			structs[structName] = append(structs[structName], *action)
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
							Type:  fmt.Sprint(field.Type),
							Rules: getRules(fieldRules),
						})
					}
				}

				if len(fields) > 0 {
					validationStructures = append(validationStructures, validationStruct{
						Name:   currType.Name.Name,
						Fields: fields,
					})
				}
			}
		}
	}

	/*for _, t := range validationStructures {
		log.Println(t)
	}*/

	tplVars["Package"] = node.Name.Name
	tplVars["Structs"] = structs
	tplVars["ValidationStructs"] = validationStructures

	tpl, err := template.ParseFiles("handlers.tmpl")
	if err != nil {
		log.Fatalln("Template parse err: ", err)
	}

	file, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}

	err = tpl.Execute(file, tplVars)
	if err != nil {
		log.Fatalln(err)
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
