package find

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/berquerant/godepdump/astutil"
	"golang.org/x/tools/go/packages"
)

func IdentFromFile(f *ast.File, identName string) (identList []*ast.Ident) {
	ast.Inspect(f, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if id.Name == identName {
				identList = append(identList, id)
			}
		}
		return true
	})
	return
}

//go:generate go run github.com/berquerant/dataclass@v0.3.1 -type "IdentTuple" -field "Ident *ast.Ident|Pkg *packages.Package" -output ident_dataclass_generated.go

func IdentFromPackage(identName string, pkgs ...*packages.Package) (identList []IdentTuple) {
	for _, pkg := range pkgs {
		for _, tree := range pkg.Syntax {
			for _, ident := range IdentFromFile(tree, identName) {
				identList = append(identList, NewIdentTuple(ident, pkg))
			}
		}
	}
	return
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

func Object(pkg *packages.Package, ident *ast.Ident) (types.Object, bool) {
	obj, ok := pkg.TypesInfo.Defs[ident]
	return obj, ok
}
