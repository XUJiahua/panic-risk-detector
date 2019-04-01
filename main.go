package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

type Report struct {
	FuncName string
	FuncPos  token.Pos
	FileName string
	Message  string
}

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "example/test.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// cache all functions
	funcMap := make(map[string]*ast.FuncDecl)
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcMap[fn.Name.Name] = fn
	}

	var reports []Report

	// walk node via depth first search
	ast.Inspect(node, func(n ast.Node) bool {
		// detect go statement
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}

		var report Report

		var funcName string
		switch v := goStmt.Call.Fun.(type) {
		case *ast.Ident:
			funcName = v.Name
			report.FuncName = v.Name
			report.FuncPos = v.NamePos
		default:
			log.Println("unexpected...")
			return true
		}

		// TODO：go 匿名函数 没考虑
		// get func from cache
		fn, ok := funcMap[funcName]
		if !ok {
			// unexpected
			report.Message = "func not found"
			reports = append(reports, report)
			return true
		}

		// empty body, ignore
		if len(fn.Body.List) == 0 {
			return true
		}

		deferStmt, ok := fn.Body.List[0].(*ast.DeferStmt)
		// first statement should be defer statement
		if !ok {
			report.Message = "first statement is not defer/recover statement"
			reports = append(reports, report)
			return true
		}

		// get source code of deferStmt
		var buf bytes.Buffer
		err := printer.Fprint(&buf, fset, deferStmt)
		if err != nil {
			panic(err)
		}
		//fmt.Println(buf.String())

		// check recover() in source code
		if !strings.Contains(buf.String(), "recover()") {
			report.Message = "recover() not exist"
			reports = append(reports, report)
		}

		return true
	})

	if len(reports) == 0 {
		fmt.Println("success")
		os.Exit(0)
	}

	fmt.Println("panic risk detected:")
	for _, report := range reports {
		fmt.Printf("=> function/method: %v, risk: %v\n", report.FuncName, report.Message)
	}
	os.Exit(2)
}
