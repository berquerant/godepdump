package def

import (
	"go/ast"

	"github.com/berquerant/godepdump/astutil"
	"github.com/berquerant/godepdump/logx"
	"golang.org/x/tools/go/packages"
)

// Node is [ValueSpec] | [FuncDecl] | [TypeSpec].
type Node interface {
	IsNode()
	ast.Node
	Pkg() *packages.Package
}

//go:generate go run github.com/berquerant/marker@v0.1.4 -method IsNode -output list_marker_generated.go -type ValueSpec,FuncDecl,TypeSpec

type ValueSpec struct {
	pkg *packages.Package
	*ast.ValueSpec
}

type FuncDecl struct {
	pkg *packages.Package
	*ast.FuncDecl
}

type TypeSpec struct {
	pkg *packages.Package
	*ast.TypeSpec
}

func (v *ValueSpec) Pkg() *packages.Package { return v.pkg }
func (f *FuncDecl) Pkg() *packages.Package  { return f.pkg }
func (t *TypeSpec) Pkg() *packages.Package  { return t.pkg }

type Lister interface {
	List(pkgs ...*packages.Package) []Node
}

func New() Lister {
	return &lister{}
}

type lister struct{}

func (l *lister) List(pkgs ...*packages.Package) []Node {
	var (
		result = []Node{}
		add    = func(n Node) {
			result = append(result, n)
		}
		debug = func(kind string, name string, node Node) {
			attrs := []logx.Attr{
				logx.S("kind", kind),
				logx.S("name", name),
				logx.S("pkg", node.Pkg().PkgPath),
			}
			attrs = append(
				attrs,
				logx.Any("begin", astutil.PosString(node.Pos(), node.Pkg().Fset)),
				logx.Any("end", astutil.PosString(node.End(), node.Pkg().Fset)),
			)
			logx.Debug("def loaded", attrs...)
		}
	)

	defer func() {
		logx.Info("def loaded", logx.I("len", len(result)))
	}()

	for _, pkg := range pkgs {
		logx.Debug("def load", logx.S("package", pkg.PkgPath))
		for _, f := range pkg.Syntax {
			for _, decl := range f.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range decl.Specs {
						switch spec := spec.(type) {
						case *ast.TypeSpec:
							n := &TypeSpec{
								pkg:      pkg,
								TypeSpec: spec,
							}
							debug("TypeSpec", spec.Name.String(), n)
							add(n)
						case *ast.ValueSpec:
							n := &ValueSpec{
								pkg:       pkg,
								ValueSpec: spec,
							}
							debug("ValueSpec", spec.Names[0].Name, n)
							add(n)
						}
					}
				case *ast.FuncDecl:
					n := &FuncDecl{
						pkg:      pkg,
						FuncDecl: decl,
					}
					debug("FuncDecl", decl.Name.String(), n)
					add(n)
				}
			}
		}
	}
	return result
}
