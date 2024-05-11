package natinfo

import (
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"github.com/libp2p/go-msgio/pbio"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"math/rand"
	"net"
)

func newRandomPort() int {
	randPort := rand.Intn(20000) + 10000
	return randPort
}

func multiAddrToUDPAddr(addr ma.Multiaddr) (net.Addr, error) {
	ipv4, err := addr.ValueForProtocol(ma.P_IP4)
	if err != nil {
		return nil, err
	}

	udpport, err := addr.ValueForProtocol(ma.P_UDP)
	if err != nil {
		return nil, err
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", ipv4+":"+udpport)
	if err != nil {
		return nil, err
	}

	return udpAddr, nil
}

func receiveMessage(reader io.Reader) (*pb.Message, error) {
	r := pbio.NewDelimitedReader(reader, maxMsgSize)
	var msg pb.Message

	if err := r.ReadMsg(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
