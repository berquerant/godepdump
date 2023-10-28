package runner

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/berquerant/godepdump/iox"
)

func ParseFile() error {
	input, err := iox.ReadStdin()
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", input, parser.SkipObjectResolution)
	if err != nil {
		return err
	}
	return ast.Print(fset, f)
}
