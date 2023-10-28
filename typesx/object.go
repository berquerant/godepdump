package typesx

import (
	"encoding/json"
	"go/types"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=ObjectType -output object_stringer_generated.go

type ObjectType int

func (t ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

const (
	Ounknown ObjectType = iota
	Obuiltin
	Oconst
	Ofunc
	Omethod
	Olabel
	Onil
	OpkgName
	OtypeName
	Ofield
	Ovar
)

func GetObjectType(obj types.Object) ObjectType {
	switch obj := obj.(type) {
	case *types.Builtin:
		return Obuiltin
	case *types.Const:
		return Oconst
	case *types.Func:
		if obj.Type().(*types.Signature).Recv() != nil {
			return Omethod
		}
		return Ofunc
	case *types.Label:
		return Olabel
	case *types.Nil:
		return Onil
	case *types.PkgName:
		return OpkgName
	case *types.TypeName:
		return OtypeName
	case *types.Var:
		if obj.IsField() {
			return Ofield
		}
		return Ovar
	}

	return Ounknown
}
