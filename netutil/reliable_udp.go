package netutil

import (
	"context"
	"errors"
	"github.com/Xib1uvXi/libp2p-nat-traverse/netutil/quicutil"
	"github.com/quic-go/quic-go"
	"net"
	"sync"
	"time"
)

var (
	streamHandlerTimeout = 30 * time.Second
)

var (
	ErrStreamHandlerNotSet = errors.New("stream handler not set")
)

type ReliableUDPStreamHandler = func(conn quic.Connection, stream quic.Stream) error

type ReliableUDPServer struct {
	ctx       context.Context
	cancelFuc context.CancelFunc

	ListenAddr     string
	UDPConn        *net.UDPConn
	quicListener   *quic.Listener
	streamHandler  ReliableUDPStreamHandler
	setHandlerOnce sync.Once
	stopC          chan struct{}
}

func NewReliableUDPServer(listenAddr string) (*ReliableUDPServer, error) {
	addr, err := net.ResolveUDPAddr("udp4", listenAddr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}

	quicListener, err := quic.Listen(udpConn, quicutil.ServerTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ReliableUDPServer{
		ctx:          ctx,
		cancelFuc:    cancel,
		ListenAddr:   listenAddr,
		UDPConn:      udpConn,
		quicListener: quicListener,
		stopC:        make(chan struct{}),
	}, nil
}

// SetStreamHandler sets the stream handler
func (s *ReliableUDPServer) SetStreamHandler(handler ReliableUDPStreamHandler) {
	s.setHandlerOnce.Do(func() {
		s.streamHandler = handler
	})
}

// Start starts the server
func (s *ReliableUDPServer) Start() error {
	if s.streamHandler == nil {
		return ErrStreamHandlerNotSet
	}

	go s.acceptLoop()

	return nil
}

// Close closes the server
func (s *ReliableUDPServer) Close() error {
	s.cancelFuc()
	close(s.stopC)

	if err := s.quicListener.Close(); err != nil {
		return err
	}

	if err := s.UDPConn.Close(); err != nil {
		return err
	}

	return nil
}

func (s *ReliableUDPServer) acceptLoop() {
	for {
		select {
		case <-s.stopC:
			logger.Info("quic accept loop stop")
			return
		default:
		}

		conn, err := s.quicListener.Accept(s.ctx)
		if err != nil {
			logger.Warnf("accept quic error: %s", err)
			continue
		}

		go s.handleConn(conn)
	}
}

// handle quic connection
func (s *ReliableUDPServer) handleConn(quicConn quic.Connection) {
	for {
		select {
		case <-s.stopC:
			logger.Info("quic accept stream loop stop")
			return
		default:
		}

		stream, err := quicConn.AcceptStream(s.ctx)
		if err != nil {
			logger.Warnf("accept stream error: %s", err)
			continue
		}

		go s.acceptSteamLoop(quicConn, stream)
	}

}

func (s *ReliableUDPServer) acceptSteamLoop(quicConn quic.Connection, stream quic.Stream) {
	defer stream.Close()
	_ = stream.SetDeadline(time.Now().Add(streamHandlerTimeout))

	if err := s.streamHandler(quicConn, stream); err != nil {
		logger.Errorf("stream handler error: %s", err)
	}

	return
}
