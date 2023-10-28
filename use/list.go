package use

import (
	"github.com/berquerant/godepdump/chanx"
	"golang.org/x/tools/go/packages"
)

//go:generate go run github.com/berquerant/dataclass@v0.3.1 -type "Node" -field "Ident *ast.Ident|Obj types.Object|Pkg *packages.Package" -output list_dataclass_generated.go

type Lister interface {
	List(pkgs ...*packages.Package) chanx.Stream[Node]
}

func New() Lister {
	return &lister{}
}

type lister struct{}

func (*lister) List(pkgs ...*packages.Package) chanx.Stream[Node] {
	resultC := make(chan Node, 100)
	go func() {
		defer close(resultC)
		for _, pkg := range pkgs {
			for ident, obj := range pkg.TypesInfo.Uses {
				resultC <- NewNode(ident, obj, pkg)
			}
		}
	}()
	return chanx.NewStream(resultC)
}
