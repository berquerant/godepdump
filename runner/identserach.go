package runner

import (
	"context"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/find"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/write"
)

type SearchIdent struct {
	Patterns     []string
	Name         *regexp.Regexp
	AnalyzeLimit int
	Exported     bool
	PackageName  *regexp.Regexp
}

func (r *SearchIdent) Run(ctx context.Context) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, r.Patterns...); err != nil {
		return err
	}

	var (
		objectMap = ref.NewObject(loader.List()...)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(r.AnalyzeLimit))
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
		ID   string            `json:"id"`
		Name string            `json:"name"`
		Pkg  *display.Package  `json:"pkg"`
		Pos  *display.Position `json:"pos"`
		Type *display.Type     `json:"type"`
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

		write.JSON(r)
	}
	return nil
}
