package runner

import (
	"context"

	"github.com/berquerant/godepdump/decls"
	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/find"
	"github.com/berquerant/godepdump/logx"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/write"
)

func SearchDecl(
	ctx context.Context,
	patterns []string,
	name string,
	analyzeLimit int,
) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, patterns...); err != nil {
		return err
	}

	var (
		declList  = def.New().List(loader.List()...)
		searcher  = ref.NewSearcher(declList)
		objectMap = ref.NewObject(loader.List()...)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(analyzeLimit))
		builder   = decls.New(loader, searcher, objectMap, analyzer)
	)

	type Result struct {
		Name string            `json:"name"`
		Pkg  *display.Package  `json:"pkg"`
		Pos  *display.Position `json:"pos"`
		Type *display.Type     `json:"type"`
		Decl *decls.Decl       `json:"decl"`
	}

	for _, ident := range find.IdentFromPackage(name, loader.List()...) {
		r := Result{
			Name: ident.Ident().String(),
			Pkg:  display.NewPackage(ident.Pkg()),
			Pos:  display.NewPosition(ident.Ident().Pos(), ident.Pkg().Fset),
		}
		if obj, ok := objectMap.FindObj(ident.Pkg(), ident.Ident()); ok {
			r.Type = analyzer.Analyze(obj.Type())
		}
		decl, err := builder.Build(r.Pkg.Path, ident.Ident().Pos())
		if err != nil {
			logx.Info("decl not found", logx.Err(err))
		} else {
			r.Decl = decl
		}

		write.JSON(r)
	}

	return nil
}
