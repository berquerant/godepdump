package packagesx

import (
	"context"

	"github.com/berquerant/godepdump/errorx"
	"github.com/berquerant/godepdump/logx"
	"golang.org/x/tools/go/packages"
)

type Loader interface {
	Load(ctx context.Context, patterns ...string) error
	List() []*packages.Package
	GetByPkgPath(pkgPath string) (*packages.Package, bool)
}

func New() Loader {
	return &loader{
		pkgs:         []*packages.Package{},
		pkgPathIndex: map[string]*packages.Package{},
	}
}

type loader struct {
	pkgs         []*packages.Package
	pkgPathIndex map[string]*packages.Package
}

const loadMode = packages.NeedTypesInfo | packages.NeedTypes | packages.NeedName |
	packages.NeedSyntax | packages.NeedImports

func (l *loader) Load(ctx context.Context, patterns ...string) error {
	logx.Info("load", logx.SS("patterns", patterns))

	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Mode:    loadMode,
	}, patterns...)
	if err != nil {
		return errorx.Errorf(err, "load %s", patterns)
	}

	l.pkgs = append(l.pkgs, pkgs...)
	for _, p := range pkgs {
		logx.Info("loaded", logx.S("package", p.PkgPath))
		l.pkgPathIndex[p.PkgPath] = p
	}
	return nil
}

func (l *loader) List() []*packages.Package { return l.pkgs }

func (l *loader) GetByPkgPath(pkgPath string) (*packages.Package, bool) {
	p, ok := l.pkgPathIndex[pkgPath]
	return p, ok
}
