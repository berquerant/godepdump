package testx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/packages"
)

type ParsedFile struct {
	File    *ast.File
	FileSet *token.FileSet
}

func (f *ParsedFile) Print() {
	ast.Print(f.FileSet, f.File)
}

func ParseFile(t *testing.T, src string) *ParsedFile {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return &ParsedFile{
		File:    f,
		FileSet: fset,
	}
}

func ParseAsPackage(t *testing.T, src string) *packages.Package {
	t.Helper()
	f := ParseFile(t, src)
	return &packages.Package{
		Fset:   f.FileSet,
		Syntax: []*ast.File{f.File},
	}
}
