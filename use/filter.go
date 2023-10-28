package use

func UniverseFilter(n Node) bool { return n.Obj() != nil && n.Obj().Pkg() == nil }
func ExportedFilter(n Node) bool { return n.Obj() != nil && n.Obj().Exported() }
