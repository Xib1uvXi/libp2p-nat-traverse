package netutil

import (
	"errors"
	"net"
	"sync"
	"time"
)

var timeout = 30 * time.Second

var ErrConnHandlerNotSet = errors.New("connection handler not set")

type WrapUDPConnHandler func(conn net.Conn) error

type WrapUDPConn struct {
	ListenAddr     string
	Conn           *net.UDPConn
	listener       net.Listener
	connHandler    WrapUDPConnHandler
	setHandlerOnce sync.Once
}

// NewWrapUDPConn creates a new WrapUDPConn
func NewWrapUDPConn(listenAddr string) (*WrapUDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", listenAddr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}

	// wrap the UDPConn in a net.Listener, use our custom listener
	ln, err := ListenFromPacketConn(udpConn)

	if err != nil {
		return nil, err
	}

	return &WrapUDPConn{
		ListenAddr: listenAddr,
		Conn:       udpConn,
		listener:   ln,
	}, nil
}

// SendUnreliableUDPMessage sends an unreliable UDP message
func (w *WrapUDPConn) SendUnreliableUDPMessage(data []byte, raddr net.Addr) error {
	_, err := w.Conn.WriteToUDP(data, raddr.(*net.UDPAddr))
	return err
}

// SetConnHandler sets the connection handler
func (w *WrapUDPConn) SetConnHandler(handler WrapUDPConnHandler) {
	w.setHandlerOnce.Do(func() {
		w.connHandler = handler
	})
}

// Start starts the listener
func (w *WrapUDPConn) Start() error {
	if w.connHandler == nil {
		return ErrConnHandlerNotSet
	}

	go w.acceptLoop()

	return nil
}

// Close closes the listener
func (w *WrapUDPConn) Close() error {
	if err := w.listener.Close(); err != nil {
		return err
	}

	return nil
}

// accept loops and accepts connections
func (w *WrapUDPConn) acceptLoop() {
	for {
		conn, err := w.listener.Accept()
		if err != nil {
			logger.Errorf("failed to accept connection: %s", err.Error())
			return
		}

		go w.handleConn(conn)
	}
}

// handleConn handles a connection
func (w *WrapUDPConn) handleConn(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	if err := w.connHandler(conn); err != nil {
		logger.Errorf("error handling connection: %s", err.Error())
		return
	}
}
