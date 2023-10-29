package runner

import (
	"context"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/decls"
	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/logx"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
	"github.com/berquerant/godepdump/use"
	"github.com/berquerant/godepdump/write"
)

type ListDeps struct {
	Patterns     []string
	AnalyzeLimit int
	Exported     bool
	PackageName  *regexp.Regexp
	IdentName    *regexp.Regexp
}

func (r *ListDeps) Run(ctx context.Context) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, r.Patterns...); err != nil {
		return err
	}

	var (
		declList  = def.New().List(loader.List()...)
		objectMap = ref.NewObject(loader.List()...)
		searcher  = ref.NewSearcher(declList)
		analyzer  = display.NewTypeAnalyzer(display.WithLimit(r.AnalyzeLimit))
		builder   = decls.New(loader, searcher, objectMap, analyzer)
		filters   = []chanx.Filter[use.Node]{
			func(x use.Node) bool { return packageNameFilter(r.PackageName).Call(x.Pkg()) },
			func(x use.Node) bool { return identExportedFilter(r.Exported).Call(x.Ident()) },
			func(x use.Node) bool { return identNameFilter(r.IdentName).Call(x.Ident()) },
		}
		stream = use.New().List(loader.List()...)
	)
	for _, f := range filters {
		stream.Filter(f)
	}

	type Use struct {
		Name string            `json:"name"`
		Pos  *display.Position `json:"pos"`
		Pkg  *display.Package  `json:"pkg"`
		Decl *decls.Decl       `json:"decl"`
	}

	type Def struct {
		ID   string            `json:"id"`
		Pos  *display.Position `json:"pos"`
		Pkg  *display.Package  `json:"pkg"`
		Type *display.Type     `json:"type"`
		Decl *decls.Decl       `json:"decl"`
	}

	type Result struct {
		Use *Use `json:"use"`
		Def *Def `json:"def"`
	}

	for useNode := range stream.C() {
		result := func() *Result {
			r := &Result{
				Use: &Use{
					Name: useNode.Ident().String(),
					Pos:  display.NewPosition(useNode.Ident().Pos(), useNode.Pkg().Fset),
					Pkg:  display.NewPackage(useNode.Pkg()),
				},
				Def: &Def{
					ID:   useNode.Obj().Id(),
					Type: analyzer.Analyze(useNode.Obj().Type()),
				},
			}
			{
				decl, err := builder.Build(r.Use.Pkg.Path, useNode.Ident().Pos())
				if err != nil {
					logx.Debug("use decl not found",
						logx.Err(err),
						logx.S("pkg", r.Use.Pkg.Path),
						logx.S("name", r.Use.Name),
						logx.Any("pos", r.Use.Pos),
					)
				} else {
					r.Use.Decl = decl
				}
			}

			defPkg := useNode.Obj().Pkg()
			if defPkg == nil {
				r.Def.Pkg = display.NewBuiltinPackage()
				return r
			}
			r.Def.Pkg = display.NewPackageFromTypes(defPkg)
			{
				decl, err := builder.Build(r.Def.Pkg.Path, useNode.Obj().Pos())
				if err != nil {
					attrs := []logx.Attr{
						logx.Err(err),
						logx.S("pkg", r.Def.Pkg.Path),
						logx.S("ident", r.Def.ID),
					}
					if pkg, ok := loader.GetByPkgPath(r.Def.Pkg.Path); ok {
						attrs = append(attrs, logx.Any("pos", display.NewPosition(useNode.Obj().Pos(), pkg.Fset)))
					}
					logx.Debug("def decl not found", attrs...)
				} else {
					r.Def.Decl = decl
				}
			}

			return r
		}()
		write.JSON(result)
	}
	return nil
}
