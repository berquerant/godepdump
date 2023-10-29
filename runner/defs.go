package runner

import (
	"context"
	"go/ast"
	"go/types"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/write"
	"golang.org/x/tools/go/packages"
)

type ListDefs struct {
	Patterns     []string
	AnalyzeLimit int
	Exported     bool
	PackageName  *regexp.Regexp
	IdentName    *regexp.Regexp
}

func (r *ListDefs) Run(ctx context.Context) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, r.Patterns...); err != nil {
		return err
	}

	analyzer := display.NewTypeAnalyzer(display.WithLimit(r.AnalyzeLimit))

	type Result struct {
		Name string            `json:"name"`
		ID   string            `json:"id"`
		Pos  *display.Position `json:"pos"`
		Type *display.Type     `json:"type"`
		Pkg  *display.Package  `json:"pkg"`
	}

	type Elem struct {
		ident *ast.Ident
		pkg   *packages.Package
		obj   types.Object
	}

	var (
		streamC = make(chan *Elem, 100)
		stream  = chanx.NewStream(streamC)
		filters = []chanx.Filter[*Elem]{
			func(x *Elem) bool { return packageNameFilter(r.PackageName).Call(x.pkg) },
			func(x *Elem) bool { return identExportedFilter(r.Exported).Call(x.ident) },
			func(x *Elem) bool { return identNameFilter(r.IdentName).Call(x.ident) },
		}
	)
	for _, f := range filters {
		stream.Filter(f)
	}

	go func() {
		defer close(streamC)
		for _, pkg := range loader.List() {
			for ident, obj := range pkg.TypesInfo.Defs {
				streamC <- &Elem{
					ident: ident,
					pkg:   pkg,
					obj:   obj,
				}
			}
		}
	}()

	for x := range stream.C() {
		p := display.NewPosition(x.ident.Pos(), x.pkg.Fset)
		if x.obj == nil {
			write.JSON(Result{
				Name: x.ident.String(),
				Pos:  p,
			})
			continue
		}
		write.JSON(Result{
			ID:   x.obj.Id(),
			Type: analyzer.Analyze(x.obj.Type()),
			Pkg:  display.NewPackageFromTypes(x.obj.Pkg()),
			Pos:  p,
		})
	}

	return nil
}
