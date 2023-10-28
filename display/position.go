package display

import (
	"go/token"
	"log/slog"
)

type Position struct {
	Pos    int    `json:"pos"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"col"`
	Offset int    `json:"offset"`
	String string `json:"string"`
}

func NewPosition(pos token.Pos, fset *token.FileSet) *Position {
	p := fset.Position(pos)
	return &Position{
		Pos:    int(pos),
		File:   p.Filename,
		Line:   p.Line,
		Column: p.Column,
		Offset: p.Offset,
		String: p.String(),
	}
}

func (p *Position) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("pos", p.Pos),
		slog.String("file", p.File),
		slog.Int("line", p.Line),
		slog.Int("column", p.Column),
		slog.Int("offset", p.Offset),
		slog.String("string", p.String),
	)
}
