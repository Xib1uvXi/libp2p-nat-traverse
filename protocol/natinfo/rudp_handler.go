package natinfo

import (
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/quic-go/quic-go"
	"net"
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
		p.ErrC <- err
		return err
	}
	p.RemoteAddr = conn.RemoteAddr()

	p.DoneC <- struct{}{}
	return nil
}

// QUICReceiveMessage quic read message
func receiveMessage(stream quic.Stream) (*pb.Message, error) {
	r := pbio.NewDelimitedReader(stream, maxMsgSize)
	var msg pb.Message

	if err := r.ReadMsg(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
