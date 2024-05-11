package natinfo

import (
	"fmt"
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"github.com/quic-go/quic-go"
	"net"
	"time"
)

type PortNegotiationRecorder struct {
	ErrC       chan error
	RemoteAddr net.Addr
	DoneC      chan struct{}
}

func NewPortNegotiationRecorder() *PortNegotiationRecorder {
	return &PortNegotiationRecorder{
		ErrC:  make(chan error),
		DoneC: make(chan struct{}),
	}
}

func (p *PortNegotiationRecorder) PortNegotiationHandler(conn quic.Connection, stream quic.Stream) error {
	msg, err := receiveMessage(stream)
	if err != nil {
		p.ErrC <- err
		return err
	}

	if msg.GetType() != pb.MsgType_PortNegotiationResponse {
		p.ErrC <- fmt.Errorf("PortNegotiationRecorder unexpected message type: %v", msg.GetType())
		return err
	}
	p.RemoteAddr = conn.RemoteAddr()

	p.DoneC <- struct{}{}
	return nil
}

func (p *PortNegotiationRecorder) Await() (net.Addr, error) {
	select {
	case <-p.DoneC:
		return p.RemoteAddr, nil
	case err := <-p.ErrC:
		return nil, err
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("PortNegotiationRecorder timeout")
	}
}
