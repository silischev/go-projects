package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

const codegenFlag = "// apigen:api"

type url string

type action string

type structData struct {
	urlsHandlers map[url]action
}

type commentData struct {
	Url  string
	Auth bool
}

var structures = make(map[string]*structData)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out)

	fmt.Fprintln(out, `import "net/http"`)

	for _, f := range node.Decls {
		switch f.(type) {
		case *ast.FuncDecl:
			currFunc := f.(*ast.FuncDecl)

			if currFunc.Doc != nil {
				for _, comment := range currFunc.Doc.List {
					if strings.HasPrefix(comment.Text, codegenFlag) {
						structureName := currFunc.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
						commentData := &commentData{}
						err := json.Unmarshal([]byte(strings.Replace(comment.Text, codegenFlag, "", -1)), commentData)
						if err != nil {
							log.Fatal(err)
						}

						val, exist := structures[structureName]
						if exist {
							val.urlsHandlers[url(commentData.Url)] = action(currFunc.Name.Name)
						} else {
							structure := &structData{}
							structure.urlsHandlers = map[url]action{url(commentData.Url): action(currFunc.Name.Name)}
							structures[structureName] = structure
						}
					}
				}

				//fmt.Println(fmt.Sprintf("%#v", currFunc.Doc.List))
			}
		}

		/* g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		} */

		//log.Println(f.(*ast.FuncDecl))

		// for _, spec := range g.Specs {
		// 	currType, ok := spec.(*ast.TypeSpec)
		// 	if !ok {
		// 		continue
		// 	}

		// 	/* currStruct, ok := currType.Type.(*ast.StructType)
		// 	if !ok {
		// 		continue
		// 	} */

		// 	log.Println(currType.Name.Name)
		// }
	}

	//fmt.Println(fmt.Sprintf("%#v", structures))
	//log.Println(structures)

	for structName, val := range structures {
		fmt.Fprintln(out)
		fmt.Fprintln(out, fmt.Sprintf(`func (structure %s) ServeHTTP(w http.ResponseWriter, req *http.Request) {`, structName))
		fmt.Fprintln(out, "\t"+`switch req.URL.Path {`)
		for url, handler := range val.urlsHandlers {
			fmt.Fprintln(out, fmt.Sprintf("\t"+`case "%s":`, url))
			fmt.Fprintln(out, fmt.Sprintf("\t\t"+`structure.%s(w, req)`, handler))
		}
		fmt.Fprintln(out, "\t"+`}`)
		fmt.Fprintln(out, `}`)
	}
}
