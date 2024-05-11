package netutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestNewWrapUDPConn(t *testing.T) {
	testMsg := "hello"
	testHit := false

	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000
	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)
	srvAddr, err := net.ResolveUDPAddr("udp4", srvAddrStr)
	assert.NoError(t, err)

	wrap, err := NewWrapUDPConn(srvAddrStr)
	assert.NoError(t, err)

	assert.EqualError(t, wrap.Start(), ErrConnHandlerNotSet.Error())

	// set the connection handler
	handler := func(conn net.Conn) error {
		var data = make([]byte, 1024)
		n, err := conn.Read(data)
		assert.NoError(t, err)
		assert.Equal(t, testMsg, string(data[:n]))
		testHit = true
		return nil
	}

	wrap.SetConnHandler(handler)

	assert.NoError(t, wrap.Start())

	clientAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", randomPort+30))
	assert.NoError(t, err)

	clientConn, err := net.DialUDP("udp4", clientAddr, srvAddr)
	assert.NoError(t, err)

	_, err = clientConn.Write([]byte(testMsg))
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	assert.True(t, testHit)

	assert.NoError(t, wrap.Conn.Close())
}

func TestWrapUDPConn_SendUnreliableUDPMessage(t *testing.T) {
	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000
	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)

	wrap, err := NewWrapUDPConn(srvAddrStr)
	assert.NoError(t, err)

	targetAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", randomPort+30))
	assert.NoError(t, err)

	err = wrap.SendUnreliableUDPMessage([]byte("hello"), targetAddr)
	assert.NoError(t, err)

}
