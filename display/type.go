package display

import (
	"go/types"

	"github.com/berquerant/godepdump/logx"
	"github.com/berquerant/godepdump/typesx"
)

type Type struct {
	Type typesx.TypeType `json:"type"`
	Map  Map             `json:"map"`
}

type Map map[string]any

func NewMap() Map { return Map(map[string]any{}) }

func (m Map) Set(k string, v any) Map {
	logx.Debug("analyze map", logx.S("key", k), logx.V("value", v))
	if _, exist := m[k]; exist {
		logx.Warn("analyze map key duplicated", logx.S("key", k), logx.V("value", v))
	}
	m[k] = v
	return m
}

func (t *Type) Set(k string, v any) *Type {
	t.Map.Set(k, v)
	return t
}

//go:generate go run github.com/berquerant/goconfig@v0.3.0 -field "Limit int" -option -output type_config_generated.go -prefix TypeAnalyze

func NewTypeAnalyzer(opt ...TypeAnalyzeConfigOption) TypeAnalyzer {
	config := NewTypeAnalyzeConfigBuilder().Limit(1).Build()
	config.Apply(opt...)
	return &typeAnalyzerImpl{
		config: config,
	}
}

type TypeAnalyzer interface {
	Analyze(t types.Type) *Type
}

type typeAnalyzerImpl struct {
	config *TypeAnalyzeConfig
}

func (a *typeAnalyzerImpl) Analyze(t types.Type) *Type {
	logx.Debug("analyze start", logx.S("string", t.String()))
	analyzer := &typeAnalyzer{
		limit:    a.config.Limit.Get(),
		initType: t,
	}
	return analyzer.analyze(t, 0)
}

type typeAnalyzer struct {
	limit    int
	initType types.Type
}

func (a *typeAnalyzer) analyze(t types.Type, count int) *Type {
	r := &Type{
		Type: typesx.GetTypeType(t),
		Map:  NewMap(),
	}

	str := t.String()
	logx.Debug("analyze",
		logx.S("string", str),
		logx.S("type", typesx.GetTypeType(t).String()),
		logx.V("v", t), logx.I("count", count),
		logx.S("init", a.initType.String()),
	)
	r.Set("string", str)

	// Limit the number of recursive analyze()
	if count >= a.limit {
		return r
	}
	count++

	switch t := t.(type) {
	case *types.Array:
		return a.analyzeArray(r, t, count)
	case *types.Basic:
		return a.analyzeBasic(r, t, count)
	case *types.Chan:
		return a.analyzeChan(r, t, count)
	case *types.Interface:
		return a.analyzeInterface(r, t, count)
	case *types.Map:
		return a.analyzeMap(r, t, count)
	case *types.Named:
		return a.analyzeNamed(r, t, count)
	case *types.Pointer:
		return a.analyzePointer(r, t, count)
	case *types.Signature:
		return a.analyzeSignature(r, t, count)
	case *types.Slice:
		return a.analyzeSlice(r, t, count)
	case *types.Struct:
		return a.analyzeStruct(r, t, count)
	case *types.Tuple:
		return a.analyzeTuple(r, t, count)
	case *types.TypeParam:
		return a.analyzeTypeParam(r, t, count)
	case *types.Union:
		return a.analyzeUnion(r, t, count)
	default:
		return nil
	}
}

/*
 * analyzeXXX() should be called via analyze().
 */

func (a *typeAnalyzer) analyzeArray(v *Type, t *types.Array, c int) *Type {
	return v.Set("elem", a.analyze(t.Elem(), c))
}

func (a *typeAnalyzer) analyzeBasic(v *Type, t *types.Basic, c int) *Type {
	return v.Set("kind", typesx.BasicKindToString(t.Kind())).
		Set("name", t.Name())
}

func (a *typeAnalyzer) analyzeChan(v *Type, t *types.Chan, c int) *Type {
	return v.Set("dir", typesx.ChanDirToString(t.Dir())).
		Set("elem", a.analyze(t.Elem(), c))
}

func (a *typeAnalyzer) analyzeInterface(v *Type, t *types.Interface, c int) *Type {
	if n := t.NumMethods(); n > 0 {
		ms := make([]any, n)
		for i := 0; i < n; i++ {
			ms[i] = a.breakDownFunc(t.Method(i), c)
		}
		v.Set("methods", ms)
	}
	if n := t.NumEmbeddeds(); n > 0 {
		es := make([]*Type, n)
		for i := 0; i < n; i++ {
			es[i] = a.analyze(t.EmbeddedType(i), c)
		}
		v.Set("embeddeds", es)
	}
	return v
}

func (a *typeAnalyzer) analyzeMap(v *Type, t *types.Map, c int) *Type {
	return v.Set("elem", a.analyze(t.Elem(), c)).
		Set("key", a.analyze(t.Key(), c))
}

func (a *typeAnalyzer) analyzeNamed(v *Type, t *types.Named, c int) *Type {
	if n := t.NumMethods(); n > 0 {
		ms := make([]any, n)
		for i := 0; i < n; i++ {
			ms[i] = a.breakDownFunc(t.Method(i), c)
		}
		v.Set("method", ms)
	}
	v.Set("obj", a.breakDownTypeName(t.Obj(), c))
	v.Set("type", a.analyze(t.Underlying(), c))
	v.Set("args", a.breakDownTypeList(t.TypeArgs(), c))

	if p := t.TypeParams(); p != nil {
		v.Set("type_params", a.breakDownTypeParamList(p, c))
	}
	return v
}

func (a *typeAnalyzer) analyzePointer(v *Type, t *types.Pointer, c int) *Type {
	return v.Set("elem", a.analyze(t.Elem(), c))
}

func (a *typeAnalyzer) analyzeSignature(v *Type, t *types.Signature, c int) *Type {
	if p := t.Params(); p != nil {
		v.Set("params", a.analyze(p, c))
	}
	if r := t.Recv(); r != nil {
		v.Set("recv", a.breakDownVar(r, c))
	}
	if s := t.RecvTypeParams(); s != nil {
		v.Set("recv_type_params", a.breakDownTypeParamList(s, c))
	}
	if r := t.Results(); r != nil {
		v.Set("result", a.analyze(r, c))
	}
	if p := t.TypeParams(); p != nil {
		v.Set("type_params", a.breakDownTypeParamList(p, c))
	}
	v.Set("variadic", t.Variadic())
	return v
}

func (a *typeAnalyzer) analyzeSlice(v *Type, t *types.Slice, c int) *Type {
	return v.Set("elem", a.analyze(t.Elem(), c))
}

func (a *typeAnalyzer) analyzeStruct(v *Type, t *types.Struct, c int) *Type {
	if n := t.NumFields(); n > 0 {
		fields := make([]Map, n)
		for i := 0; i < n; i++ {
			fields[i] = a.breakDownVar(t.Field(i), c)
		}
		v.Set("fields", fields)
	}
	return v
}

func (a *typeAnalyzer) analyzeTuple(v *Type, t *types.Tuple, c int) *Type {
	ts := make([]Map, t.Len())
	for i := 0; i < t.Len(); i++ {
		ts[i] = a.breakDownVar(t.At(i), c)
	}
	v.Set("vars", ts)
	return v
}

func (a *typeAnalyzer) analyzeTypeParam(v *Type, t *types.TypeParam, c int) *Type {
	return v.Set("constraint", a.analyze(t.Constraint(), c)).
		Set("index", t.Index()).
		Set("obj", a.breakDownTypeName(t.Obj(), c)).
		Set("type", a.analyze(t.Underlying(), c))
}

func (a *typeAnalyzer) analyzeUnion(v *Type, t *types.Union, c int) *Type {
	terms := make([]Map, t.Len())
	for i := 0; i < t.Len(); i++ {
		terms[i] = a.breakDownTerm(t.Term(i), c)
	}
	v.Set("terms", terms)
	return v
}

/*
 * breakDownXXX() should be called via analyzeXXX().
 */

func (a *typeAnalyzer) breakDownTerm(t *types.Term, c int) Map {
	logx.Debug("break down term", logx.S("t", t.String()), logx.V("v", t))
	return NewMap().
		Set("tilde", t.Tilde()).
		Set("type", a.analyze(t.Type(), c))
}

func (a *typeAnalyzer) breakDownTypeList(t *types.TypeList, c int) []*Type {
	n := t.Len()
	args := make([]*Type, n)
	for i := 0; i < n; i++ {
		logx.Debug("break down type list", logx.I("index", i), logx.S("t", t.At(i).String()), logx.V("v", t))
		args[i] = a.analyze(t.At(i), c)
	}
	return args
}

func (a *typeAnalyzer) breakDownTypeParamList(t *types.TypeParamList, c int) []*Type {
	n := t.Len()
	args := make([]*Type, n)
	for i := 0; i < n; i++ {
		logx.Debug("break down type param list", logx.I("index", i), logx.S("t", t.At(i).String()), logx.V("v", t))
		args[i] = a.analyze(t.At(i), c)
	}
	return args
}

func (a *typeAnalyzer) breakDownFunc(t *types.Func, c int) Map {
	logx.Debug("break down func", logx.S("t", t.String()), logx.V("v", t))
	return NewMap().
		Set("fullname", t.FullName()).
		Set("name", t.Name()).
		Set("exported", t.Exported())
}

func (a *typeAnalyzer) breakDownTypeName(t *types.TypeName, c int) Map {
	logx.Debug("break down type name", logx.S("t", t.String()), logx.V("v", t))
	return NewMap().
		Set("id", t.Id()).
		Set("name", t.Name()).
		Set("alias", t.IsAlias()).
		Set("exported", t.Exported())
}

func (a *typeAnalyzer) breakDownVar(t *types.Var, c int) Map {
	logx.Debug("break down var", logx.S("t", t.String()), logx.V("v", t))
	return NewMap().
		Set("id", t.Id()).
		Set("name", t.Name()).
		Set("field", t.IsField()).
		Set("embedded", t.Embedded()).
		Set("exported", t.Exported()).
		Set("type", a.analyze(t.Type(), c))
}
