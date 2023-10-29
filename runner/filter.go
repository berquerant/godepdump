package runner

import (
	"go/ast"
	"regexp"

	"github.com/berquerant/godepdump/chanx"
	"github.com/berquerant/godepdump/typesx"
	"golang.org/x/tools/go/packages"
)

func identExportedFilter(enabled bool) chanx.Filter[*ast.Ident] {
	if enabled {
		return typesx.IdentExportedFilter
	}
	return nil
}

func packageNameFilter(r *regexp.Regexp) chanx.Filter[*packages.Package] {
	if r != nil {
		return typesx.PackageNameFilter(r)
	}
	return nil
}

func identNameFilter(r *regexp.Regexp) chanx.Filter[*ast.Ident] {
	if r != nil {
		return typesx.IdentNameFilter(r)
	}
	return nil
}
