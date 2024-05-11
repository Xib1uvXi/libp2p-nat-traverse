package natinfo

import (
	"context"
	"fmt"
	"github.com/Xib1uvXi/libp2p-nat-traverse/netutil"
	"github.com/Xib1uvXi/libp2p-nat-traverse/netutil/quicutil"
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/quic-go/quic-go"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestNewPortNegotiationRecorder(t *testing.T) {
	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000
	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)

	srvAddr, err := net.ResolveUDPAddr("udp", srvAddrStr)
	assert.NoError(t, err)

	server, err := netutil.NewReliableUDPServer(srvAddrStr)
	assert.NoError(t, err)

	recorder := NewPortNegotiationRecorder()
	server.SetStreamHandler(recorder.PortNegotiationHandler)

	assert.NoError(t, server.Start())
	// client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", randomPort+30))
	assert.NoError(t, err)

	clientUDPConn, err := net.ListenUDP("udp", clientAddr)

	testConn, err := quic.Dial(ctx, clientUDPConn, srvAddr, quicutil.ClientTLSConfig(), nil)
	assert.NoError(t, err)

	stream, err := testConn.OpenStreamSync(ctx)
	assert.NoError(t, err)

	assert.NoError(t, pbio.NewDelimitedWriter(stream).WriteMsg(&pb.Message{Type: pb.MsgType_PortNegotiationResponse}))

	raddr, err := recorder.Await()
	assert.NoError(t, err)

	assert.Equal(t, testConn.LocalAddr().String(), raddr.String())

	assert.NoError(t, server.Close())
}
