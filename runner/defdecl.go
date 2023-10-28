package runner

import (
	"context"
	"go/ast"

	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/write"
)

func ListDefDecls(ctx context.Context, patterns []string, analyzeLimit int) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, patterns...); err != nil {
		return err
	}

	type Result struct {
		Kind  string            `json:"kind"`
		Name  string            `json:"name"`
		Begin *display.Position `json:"begin"`
		End   *display.Position `json:"end"`
		Pkg   *display.Package  `json:"pkg"`
		Type  *display.Type     `json:"type"`
	}

	var (
		declList  = def.New().List(loader.List()...)
		objectMap = ref.NewObject(loader.List()...)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(analyzeLimit))

		newResult = func(kind string, ident *ast.Ident, node def.Node) *Result {
			obj, _ := objectMap.FindObj(node.Pkg(), ident)

			return &Result{
				Kind:  kind,
				Name:  ident.String(),
				Begin: display.NewPosition(node.Pos(), node.Pkg().Fset),
				End:   display.NewPosition(node.End(), node.Pkg().Fset),
				Pkg:   display.NewPackage(node.Pkg()),
				Type:  analyzer.Analyze(obj.Type()),
			}
		}
	)

	for _, decl := range declList {
		switch decl := decl.(type) {
		case *def.ValueSpec:
			for _, id := range decl.Names {
				r := newResult(
					"ValueSpec",
					id,
					decl,
				)
				write.JSON(r)
			}
		case *def.TypeSpec:
			write.JSON(newResult(
				"TypeSpec",
				decl.Name,
				decl,
			))
		case *def.FuncDecl:
			write.JSON(newResult(
				"FuncDecl",
				decl.Name,
				decl,
			))
		}
	}
	return nil
}
