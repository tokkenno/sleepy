package router

import (
	"github.com/tokkenno/sleepy/network/kad"
	"github.com/tokkenno/sleepy/types"
	"math/rand"
	"net"
	"testing"
)

func TestZone_AddPeer(t *testing.T) {
	routerId := types.NewUInt128FromInt(0xff00ff)
	randGen := rand.New(rand.NewSource(0))
	zone := NewRouter(routerId)

	// Add maxBucketSize + 1 to force new bucket creation
	for i := 0; i < maxBucketSize+1; i++ {
		peer := kad.NewPeer(types.NewUInt128(randGen.Uint64(), randGen.Uint64()))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		err := zone.AddPeer(peer)

		if err != nil {
			t.Errorf("Unexpected error when add peer: %s", err.Error())
		} else if zone.ContainsPeer(peer.Id()) {
			peerIn, err := zone.GetPeer(peer.Id())
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if peerIn == nil {
				t.Errorf("Peer can't be null")
			} else if !peerIn.Equal(peer) {
				t.Errorf("Z peer not equal, 0x%s expected, 0x%s found", peer.Id().ToHexString(), peerIn.Id().ToHexString())
			}
		} else {
			t.Errorf("Zone must contains the peer, but not found.")
		}
	}
}

func TestZone_AddSamePeer(t *testing.T) {
	routerId := types.NewUInt128FromInt(0xff00ff)
	zone := NewRouter(routerId)
	randGen := rand.New(rand.NewSource(0))

	// Add maxBucketSize + 1 to force new bucket creation
	for i := 0; i < maxBucketSize+5; i++ {
		peer := kad.NewPeer(types.NewUInt128FromInt(0))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		err := zone.AddPeer(peer)

		if err != nil {
			t.Errorf("Unexpected error when add peer: %s", err.Error())
		}
	}

	if zone.CountPeers() > 1 {
		t.Errorf("This zone only must have %d peers, %d given", 1, zone.CountPeers())
	}
}

func TestZone_ContainsPeer(t *testing.T) {
	routerId := types.NewUInt128FromInt(0xff00ff00)
	router := NewRouter(routerId)

	peerId := types.NewUInt128FromInt(0xff00ff)
	peer := kad.NewPeer(peerId)

	err := router.AddPeer(peer)

	if err != nil {
		t.Errorf("Unexpected error when add peer: %s", err.Error())
	} else if !router.ContainsPeer(peerId) {
		t.Errorf("Router must contains the peer: %s", peerId.ToHexString())
	}
}

func TestZone_CountPeers(t *testing.T) {
	routerId := types.NewUInt128FromInt(0xff00ff)
	randGen := rand.New(rand.NewSource(0))
	zone := NewRouter(routerId)

	for i := 1; i <= maxBucketSize; i++ {
		peer := kad.NewPeer(types.NewUInt128FromInt(i))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		err := zone.AddPeer(peer)

		if err != nil {
			t.Errorf("Unexpected error when add peer: %s", err.Error())
		} else if zone.CountPeers() != i {
			t.Errorf("Zone must contains %d peers, %d found", i, zone.CountPeers())
		}
	}
}
