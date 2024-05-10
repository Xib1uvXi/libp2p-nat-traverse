package nattype

const (
	UnKnown NATType = iota
	None
	FullCone
	RestrictedCone
	PortRestrictedCone
	Symmetric
)

type NATType int

func (n NATType) String() string {
	switch n {
	case FullCone:
		return "FullCone"
	case RestrictedCone:
		return "RestrictedCone"
	case PortRestrictedCone:
		return "PortRestrictedCone"
	case Symmetric:
		return "Symmetric"
	case None:
		return "None"
	default:
		return "UnKnown"
	}
}
