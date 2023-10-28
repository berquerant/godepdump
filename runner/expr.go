package runner

import (
	"go/ast"
	"go/parser"

	"github.com/berquerant/godepdump/iox"
)

func ParseExpr() error {
	input, err := iox.ReadStdin()
	if err != nil {
		return err
	}
	expr, err := parser.ParseExpr(input)
	if err != nil {
		return err
	}
	return ast.Print(nil, expr)
}
