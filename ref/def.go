package ref

import (
	"go/ast"
	"go/types"

	"github.com/berquerant/godepdump/astutil"
	"github.com/berquerant/godepdump/logx"
	"golang.org/x/tools/go/packages"
)

type Object interface {
	FindObj(pkg *packages.Package, ident *ast.Ident) (types.Object, bool)
}

func NewObject(pkgs ...*packages.Package) Object {
	return &object{
		d: newPkgObjectMap(pkgs),
	}
}

type pkgObjectMap map[string]map[*ast.Ident]types.Object

func newPkgObjectMap(pkgs []*packages.Package) pkgObjectMap {
	var (
		d     = make(map[string]map[*ast.Ident]types.Object, len(pkgs))
		count int
	)
	for _, pkg := range pkgs {
		logx.Debug("ref object loaded", logx.S("pkg", pkg.PkgPath))
		d[pkg.ID] = pkg.TypesInfo.Defs
		count += len(d[pkg.ID])
	}
	logx.Info("ref object loaded", logx.I("count", count))
	return d
}

type object struct {
	d pkgObjectMap
}

func (o *object) FindObj(pkg *packages.Package, ident *ast.Ident) (obj types.Object, found bool) {
	defer func() {
		attrs := []logx.Attr{
			logx.S("pkg", pkg.PkgPath),
			logx.S("ident", ident.String()),
			logx.B("found", found),
			logx.Any("pos", astutil.PosString(ident.Pos(), pkg.Fset)),
		}
		logx.Debug("ref find obj", attrs...)
	}()

	objectList, ok := o.d[pkg.ID]
	if !ok {
		return
	}
	obj, found = objectList[ident]
	return
}
