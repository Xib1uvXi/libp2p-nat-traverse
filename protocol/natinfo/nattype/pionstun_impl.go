package nattype

import (
	"errors"
	"github.com/pion/stun/v2"
	"net"
	"time"
)

const (
	stunServer = "stun.miwifi.com:3478"
	//stunServer = "stun.syncthing.net:3478"
	//stunServer        = "stun.voipgate.com:3478"
	timeout           = 3
	messageHeaderSize = 20
)

const (
	noNat                   = "no nat"
	endpointIndependent     = "endpoint independent"
	addressDependent        = "address dependent"
	addressAndPortDependent = "address and port dependent"
)

var (
	errResponseMessage = errors.New("error reading from response message channel")
	errTimedOut        = errors.New("timed out waiting for response")
	errNoOtherAddress  = errors.New("no OTHER-ADDRESS in message")
)

func NewPionStunImpl() Tester {
	return &PionStunImpl{}
}

type PionStunImpl struct {
}

func (p *PionStunImpl) GetNATType() (NATType, error) {
	mappingResult, err := mappingTests(stunServer)
	if err != nil {
		return UnKnown, err
	}

	filteringResult, err := filteringTests(stunServer)
	if err != nil {
		return UnKnown, err
	}

	// Full Cone
	if mappingResult == filteringResult && mappingResult == endpointIndependent {
		return FullCone, nil
	}

	// Restricted Cone
	if mappingResult == endpointIndependent && filteringResult == addressDependent {
		return RestrictedCone, nil
	}

	// Port Restricted Cone
	if mappingResult == endpointIndependent && filteringResult == addressAndPortDependent {
		return PortRestrictedCone, nil
	}

	// Symmetric
	if mappingResult == addressAndPortDependent && filteringResult == addressAndPortDependent {
		return Symmetric, nil
	}

	return UnKnown, nil
}

// RFC5780: 4.3.  Determining NAT Mapping Behavior
func mappingTests(addrStr string) (string, error) {
	mapTestConn, err := connect(addrStr)
	if err != nil {
		logger.Warnf("Error creating STUN connection: %s", err)
		return "", err
	}

	defer mapTestConn.Close()

	// Test I: Regular binding request
	logger.Info("Mapping Test I: Regular binding request")
	request := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	resp, err := mapTestConn.roundTrip(request, mapTestConn.RemoteAddr)
	if err != nil {
		return "", err
	}

	// Parse response message for XOR-MAPPED-ADDRESS and make sure OTHER-ADDRESS valid
	resps1 := parse(resp)
	if resps1.xorAddr == nil || resps1.otherAddr == nil {
		logger.Info("Error: NAT discovery feature not supported by this server")
		return "", errNoOtherAddress
	}
	addr, err := net.ResolveUDPAddr("udp4", resps1.otherAddr.String())
	if err != nil {
		logger.Infof("Failed resolving OTHER-ADDRESS: %v", resps1.otherAddr)
		return "", err
	}
	mapTestConn.OtherAddr = addr
	logger.Infof("Received XOR-MAPPED-ADDRESS: %v", resps1.xorAddr)

	// Assert mapping behavior
	if resps1.xorAddr.String() == mapTestConn.LocalAddr.String() {
		logger.Debug("=> NAT mapping behavior: endpoint independent (no NAT)")
		return noNat, nil
	}

	// Test II: Send binding request to the other address but primary port
	logger.Info("Mapping Test II: Send binding request to the other address but primary port")
	oaddr := *mapTestConn.OtherAddr
	oaddr.Port = mapTestConn.RemoteAddr.Port
	resp, err = mapTestConn.roundTrip(request, &oaddr)
	if err != nil {
		return "", err
	}

	// Assert mapping behavior
	resps2 := parse(resp)
	logger.Infof("Received XOR-MAPPED-ADDRESS: %v", resps2.xorAddr)
	if resps2.xorAddr.String() == resps1.xorAddr.String() {
		logger.Debug("=> NAT mapping behavior: endpoint independent")
		return endpointIndependent, nil
	}

	// Test III: Send binding request to the other address and port
	logger.Info("Mapping Test III: Send binding request to the other address and port")
	resp, err = mapTestConn.roundTrip(request, mapTestConn.OtherAddr)
	if err != nil {
		return "", err
	}

	// Assert mapping behavior
	resps3 := parse(resp)
	logger.Infof("Received XOR-MAPPED-ADDRESS: %v", resps3.xorAddr)
	if resps3.xorAddr.String() == resps2.xorAddr.String() {
		logger.Debug("=> NAT mapping behavior: address dependent")
		return addressDependent, nil
	} else {
		logger.Debug("=> NAT mapping behavior: address and port dependent")
		return addressAndPortDependent, nil
	}
}

// RFC5780: 4.4.  Determining NAT Filtering Behavior
func filteringTests(addrStr string) (string, error) {
	mapTestConn, err := connect(addrStr)
	if err != nil {
		logger.Warnf("Error creating STUN connection: %s", err)
		return "", err
	}

	defer mapTestConn.Close()

	// Test I: Regular binding request
	logger.Info("Filtering Test I: Regular binding request")
	request := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	resp, err := mapTestConn.roundTrip(request, mapTestConn.RemoteAddr)
	if err != nil || errors.Is(err, errTimedOut) {
		return "", err
	}
	resps := parse(resp)
	if resps.xorAddr == nil || resps.otherAddr == nil {
		logger.Warn("Error: NAT discovery feature not supported by this server")
		return "", errNoOtherAddress
	}
	addr, err := net.ResolveUDPAddr("udp4", resps.otherAddr.String())
	if err != nil {
		logger.Infof("Failed resolving OTHER-ADDRESS: %v", resps.otherAddr)
		return "", err
	}
	mapTestConn.OtherAddr = addr

	// Test II: Request to change both IP and port
	logger.Info("Filtering Test II: Request to change both IP and port")
	request = stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	request.Add(stun.AttrChangeRequest, []byte{0x00, 0x00, 0x00, 0x06})

	resp, err = mapTestConn.roundTrip(request, mapTestConn.RemoteAddr)
	if err == nil {
		parse(resp) // just to print out the resp
		logger.Debug("=> NAT filtering behavior: endpoint independent")
		return endpointIndependent, nil
	} else if !errors.Is(err, errTimedOut) {
		return "", err // something else went wrong
	}

	// Test III: Request to change port only
	logger.Info("Filtering Test III: Request to change port only")
	request = stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	request.Add(stun.AttrChangeRequest, []byte{0x00, 0x00, 0x00, 0x02})

	resp, err = mapTestConn.roundTrip(request, mapTestConn.RemoteAddr)
	if err == nil {
		parse(resp) // just to print out the resp
		logger.Debug("=> NAT filtering behavior: address dependent")
		return addressDependent, nil
	} else if errors.Is(err, errTimedOut) {
		logger.Debug("=> NAT filtering behavior: address and port dependent")
		return addressAndPortDependent, nil
	}

	return "", err
}

type stunServerConn struct {
	conn        net.PacketConn
	LocalAddr   net.Addr
	RemoteAddr  *net.UDPAddr
	OtherAddr   *net.UDPAddr
	messageChan chan *stun.Message
}

func (c *stunServerConn) Close() error {
	return c.conn.Close()
}

// Send request and wait for response or timeout
func (c *stunServerConn) roundTrip(msg *stun.Message, addr net.Addr) (*stun.Message, error) {
	_ = msg.NewTransactionID()
	logger.Infof("Sending to %v: (%v bytes)", addr, msg.Length+messageHeaderSize)
	logger.Debugf("%v", msg)
	for _, attr := range msg.Attributes {
		logger.Debugf("\t%v (l=%v)", attr, attr.Length)
	}
	_, err := c.conn.WriteTo(msg.Raw, addr)
	if err != nil {
		logger.Warnf("Error sending request to %v", addr)
		return nil, err
	}

	// Wait for response or timeout
	select {
	case m, ok := <-c.messageChan:
		if !ok {
			return nil, errResponseMessage
		}
		return m, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		logger.Infof("Timed out waiting for response from server %v", addr)
		return nil, errTimedOut
	}
}

// Given an address string, returns a StunServerConn
func connect(addrStr string) (*stunServerConn, error) {
	logger.Infof("Connecting to STUN server: %s", addrStr)
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		logger.Warnf("Error resolving address: %s", err)
		return nil, err
	}

	c, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, err
	}
	logger.Infof("Local address: %s", c.LocalAddr())
	logger.Infof("Remote address: %s", addr.String())

	mChan := listen(c)

	return &stunServerConn{
		conn:        c,
		LocalAddr:   c.LocalAddr(),
		RemoteAddr:  addr,
		messageChan: mChan,
	}, nil
}

func listen(conn *net.UDPConn) (messages chan *stun.Message) {
	messages = make(chan *stun.Message)
	go func() {
		for {
			buf := make([]byte, 1024)

			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				close(messages)
				return
			}
			logger.Infof("Response from %v: (%v bytes)", addr, n)
			buf = buf[:n]

			m := new(stun.Message)
			m.Raw = buf
			err = m.Decode()
			if err != nil {
				logger.Infof("Error decoding message: %v", err)
				close(messages)
				return
			}

			messages <- m
		}
	}()
	return
}

// Parse a STUN message
func parse(msg *stun.Message) (ret struct {
	xorAddr    *stun.XORMappedAddress
	otherAddr  *stun.OtherAddress
	respOrigin *stun.ResponseOrigin
	mappedAddr *stun.MappedAddress
	software   *stun.Software
},
) {
	ret.mappedAddr = &stun.MappedAddress{}
	ret.xorAddr = &stun.XORMappedAddress{}
	ret.respOrigin = &stun.ResponseOrigin{}
	ret.otherAddr = &stun.OtherAddress{}
	ret.software = &stun.Software{}
	if ret.xorAddr.GetFrom(msg) != nil {
		ret.xorAddr = nil
	}
	if ret.otherAddr.GetFrom(msg) != nil {
		ret.otherAddr = nil
	}
	if ret.respOrigin.GetFrom(msg) != nil {
		ret.respOrigin = nil
	}
	if ret.mappedAddr.GetFrom(msg) != nil {
		ret.mappedAddr = nil
	}
	if ret.software.GetFrom(msg) != nil {
		ret.software = nil
	}
	logger.Debugf("%v", msg)
	logger.Debugf("\tMAPPED-ADDRESS:     %v", ret.mappedAddr)
	logger.Debugf("\tXOR-MAPPED-ADDRESS: %v", ret.xorAddr)
	logger.Debugf("\tRESPONSE-ORIGIN:    %v", ret.respOrigin)
	logger.Debugf("\tOTHER-ADDRESS:      %v", ret.otherAddr)
	logger.Debugf("\tSOFTWARE: %v", ret.software)
	for _, attr := range msg.Attributes {
		switch attr.Type {
		case
			stun.AttrXORMappedAddress,
			stun.AttrOtherAddress,
			stun.AttrResponseOrigin,
			stun.AttrMappedAddress,
			stun.AttrSoftware:
			break
		default:
			logger.Debugf("\t%v (l=%v)", attr, attr.Length)
		}
	}
	return ret
}
