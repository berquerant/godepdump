package runner

import (
	"context"

	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/use"
	"github.com/berquerant/godepdump/write"
)

func ListUses(ctx context.Context, patterns []string, analyzeLimit int) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, patterns...); err != nil {
		return err
	}

	analyzer := display.NewTypeAnalyzer(display.WithLimit(analyzeLimit))

	type Result struct {
		ID     string            `json:"id"`
		Pos    *display.Position `json:"pos"`
		Pkg    *display.Package  `json:"pkg"`
		Type   *display.Type     `json:"type"`
		ObjPkg *display.Package  `json:"objpkg"`
		ObjPos *display.Position `json:"objpos"`
	}

	for useNode := range use.New().List(loader.List()...).C() {
		var (
			ident = useNode.Ident()
			pkg   = useNode.Pkg()
			obj   = useNode.Obj()
		)

		v := Result{
			ID:  ident.String(),
			Pos: display.NewPosition(ident.Pos(), pkg.Fset),
			Pkg: display.NewPackage(pkg),
		}
		if obj != nil {
			v.Type = analyzer.Analyze(obj.Type())

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
