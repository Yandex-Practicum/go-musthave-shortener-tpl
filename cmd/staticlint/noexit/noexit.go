// Package noexit предоставляет анализатор noexit, который запрещает прямой вызов os.Exit в main.
package noexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// NoExitAnalyzer запрещает использование os.Exit в main.
var NoExitAnalyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "запрещает прямой вызов os.Exit в main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Проверяем, что это файл main
		if pass.Pkg.Name() == "main" && pass.Fset.File(file.Pos()).Name() == "main.go" {
			ast.Inspect(file, func(n ast.Node) bool {
				// Поиск функции main
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
					// Поиск os.Exit
					ast.Inspect(fn.Body, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
								switch pkgIdent, ok := sel.X.(*ast.Ident); {
								case pkgIdent.Name == "os" && ok && sel.Sel.Name == "Exit":
									pass.Reportf(call.Pos(), "нельзя использовать os.Exit в main")
								}
							}
						}
						return true
					})
				}
				return true
			})
		}
	}
	return nil, nil
}
