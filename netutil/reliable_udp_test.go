package netutil

import (
	"context"
	"fmt"
	"github.com/Xib1uvXi/libp2p-nat-traverse/netutil/quicutil"
	"github.com/quic-go/quic-go"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestNewReliableUDPServer(t *testing.T) {
	testMsg := "hello"
	testHit := false

	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000

	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)

	server, err := NewReliableUDPServer(srvAddrStr)
	assert.NoError(t, err)

	serverHandler := func(conn quic.Connection, stream quic.Stream) error {
		var msg = make([]byte, 1024)

		n, err := stream.Read(msg)
		assert.NoError(t, err)

		assert.Equal(t, testMsg, string(msg[:n]))
		testHit = true
		return nil
	}

	server.SetStreamHandler(serverHandler)

	assert.NoError(t, server.Start())

	// client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testConn, err := quic.DialAddr(ctx, srvAddrStr, quicutil.ClientTLSConfig(), nil)
	assert.NoError(t, err)

	stream, err := testConn.OpenStreamSync(ctx)
	assert.NoError(t, err)

	_, err = stream.Write([]byte(testMsg))
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	assert.True(t, testHit)

	assert.NoError(t, server.Close())
}
