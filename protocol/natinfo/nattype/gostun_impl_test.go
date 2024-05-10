package nattype

import (
	"github.com/ccding/go-stun/stun"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoStunImpl_GetNATType(t *testing.T) {
	t.Skip("Skip this test because it will fail in CI/CD environment")
	impl := &GoStunImpl{}

	result, err := impl.GetNATType()
	assert.NoError(t, err)

	assert.Equal(t, Symmetric, result)
}

func TestGoStunImpl_GetNATType2(t *testing.T) {
	t.Skip("Skip this test because it will fail in CI/CD environment")
	clinet := stun.NewClient()
	//clinet.SetServerHost("stun.miwifi.com", 3478)
	//clinet.SetServerHost("stun.miwifi.com", 3478)
	clinet.SetServerHost("stun.syncthing.net", 3478)
	nat, ip, err := clinet.Discover()
	assert.NoError(t, err)

	t.Logf("NAT: %v, IP: %v", nat, ip)

}
