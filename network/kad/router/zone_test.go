package router

import (
	"testing"
	"github.com/reimashi/sleepy/types"
	"github.com/reimashi/sleepy/network/kad"
	"net"
	"math/rand"
)

func TestZone_AddPeer(t *testing.T) {
	routerId := *types.NewUInt128FromInt(0xff00ff)
	randGen := rand.New(rand.NewSource(0))
	zone := FromId(routerId)

	for i := 0; i < maxSize + 1; i++ {
		peer := kad.NewPeer(*types.NewUInt128FromInt(i))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		zone.AddPeer(peer)

		if zone.ContainsPeer(*peer.GetId()) {
			peerIn, err := zone.GetPeer(*peer.GetId())
			if err != nil || peerIn == nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if !peerIn.Equal(*peer) {
				t.Errorf("Z peer not equal, 0x%s expected, 0x%s found", peer.GetId().ToHexString(), peerIn.GetId().ToHexString())
			}
		} else {
			t.Errorf("Zone must contains the peer, but not found.")
		}
	}
}

func TestZone_CountPeers(t *testing.T) {
	routerId := *types.NewUInt128FromInt(0xff00ff)
	randGen := rand.New(rand.NewSource(0))
	zone := FromId(routerId)

	for i := 0; i < maxSize + 1; i++ {
		peer := kad.NewPeer(*types.NewUInt128FromInt(i))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		zone.AddPeer(peer)

		if zone.CountPeers() != i+1 {
			t.Errorf("Zone must contains %d peers, %d found", i+1, zone.CountPeers())
		}
	}
}