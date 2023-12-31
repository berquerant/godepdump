// Code generated by "dataclass -type IdentTuple -field Ident *ast.Ident|Pkg *packages.Package -output ident_dataclass_generated.go"; DO NOT EDIT.

package find

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type IdentTuple interface {
	Ident() *ast.Ident
	Pkg() *packages.Package
}
type identTuple struct {
	ident *ast.Ident
	pkg   *packages.Package
}

func (s *identTuple) Ident() *ast.Ident      { return s.ident }
func (s *identTuple) Pkg() *packages.Package { return s.pkg }
func NewIdentTuple(
	ident *ast.Ident,
	pkg *packages.Package,
) IdentTuple {
	return &identTuple{
		ident: ident,
		pkg:   pkg,
	}
}
