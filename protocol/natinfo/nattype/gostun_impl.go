package nattype

import "github.com/ccding/go-stun/stun"

type GoStunImpl struct {
}

func (g *GoStunImpl) GetNATType() (NATType, error) {
	nat, _, err := stun.NewClient().Discover()
	if err != nil {
		return UnKnown, err
	}

	switch nat {
	case stun.NATNone:
		return None, nil
	case stun.NATFull:
		return FullCone, nil
	case stun.NATRestricted:
		return RestrictedCone, nil
	case stun.NATPortRestricted:
		return PortRestrictedCone, nil
	case stun.NATSymetric:
		return Symmetric, nil
	default:
		return UnKnown, nil
	}
}
