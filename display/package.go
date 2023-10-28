package display

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewPackageFromTypes(pkg *types.Package) *Package {
	return &Package{
		Name: pkg.Name(),
		Path: pkg.Path(),
	}
}

func NewPackage(pkg *packages.Package) *Package {
	return &Package{
		Name: pkg.Name,
		Path: pkg.PkgPath,
	}
}

func NewBuiltinPackage() *Package {
	return &Package{
		Name: "builtin",
		Path: "builtin",
	}
}
