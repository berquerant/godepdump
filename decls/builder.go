package decls

import (
	"encoding/json"
	"errors"
	"go/token"

	"github.com/berquerant/godepdump/astutil"
	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/errorx"
	"github.com/berquerant/godepdump/find"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/ref"
)

type Decl struct {
	Kind  Kind              `json:"kind"`
	ID    string            `json:"id"`
	Begin *display.Position `json:"begin"`
	End   *display.Position `json:"end"`
	Type  *display.Type     `json:"type"`
}

type Builder interface {
	// Build finds the top-level function, type and variables definition to which pos belongs.
	Build(pkgPath string, pos token.Pos) (*Decl, error)
}

func New(
	loader packagesx.Loader,
	searcher ref.Searcher,
	objectMap ref.Object,
	analyzer display.TypeAnalyzer,
) Builder {
	return &builder{
		loader:    loader,
		searcher:  searcher,
		objectMap: objectMap,
		analyzer:  analyzer,
	}
}

type builder struct {
	loader    packagesx.Loader
	searcher  ref.Searcher
	objectMap ref.Object
	analyzer  display.TypeAnalyzer
}

func (b *builder) Build(pkgPath string, pos token.Pos) (*Decl, error) {
	pkg, ok := b.loader.GetByPkgPath(pkgPath)
	if !ok {
		return nil, errorx.Errorf(ErrPackageNotFound, "path: %s", pkgPath)
	}
	node, ok := b.searcher.Search(pkg, pos)
	if !ok {
		return nil, errorx.Errorf(ErrObjectNotFound, "path: %s, pos: %s", pkgPath, astutil.PosString(pos, pkg.Fset))
	}

	r := &Decl{
		Begin: display.NewPosition(node.Pos(), pkg.Fset),
		End:   display.NewPosition(node.End(), pkg.Fset),
	}

	switch node := node.(type) {
	case *def.TypeSpec:
		r.Kind = KtypeSpec
		r.ID = node.Name.String()
		if obj, ok := b.objectMap.FindObj(node.Pkg(), node.Name); ok {
			r.Type = b.analyzer.Analyze(obj.Type())
		}
	case *def.FuncDecl:
		r.Kind = KfuncDecl
		r.ID = node.Name.String()
		if obj, ok := b.objectMap.FindObj(node.Pkg(), node.Name); ok {
			r.Type = b.analyzer.Analyze(obj.Type())
		}
	case *def.ValueSpec:
		r.Kind = KvalueSpec
		if index, ok := find.ValueSpecIndex(node.ValueSpec, pos); ok {
			name := node.Names[index]
			r.ID = name.String()
			if obj, ok := b.objectMap.FindObj(node.Pkg(), name); ok {
				r.Type = b.analyzer.Analyze(obj.Type())
			}
		}
	}

	return r, nil
}

var (
	ErrPackageNotFound = errors.New("PackageNotFound")
	ErrObjectNotFound  = errors.New("ObjectNotFound")
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Kind -output kind_stringer_generated.go

type Kind int

const (
	KtypeSpec Kind = iota
	KfuncDecl
	KvalueSpec
)

func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}
