package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

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
			log.Println(currFunc.Name)
		}

		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}

		//log.Println(f.(*ast.FuncDecl))

		/* currFunc, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		} else {
			log.Println(currFunc.Name)
		} */

		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			/* currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue
			} */

			log.Println(currType.Name.Name)
		}
	}

	//fmt.Fprintln(`func ServeHTTP(w http.ResponseWriter, req *http.Request) {`)
	//fmt.Fprintln(``)
	//fmt.Fprintln(`}`)
}
