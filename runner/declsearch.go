package runner

import (
	"context"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/decls"
	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/find"
	"github.com/berquerant/godepdump/logx"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/write"
)

type SearchDecl struct {
	Patterns     []string
	Name         *regexp.Regexp
	AnalyzeLimit int
	Exported     bool
	PackageName  *regexp.Regexp
}

func (r *SearchDecl) Run(ctx context.Context) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, r.Patterns...); err != nil {
		return err
	}

	var (
		declList  = def.New().List(loader.List()...)
		searcher  = ref.NewSearcher(declList)
		objectMap = ref.NewObject(loader.List()...)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(r.AnalyzeLimit))
		builder   = decls.New(loader, searcher, objectMap, analyzer)
		stream    = find.IdentFromPackage(r.Name, loader.List()...)
		filters   = []chanx.Filter[find.IdentTuple]{
			func(x find.IdentTuple) bool { return packageNameFilter(r.PackageName).Call(x.Pkg()) },
			func(x find.IdentTuple) bool { return identExportedFilter(r.Exported).Call(x.Ident()) },
		}
	)
	for _, f := range filters {
		stream.Filter(f)
	}

	type Result struct {
		Name string            `json:"name"`
		ID   string            `json:"id"`
		Pkg  *display.Package  `json:"pkg"`
		Pos  *display.Position `json:"pos"`
		Type *display.Type     `json:"type"`
		Decl *decls.Decl       `json:"decl"`
	}

	for ident := range stream.C() {
		r := Result{
			Name: ident.Ident().String(),
			Pkg:  display.NewPackage(ident.Pkg()),
			Pos:  display.NewPosition(ident.Ident().Pos(), ident.Pkg().Fset),
		}
		if obj, ok := objectMap.FindObj(ident.Pkg(), ident.Ident()); ok {
			r.ID = obj.Id()
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
