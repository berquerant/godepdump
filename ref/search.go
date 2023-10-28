package ref

import (
	"fmt"
	"go/token"
	"log/slog"
	"strings"

	"github.com/berquerant/godepdump/astutil"
	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/logx"
	"golang.org/x/tools/go/packages"
)

type Searcher interface {
	// Search searches the definition to which pos belongs.
	Search(pkg *packages.Package, pos token.Pos) (def.Node, bool)
}

func NewSearcher(nodeList []def.Node) Searcher {
	return &searcher{
		d: newPkgNodeListMap(nodeList),
	}
}

type searcher struct {
	d pkgNodeListMap
}

func (s *searcher) Search(pkg *packages.Package, pos token.Pos) (result def.Node, found bool) {
	defer func() {
		attrs := []logx.Attr{
			logx.S("pkg", pkg.PkgPath),
			logx.B("found", found),
			logx.Any("pos", astutil.PosString(pos, pkg.Fset)),
		}
		if result != nil {
			attrs = append(
				attrs,
				logx.S("result", fmt.Sprintf("%T", result)),
				logx.Any("begin", astutil.PosString(result.Pos(), pkg.Fset)),
				logx.Any("end", astutil.PosString(result.End(), pkg.Fset)),
			)
		}
		logx.Debug("ref searcher search", attrs...)
	}()

	nodeList, ok := s.d[pkg.ID]
	if !ok {
		return
	}

	logx.Debug("ref searcher search pkg", logx.S("pkg", pkg.PkgPath), logx.I("nodes", len(nodeList)))

	for _, p := range nodeList {
		if p.in(pos) {
			result = p.node
			found = true
			return
		}
	}

	return
}

type pkgNodeListMap map[string][]*posRange

func newPkgNodeListMap(nodeList []def.Node) pkgNodeListMap {
	result := map[string][]*posRange{}
	for _, x := range nodeList {
		id := x.Pkg().ID
		if p, ok := newPosRange(x); ok {
			result[id] = append(result[id], p)
			logx.Debug("ref searcher loaded", logx.Any("pos", p))
		}
	}
	var count int
	for _, x := range result {
		logx.Debug("ref searcher loaded", logx.S("pkg", x[0].node.Pkg().PkgPath), logx.I("len", len(x)))
		count += len(x)
	}
	logx.Info("ref searcher loaded", logx.I("given", len(nodeList)), logx.I("accepted", count))
	return result
}

func newPosRange(node def.Node) (*posRange, bool) {
	r := &posRange{
		node: node,
	}
	switch node := node.(type) {
	case *def.FuncDecl:
		r.kind = "FuncDecl"
		r.name = node.Name.String()
		if node.Type.TypeParams != nil {
			r.begin = node.Type.TypeParams.Opening // generic parameter beginning
		} else {
			r.begin = node.Type.Params.Opening // argument beginning
		}
		r.end = node.Body.Rbrace
	case *def.TypeSpec:
		r.kind = "TypeSpec"
		r.name = node.Name.String()
		r.begin = node.Type.Pos()
		r.end = node.Type.End()
	case *def.ValueSpec:
		r.kind = "ValueSpec"
		// begin
		switch {
		case node.Type != nil:
			r.begin = node.Type.Pos()
		case len(node.Values) > 0:
			r.begin = node.Values[0].Pos()
		default:
			r.begin = node.Names[len(node.Names)-1].End()
		}
		// end
		switch {
		case len(node.Values) > 0:
			r.end = node.Values[len(node.Values)-1].End()
		case node.Type != nil:
			r.end = node.Type.End()
		default:
			r.end = node.Names[len(node.Names)-1].End()
		}
		// name
		names := make([]string, len(node.Names))
		for i, name := range node.Names {
			names[i] = name.String()
		}
		r.name = strings.Join(names, ",")
	default:
		return nil, false
	}

	r.beginString = astutil.PosString(r.begin, node.Pkg().Fset)
	r.endString = astutil.PosString(r.end, node.Pkg().Fset)
	return r, true
}

type posRange struct {
	node        def.Node
	kind        string // for log
	name        string // for log
	begin       token.Pos
	beginString string // for log
	end         token.Pos
	endString   string // for log
}

func (r *posRange) in(pos token.Pos) bool { return r.begin <= pos && pos <= r.end }

func (r *posRange) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("pkg", r.node.Pkg().PkgPath),
		slog.String("kind", r.kind),
		slog.String("name", r.name),
		slog.String("begin", r.beginString),
		slog.String("end", r.endString),
	)
}
