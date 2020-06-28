package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const codegenPrefix = "apigen:api"

type action struct {
	URL    string
	Method string
	Auth   bool
}

func main() {
	tplVars := make(map[string]interface{})
	structs := make(map[string][]action)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range node.Decls {
		function, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		comment := function.Doc
		if comment == nil {
			continue
		}

		if !strings.HasPrefix(comment.Text(), codegenPrefix) {
			continue
		}

		params := comment.Text()[len(codegenPrefix):len(comment.Text())]
		action := &action{Method: http.MethodGet}
		err := json.Unmarshal([]byte(params), action)
		if err != nil {
			log.Fatalln("Unmarshal err: ", err)
		}

		structName := function.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		structs[structName] = append(structs[structName], *action)
	}

	tplVars["Package"] = node.Name.Name
	tplVars["Structs"] = structs

	tpl, err := template.ParseFiles("handlers.tpl")
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
