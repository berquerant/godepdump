package find

import (
	"go/ast"
	"go/token"
	"regexp"

	"github.com/berquerant/godepdump/astutil"
	"github.com/berquerant/godepdump/chanx"
	"golang.org/x/tools/go/packages"
)

func IdentFromFile(f *ast.File, identName *regexp.Regexp) chanx.Stream[*ast.Ident] {
	resultC := make(chan *ast.Ident, 100)
	go func() {
		defer close(resultC)
		ast.Inspect(f, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				if identName.MatchString(id.Name) {
					resultC <- id
				}
			}
			return true
		})
	}()
	return chanx.NewStream(resultC)
}

//go:generate go run github.com/berquerant/dataclass@v0.3.1 -type "IdentTuple" -field "Ident *ast.Ident|Pkg *packages.Package" -output ident_dataclass_generated.go

func IdentFromPackage(identName *regexp.Regexp, pkgs ...*packages.Package) chanx.Stream[IdentTuple] {
	resultC := make(chan IdentTuple)
	go func() {
		defer close(resultC)
		for _, pkg := range pkgs {
			for _, tree := range pkg.Syntax {
				for ident := range IdentFromFile(tree, identName).C() {
					resultC <- NewIdentTuple(ident, pkg)
				}
			}
		}
	}()
	return chanx.NewStream(resultC)
}

func ValueSpecIndex(vs *ast.ValueSpec, pos token.Pos) (index int, found bool) {
	if vs.Type != nil {
		if astutil.InPos(vs.Type, pos) {
			found = true
			return
		}
	}
	for i, nm := range vs.Names {
		if astutil.InPos(nm, pos) {
			index = i
			found = true
			return
		}
	}
	for i, v := range vs.Values {
		if astutil.InPos(v, pos) {
			index = i
			found = true
			return
		}
	}
	return
}
