package runner

import (
	"context"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/use"
	"github.com/berquerant/godepdump/write"
)

type ListUses struct {
	Patterns     []string
	AnalyzeLimit int
	Exported     bool
	PackageName  *regexp.Regexp
	IdentName    *regexp.Regexp
}

func (r *ListUses) Run(ctx context.Context) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, r.Patterns...); err != nil {
		return err
	}

	var (
		analyzer = display.NewTypeAnalyzer(display.WithLimit(r.AnalyzeLimit))
		filters  = []chanx.Filter[use.Node]{
			func(x use.Node) bool { return packageNameFilter(r.PackageName).Call(x.Pkg()) },
			func(x use.Node) bool { return identExportedFilter(r.Exported).Call(x.Ident()) },
			func(x use.Node) bool { return identNameFilter(r.IdentName).Call(x.Ident()) },
		}
		stream = use.New().List(loader.List()...)
	)
	for _, f := range filters {
		stream.Filter(f)
	}

	type Result struct {
		Name   string            `json:"name"`
		ID     string            `json:"id"`
		Pos    *display.Position `json:"pos"`
		Pkg    *display.Package  `json:"pkg"`
		Type   *display.Type     `json:"type"`
		ObjPkg *display.Package  `json:"objpkg"`
		ObjPos *display.Position `json:"objpos"`
	}

	for useNode := range stream.C() {
		var (
			ident = useNode.Ident()
			pkg   = useNode.Pkg()
			obj   = useNode.Obj()
		)

		v := Result{
			Name: ident.String(),
			Pos:  display.NewPosition(ident.Pos(), pkg.Fset),
			Pkg:  display.NewPackage(pkg),
		}
		if obj != nil {
			v.Type = analyzer.Analyze(obj.Type())
			v.ID = obj.Id()

			if objPkg := obj.Pkg(); objPkg != nil {
				v.ObjPkg = display.NewPackageFromTypes(obj.Pkg())

				if objPkg, ok := loader.GetByPkgPath(obj.Pkg().Path()); ok {
					v.ObjPos = display.NewPosition(obj.Pos(), objPkg.Fset)
				}
			}
		}

		write.JSON(v)
	}

	return nil
}
