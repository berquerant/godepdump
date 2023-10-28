package typesx

import (
	"encoding/json"
	"go/types"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=TypeType -output type_stringer_generated.go

type TypeType int

func (t TypeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

const (
	Tunknown TypeType = iota
	Tarray
	Tbasic
	Tchan
	Tinterface
	Tmap
	Tnamed
	Tpointer
	Tsignature
	Tslice
	Tstruct
	Ttuple
	TtypeParam
	Tunion
)

func GetTypeType(typ types.Type) TypeType {
	switch typ.(type) {
	case *types.Array:
		return Tarray
	case *types.Basic:
		return Tbasic
	case *types.Chan:
		return Tchan
	case *types.Interface:
		return Tinterface
	case *types.Map:
		return Tmap
	case *types.Named:
		return Tnamed
	case *types.Pointer:
		return Tpointer
	case *types.Signature:
		return Tsignature
	case *types.Slice:
		return Tslice
	case *types.Struct:
		return Tstruct
	case *types.Tuple:
		return Ttuple
	case *types.TypeParam:
		return TtypeParam
	case *types.Union:
		return Tunion
	}

	return Tunknown
}

func BasicKindToString(k types.BasicKind) string {
	switch k {
	case types.Bool:
		return "Bool"
	case types.Int:
		return "Int"
	case types.Int8:
		return "Int8"
	case types.Int16:
		return "Int16"
	case types.Int32:
		return "Int32"
	case types.Int64:
		return "Int64"
	case types.Uint:
		return "Uint"
	case types.Uint8:
		return "Uint8"
	case types.Uint16:
		return "Uint16"
	case types.Uint32:
		return "Uint32"
	case types.Uint64:
		return "Uint64"
	case types.Uintptr:
		return "Uintptr"
	case types.Float32:
		return "Float32"
	case types.Float64:
		return "Float64"
	case types.Complex64:
		return "Complex64"
	case types.Complex128:
		return "Complex128"
	case types.String:
		return "String"
	case types.UnsafePointer:
		return "UnsafePointer"
	case types.UntypedBool:
		return "UntypedBool"
	case types.UntypedInt:
		return "UntypedInt"
	case types.UntypedRune:
		return "UntypedRune"
	case types.UntypedFloat:
		return "UntypedFloat"
	case types.UntypedComplex:
		return "UntypedComplex"
	case types.UntypedString:
		return "UntypedString"
	case types.UntypedNil:
		return "UntypedNil"
	default:
		return "Invalid"
	}
}

func ChanDirToString(d types.ChanDir) string {
	switch d {
	case types.SendOnly:
		return "Send"
	case types.RecvOnly:
		return "Recv"
	default:
		return "SendRecv"
	}
}
