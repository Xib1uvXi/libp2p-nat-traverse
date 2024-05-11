package netutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestListenFromPacketConn(t *testing.T) {
	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000
	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)

	srvAddr, err := net.ResolveUDPAddr("udp4", srvAddrStr)
	assert.NoError(t, err)

	udpConn, err := net.ListenUDP("udp4", srvAddr)
	assert.NoError(t, err)

	ln, err := ListenFromPacketConn(udpConn)
	assert.NoError(t, err)

	assert.NoError(t, ln.Close())
}

func TestListen(t *testing.T) {
	testMsg := "hello"
	testHit := false

	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000
	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)
	srvAddr, err := net.ResolveUDPAddr("udp4", srvAddrStr)
	assert.NoError(t, err)

	ln, err := Listen("udp4", srvAddr)
	assert.NoError(t, err)

	go func() {
		conn, err := ln.Accept()
		assert.NoError(t, err)

		var data = make([]byte, 1024)
		n, err := conn.Read(data)
		assert.NoError(t, err)
		assert.Equal(t, testMsg, string(data[:n]))
		testHit = true
	}()

	clientAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", randomPort+30))
	assert.NoError(t, err)

	clientConn, err := net.DialUDP("udp4", clientAddr, srvAddr)
	assert.NoError(t, err)

	_, err = clientConn.Write([]byte(testMsg))
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	assert.True(t, testHit)

	assert.NoError(t, ln.Close())
}
