package runner

import (
	"context"

	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/find"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/write"
)

func SearchIdent(ctx context.Context, patterns []string, name string, analyzeLimit int) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, patterns...); err != nil {
		return err
	}

	var (
		objectMap = ref.NewObject(loader.List()...)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(analyzeLimit))
	)

	type Result struct {
		ID   string            `json:"id"`
		Name string            `json:"name"`
		Pkg  *display.Package  `json:"pkg"`
		Pos  *display.Position `json:"pos"`
		Type *display.Type     `json:"type"`
	}

	for _, ident := range find.IdentFromPackage(name, loader.List()...) {
		r := Result{
			Name: ident.Ident().String(),
			Pkg:  display.NewPackage(ident.Pkg()),
			Pos:  display.NewPosition(ident.Ident().Pos(), ident.Pkg().Fset),
		}
		if obj, ok := objectMap.FindObj(ident.Pkg(), ident.Ident()); ok {
			r.ID = obj.Id()
			r.Type = analyzer.Analyze(obj.Type())
		}

		write.JSON(r)
	}
	return nil
}
