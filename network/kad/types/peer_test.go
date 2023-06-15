package types

import (
	"net"
	"sleepy/types"
	"testing"
)

func TestPeer_GetIP(t *testing.T) {
	testIp := net.ParseIP("100.101.102.103")
	peer := &peerImp{
		id: types.NewUInt128FromInt(1),
		ip: testIp,
	}

	if !testIp.Equal(peer.GetIP()) {
		t.Errorf("IP missmatch")
	}
}

func TestPeer_SetIP(t *testing.T) {
	testIp := net.ParseIP("100.101.102.103")
	peer := &peerImp{
		id: types.NewUInt128FromInt(1),
	}

	peer.SetIP(testIp, true)

	if !testIp.Equal(peer.ip) || !peer.IsIPVerified() {
		t.Errorf("IP os verify state missmatch")
	}
}
