package natinfo

import (
	"fmt"
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"google.golang.org/protobuf/proto"
	"net"
)

type PortRestrictedChecker struct {
	ErrC  chan error
	DoneC chan struct{}
}

func NewPortRestrictedChecker() *PortRestrictedChecker {
	return &PortRestrictedChecker{
		ErrC:  make(chan error),
		DoneC: make(chan struct{}),
	}
}

func (p *PortRestrictedChecker) PortRestrictedHandler(conn net.Conn) error {
	var data = make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		p.ErrC <- err
		return err
	}
	msg := &pb.Message{}

	err = proto.Unmarshal(data[:n], msg)
	if err != nil {
		p.ErrC <- err
		return err
	}

	if msg.GetType() != pb.MsgType_ServerPortChangeTestResponse {
		p.ErrC <- fmt.Errorf("PortRestrictedChecker unexpected message type: %v", msg.GetType())
		return err
	}

	p.DoneC <- struct{}{}
	return nil
}
