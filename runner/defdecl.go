package runner

import (
	"context"
	"go/ast"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/write"
)

type ListDefDecls struct {
	Patterns     []string
	AnalyzeLimit int
	Exported     bool
	PackageName  *regexp.Regexp
	IdentName    *regexp.Regexp
}

func (r *ListDefDecls) Run(ctx context.Context) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, r.Patterns...); err != nil {
		return err
	}

	type Result struct {
		Kind  string            `json:"kind"`
		Name  string            `json:"name"`
		ID    string            `json:"id"`
		Begin *display.Position `json:"begin"`
		End   *display.Position `json:"end"`
		Pkg   *display.Package  `json:"pkg"`
		Type  *display.Type     `json:"type"`
	}

	var (
		declList  = def.New().List(loader.List()...)
		objectMap = ref.NewObject(loader.List()...)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(r.AnalyzeLimit))

		newResult = func(kind string, ident *ast.Ident, node def.Node) *Result {
			obj, _ := objectMap.FindObj(node.Pkg(), ident)

			return &Result{
				Kind:  kind,
				Name:  ident.String(),
				ID:    obj.Id(),
				Begin: display.NewPosition(node.Pos(), node.Pkg().Fset),
				End:   display.NewPosition(node.End(), node.Pkg().Fset),
				Pkg:   display.NewPackage(node.Pkg()),
				Type:  analyzer.Analyze(obj.Type()),
			}
		}
	)

	type Elem struct {
		kind  string
		ident *ast.Ident
		node  def.Node
	}
	var (
		streamC = make(chan *Elem, 100)
		stream  = chanx.NewStream(streamC)
		filters = []chanx.Filter[*Elem]{
			func(x *Elem) bool { return packageNameFilter(r.PackageName).Call(x.node.Pkg()) },
			func(x *Elem) bool { return identExportedFilter(r.Exported).Call(x.ident) },
			func(x *Elem) bool { return identNameFilter(r.IdentName).Call(x.ident) },
		}
	)
	for _, f := range filters {
		stream.Filter(f)
	}

	go func() {
		defer close(streamC)
		for _, decl := range declList {
			switch decl := decl.(type) {
			case *def.ValueSpec:
				for _, id := range decl.Names {
					streamC <- &Elem{
						kind:  "ValueSpec",
						ident: id,
						node:  decl,
					}
				}
			case *def.TypeSpec:
				streamC <- &Elem{
					kind:  "TypeSpec",
					ident: decl.Name,
					node:  decl,
				}
			case *def.FuncDecl:
				streamC <- &Elem{
					kind:  "FuncDecl",
					ident: decl.Name,
					node:  decl,
				}
			}
		}
	}()

	for x := range stream.C() {
		write.JSON(newResult(x.kind, x.ident, x.node))
	}

	return nil
}
