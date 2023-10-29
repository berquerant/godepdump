package typesx

import (
	"go/ast"
	"go/types"
	"regexp"

	"golang.org/x/tools/go/packages"
)

func UniverseFilter(obj types.Object) bool {
	return obj != nil && obj.Pkg() == nil
}

func ObjectExportedFilter(obj types.Object) bool {
	return obj != nil && obj.Exported()
}

func IdentExportedFilter(ident *ast.Ident) bool {
	return ident.IsExported()
}

func RegexpFilter(r *regexp.Regexp) func(string) bool {
	return r.MatchString
}

func PackageNameFilter(r *regexp.Regexp) func(*packages.Package) bool {
	f := RegexpFilter(r)
	return func(pkg *packages.Package) bool {
		return f(pkg.Name)
	}
}

func ObjectNameFilter(r *regexp.Regexp) func(types.Object) bool {
	f := RegexpFilter(r)
	return func(obj types.Object) bool {
		return obj != nil && f(obj.Name())
	}
}

func IdentNameFilter(r *regexp.Regexp) func(*ast.Ident) bool {
	f := RegexpFilter(r)
	return func(ident *ast.Ident) bool {
		return f(ident.Name)
	}
}
