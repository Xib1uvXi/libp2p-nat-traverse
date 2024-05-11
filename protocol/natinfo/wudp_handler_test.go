package natinfo

import (
	"fmt"
	"github.com/Xib1uvXi/libp2p-nat-traverse/netutil"
	"github.com/Xib1uvXi/libp2p-nat-traverse/protocol/natinfo/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestNewPortRestrictedChecker(t *testing.T) {
	rand.New(rand.NewSource(time.Now().Unix()))
	randomPort := rand.Intn(10000) + 10000

	srvAddrStr := fmt.Sprintf("127.0.0.1:%d", randomPort)
	srvAddr, err := net.ResolveUDPAddr("udp4", srvAddrStr)
	assert.NoError(t, err)

	wrapUDP, err := netutil.NewWrapUDPConn(srvAddrStr)
	assert.NoError(t, err)

	portRestrictedChecker := NewPortRestrictedChecker()
	wrapUDP.SetConnHandler(portRestrictedChecker.PortRestrictedHandler)

	assert.NoError(t, wrapUDP.Start())

	clientAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", randomPort+30))
	assert.NoError(t, err)

	clientConn, err := net.ListenUDP("udp4", clientAddr)
	assert.NoError(t, err)

	data, err := proto.Marshal(&pb.Message{Type: pb.MsgType_ServerPortChangeTestResponse})
	assert.NoError(t, err)

	clientConn.WriteToUDP(data, srvAddr)
	_, err = clientConn.WriteToUDP(data, srvAddr)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	select {
	case <-time.After(2 * time.Second):
		assert.Fail(t, "timeout")
	case err := <-portRestrictedChecker.ErrC:
		assert.Fail(t, err.Error())
	case <-portRestrictedChecker.DoneC:
	}

	assert.NoError(t, wrapUDP.Close())
}
