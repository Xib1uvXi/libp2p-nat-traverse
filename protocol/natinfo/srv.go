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

	var natInfo *pb.NATTypeInfo

	if newAddr.String() != oldAddr.String() {
		// nat type is symmetric
		natInfo, err = t.symmetricNatType(oldAddr, newAddr)
		if err != nil {
			logger.Errorf("failed to check symmetric nat type: %v", err)
			s.Reset()
			return
		}
	} else {
		// nat type is not symmetric
		natInfo, err = t.coneCheck(oldAddr)
		if err != nil {
			logger.Errorf("failed to check cone nat type: %v", err)
			s.Reset()
			return
		}
	}

	if err := t.sendResult(natInfo, writer); err != nil {
		logger.Errorf("failed to send nat type result: %v", err)
		s.Reset()
		return
	}

	logger.Debugf("send nat type result: %s, PeerID: %s", natInfo.NatType, remotePeerID.String())
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
	remoteAddr, err := recorder.Await()
	if err != nil {
		logger.Errorf("failed to await port negotiation: %v", err)
		return nil, err
	}

	return remoteAddr, nil
}

// handle symmetric nat type
func (t *NATTypeObserServer) symmetricNatType(oldAddr, newAddr net.Addr) (*pb.NATTypeInfo, error) {
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

	return natInfo, nil
}

// check port restricted nat
func (t *NATTypeObserServer) coneCheck(raddr net.Addr) (*pb.NATTypeInfo, error) {
	// random listen a port
	randomPort := newRandomPort()
	wrapUDP, err := netutil.NewWrapUDPConn(fmt.Sprintf("%s:%d", t.PublicIP, randomPort))
	if err != nil {
		logger.Errorf("failed to create wrap udp server: %v", err)
		return nil, err
	}

	// set handler
	checker := NewPortRestrictedChecker()
	wrapUDP.SetConnHandler(checker.PortRestrictedHandler)

	if err := wrapUDP.Start(); err != nil {
		logger.Errorf("failed to start wrap udp server: %v", err)
		return nil, err
	}
	defer wrapUDP.Close()

	// send Unreliable udp Message to client
	data, err := proto.Marshal(&pb.Message{Type: pb.MsgType_ServerPortChangeTest})
	if err != nil {
		logger.Errorf("failed to marshal port change test msg: %v", err)
		return nil, err
	}

	// send unreliable udp msg 2 times
	_ = wrapUDP.SendUnreliableUDPMessage(data, raddr)
	_ = wrapUDP.SendUnreliableUDPMessage(data, raddr)

	// block until client resp, timeout 2s
	select {
	case err := <-checker.ErrC:
		logger.Errorf("port restricted check error: %v", err)
		return nil, err
	case <-time.After(2 * time.Second):
		// port restricted
		return &pb.NATTypeInfo{
			NatType: pb.NATType_PortRestrictedCone,
		}, nil
	case <-checker.DoneC:
		// full or restricted
		return &pb.NATTypeInfo{
			NatType: pb.NATType_FullOrRestrictedCone,
		}, nil
	}
}

func (t *NATTypeObserServer) sendResult(natInfo *pb.NATTypeInfo, writer pbio.WriteCloser) error {
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
