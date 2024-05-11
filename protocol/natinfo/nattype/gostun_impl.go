package nattype

import "github.com/ccding/go-stun/stun"

type GoStunImpl struct {
}

func (g *GoStunImpl) GetNATType() (NATType, error) {
	clinet := stun.NewClient()
	clinet.SetServerHost("stun.syncthing.net", 3478)
	nat, _, err := clinet.Discover()
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
