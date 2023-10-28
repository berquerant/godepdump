package astutil

import (
	"go/ast"
	"go/token"
)

// InPos returns true if pos is in node.
func InPos(node ast.Node, pos token.Pos) bool { return node.Pos() <= pos && pos <= node.End() }

// PosString converts pos to human-readable string.
func PosString(pos token.Pos, fset *token.FileSet) string { return fset.Position(pos).String() }
