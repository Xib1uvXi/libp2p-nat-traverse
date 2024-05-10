package natinfo

import (
	"fmt"
	"github.com/Xib1uvXi/libp2p-nat-traverse/netutil"
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-msgio/pbio"
	"google.golang.org/protobuf/proto"
	"math"
	"net"
	"time"
)

var logger = log.Logger("nat-traverse/protocol/natinfo")

var streamHandlerTimeout = 60 * time.Second

const maxMsgSize = 4 * 1024

type NATTypeObserServer struct {
	PublicIP string
	host     host.Host
}

func (t *NATTypeObserServer) handleNewStream(s network.Stream) {
	defer s.Close()
	_ = s.SetDeadline(time.Now().Add(streamHandlerTimeout))

	remotePeerID := s.Conn().RemotePeer()
	logger.Debugf("new stream from %s", remotePeerID)

	reader := pbio.NewDelimitedReader(s, maxMsgSize)
	writer := pbio.NewDelimitedWriter(s)

	// port negotiation
	newAddr, err := t.portNegotiation(reader, writer)
	if err != nil {
		logger.Errorf("failed to handle test req: %v", err)
		s.Reset()
		return
	}

	oldAddr, err := multiAddrToUDPAddr(s.Conn().RemoteMultiaddr())
	if err != nil {
		logger.Errorf("failed to convert multiaddr to udp addr: %v", err)
		s.Reset()
		return
	}

	if newAddr.String() != oldAddr.String() {
		// nat type is symmetric
		if err := t.symmetricNatType(oldAddr, newAddr, writer); err != nil {
			logger.Errorf("failed to handle symmetric nat type: %v", err)
			s.Reset()
			return
		}

		return
	}
}

// handle req msg, response with PortNegotiation port
func (t *NATTypeObserServer) portNegotiation(reader pbio.ReadCloser, writer pbio.WriteCloser) (net.Addr, error) {
	var msg pb.Message

	if err := reader.ReadMsg(&msg); err != nil {
		logger.Errorf("failed to read test req msg: %v", err)
		return nil, err
	}

	// check msg type is TestNatType
	if msg.GetType() != pb.MsgType_TestNatType {
		logger.Errorf("expect test req msg type, got %s", msg.GetType())
		return nil, fmt.Errorf("expect test req msg type, got %s", msg.GetType())
	}

	// random listen a port, wait for client to connect, check client public addr
	randomPort := newRandomPort()

	// new reliable server
	tmpRUDP, err := netutil.NewReliableUDPServer(fmt.Sprintf("%s:%d", t.PublicIP, randomPort))
	if err != nil {
		logger.Errorf("failed to create reliable udp server: %v", err)
		return nil, err
	}

	defer tmpRUDP.Close()

	recorder := NewPortNegotiationRecorder()
	tmpRUDP.SetStreamHandler(recorder.PortNegotiationHandler)
	if err := tmpRUDP.Start(); err != nil {
		logger.Errorf("failed to start reliable udp server: %v", err)
		return nil, err
	}

	// send response msg
	resMsg := &pb.Message{
		Type: pb.MsgType_PortNegotiation,
		Data: []byte(fmt.Sprintf("%d", randomPort)),
	}

	if err := writer.WriteMsg(resMsg); err != nil {
		logger.Errorf("failed to write port negotiation msg: %v", err)
		return nil, err
	}

	// block until client connect, timeout 5s
	select {
	case err := <-recorder.ErrC:
		logger.Errorf("port negotiation error: %v", err)
		return nil, err

	case <-time.After(5 * time.Second):
		logger.Errorf("get port negotiation resptimeout")
		return nil, fmt.Errorf("get port negotiation resp timeout")

	case <-recorder.DoneC:
		logger.Debugf("port negotiation done")
	}

	return recorder.RemoteAddr, nil
}

// handle symmetric nat type
func (t *NATTypeObserServer) symmetricNatType(oldAddr, newAddr net.Addr, writer pbio.WriteCloser) error {
	var changeRule pb.PortChangeType
	// observe port changes
	// changes range
	portRange := math.Abs(float64(oldAddr.(*net.UDPAddr).Port - newAddr.(*net.UDPAddr).Port))
	if portRange <= 100 {
		changeRule = pb.PortChangeType_Linear
	} else {
		changeRule = pb.PortChangeType_Random
	}

	natInfo := &pb.NATTypeInfo{
		NatType:           pb.NATType_Symmetric,
		UdpPortChangeRule: changeRule,
	}

	data, err := proto.Marshal(natInfo)
	if err != nil {
		logger.Errorf("failed to marshal nat info: %v", err)
		return err
	}

	// send result
	resMsg := &pb.Message{
		Type: pb.MsgType_NatTypeResult,
		Data: data,
	}

	if err := writer.WriteMsg(resMsg); err != nil {
		logger.Errorf("failed to write nat type result msg: %v", err)
		return err
	}

	return nil
}
