package def_test

import (
	"testing"

	"github.com/berquerant/godepdump/def"
	"github.com/berquerant/godepdump/testx"
	"github.com/stretchr/testify/assert"
)

func TestLister(t *testing.T) {
	for _, tc := range []listerTestcase{
		{
			name: "a func",
			src: `package p
func main() {
  println("hello")
}`,
			want: newAssertNodeList(
				newAssertFuncDecl("main"),
			),
		},
		{
			name: "a type",
			src: `package p
type Yen int`,
			want: newAssertNodeList(
				newAssertTypeSpec("Yen"),
			),
		},
		{
			name: "a var",
			src: `package p
var Global = "north"`,
			want: newAssertNodeList(
				newAssertValueSpec("Global"),
			),
		},
		{
			name: "funcs",
			src: `package p
func F1(_ string) bool { return true }
func F2(_ int) string { return "true" }`,
			want: newAssertNodeList(
				newAssertFuncDecl("F1"),
				newAssertFuncDecl("F2"),
			),
		},
		{
			name: "types",
			src: `package p
type planet string
type Secondary = uintptr
`,
			want: newAssertNodeList(
				newAssertTypeSpec("planet"),
				newAssertTypeSpec("Secondary"),
			),
		},
		{
			name: "vars",
			src: `package p
var V1 = 1
var V2, V3 = 2, 3
const C1 = 0`,
			want: newAssertNodeList(
				newAssertValueSpec("V1"),
				newAssertValueSpec("V2", "V3"),
				newAssertValueSpec("C1"),
			),
		},
		{
			name: "a method",
			src: `package p
type Empty struct{}
func (*Empty) String() string {return "null"}`,
			want: newAssertNodeList(
				newAssertTypeSpec("Empty"),
				newAssertFuncDecl("String"),
			),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pkg := testx.ParseAsPackage(t, tc.src)
			got := def.New().List(pkg)
			if !assert.Equal(t, len(tc.want), len(got)) {
				return
			}
			for i, w := range tc.want {
				w.Assert(t, got[i])
			}
		})
	}
}

type listerTestcase struct {
	name string
	src  string
	want []assertNode
}

func newAssertNodeList(node ...assertNode) []assertNode {
	return node
}

type assertNode interface {
	Assert(t *testing.T, node def.Node)
}

func newAssertValueSpec(names ...string) *assertValueSpec {
	return &assertValueSpec{
		names: names,
	}
}

type assertValueSpec struct {
	names []string
}

func (a *assertValueSpec) Assert(t *testing.T, node def.Node) {
	v, ok := node.(*def.ValueSpec)
	if !assert.True(t, ok) {
		return
	}
	names := make([]string, len(v.Names))
	for i, n := range v.Names {
		names[i] = n.String()
	}
	assert.Equal(t, a.names, names)
}

func newAssertTypeSpec(name string) *assertTypeSpec {
	return &assertTypeSpec{
		name: name,
	}
}

type assertTypeSpec struct {
	name string
}

func (a *assertTypeSpec) Assert(t *testing.T, node def.Node) {
	v, ok := node.(*def.TypeSpec)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, a.name, v.Name.String())
}

func newAssertFuncDecl(name string) *assertFuncDecl {
	return &assertFuncDecl{
		name: name,
	}
}

type assertFuncDecl struct {
	name string
}

func (a *assertFuncDecl) Assert(t *testing.T, node def.Node) {
	v, ok := node.(*def.FuncDecl)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, a.name, v.Name.String())
}
