package runner

import (
	"context"

	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/write"
)

func ListDefs(ctx context.Context, patterns []string, analyzeLimit int) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, patterns...); err != nil {
		return err
	}

	analyzer := display.NewTypeAnalyzer(display.WithLimit(analyzeLimit))

	type Result struct {
		Name string            `json:"name`
		ID   string            `json:"id"`
		Pos  *display.Position `json:"pos"`
		Type *display.Type     `json:"type"`
		Pkg  *display.Package  `json:"pkg"`
	}

	for _, pkg := range loader.List() {
		for ident, obj := range pkg.TypesInfo.Defs {
			p := display.NewPosition(ident.Pos(), pkg.Fset)
			if obj == nil {
				write.JSON(Result{
					Name: ident.String(),
					Pos:  p,
				})
				continue
			}
			write.JSON(Result{
				ID:   obj.Id(),
				Type: analyzer.Analyze(obj.Type()),
				Pkg:  display.NewPackageFromTypes(obj.Pkg()),
				Pos:  p,
			})
		}
	}
	return nil
}
