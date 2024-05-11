package natinfo

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_multiAddrToUDPAddr(t *testing.T) {
	maAddrStr := "/ip4/127.0.0.1/udp/1234/quic-v1/p2p/12D3KooWRv6jAyvZc6nTZ14ixHppwuN4VM52S3UnDmwzHCTR93Jx"

	maAddr, err := ma.NewMultiaddr(maAddrStr)
	assert.NoError(t, err)

	udpAddr, err := multiAddrToUDPAddr(maAddr)

	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:1234", udpAddr.String())
}
