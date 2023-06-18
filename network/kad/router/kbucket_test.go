package router

import (
	"math/rand"
	"net"
	types2 "sleepy/network/kad/types"
	"sleepy/types"
	"strconv"
	"testing"
)

func TestKBucket_AddPeer(t *testing.T) {
	peer := types2.NewPeer(types.NewUInt128FromInt(1))
	kBucket := newKBucket()
	err := kBucket.AddPeer(peer)

	if err != nil {
		t.Errorf("Error when add peer to K-Bucket")
	}

	if kBucket.CountPeers() != 1 {
		t.Errorf("K-Bucket must contains a unique peer, %d found", kBucket.CountPeers())
	}

	peerIn := kBucket.peers[0]
	if !peerIn.Equal(peer) {
		t.Errorf("K-Bucket peer not equal, 0x%s expected, 0x%s found", peer.GetID().ToHexString(), peerIn.GetID().ToHexString())
	}
}

func TestKBucket_AddMultiplePeer(t *testing.T) {
	kBucket := newKBucket()
	randGen := rand.New(rand.NewSource(0))

	for i := 0; i < maxBucketSize; i++ {
		peer := types2.NewPeer(types.NewUInt128FromInt(i))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		kBucket.AddPeer(peer)

		if kBucket.CountPeers() != i+1 {
			t.Errorf("K-Bucket must contains %d peers, %d found", i+1, kBucket.CountPeers())
		}
	}
}

func TestKBucket_Count(t *testing.T) {
	peer := types2.NewPeer(types.NewUInt128FromInt(1))
	kBucket := newKBucket()
	kBucket.AddPeer(peer)

	if kBucket.CountPeers() != 1 {
		t.Errorf("K-Bucket must contains a unique peer, %d found", kBucket.CountPeers())
	}

	kBucket.AddPeer(peer)

	if kBucket.CountPeers() != 1 {
		t.Errorf("K-Bucket must contains a unique peer, %d found", kBucket.CountPeers())
	}

	otherPeer := types2.NewPeer(types.NewUInt128FromInt(2))
	kBucket.AddPeer(otherPeer)

	if kBucket.CountPeers() != 2 {
		t.Errorf("K-Bucket must contains two unique peer, %d found", kBucket.CountPeers())
	}
}

func TestKBucket_GetPeer(t *testing.T) {
	id := types.NewUInt128FromInt(1)
	peer := types2.NewPeer(id)
	kBucket := newKBucket()
	kBucket.AddPeer(peer)

	gpeer, err := kBucket.GetPeer(id)
	if err != nil {
		t.Errorf(err.Error())
	} else if gpeer == nil {
		t.Errorf("Peer not found in K-Bucket")
	} else if !gpeer.Equal(peer) {
		t.Errorf("Peer retrieved is distinct than the original")
	}
}

func TestKBucket_GetPeerByAddr(t *testing.T) {
	testIp := net.ParseIP("100.101.102.103")
	peer := types2.NewPeer(types.NewUInt128FromInt(1))
	peer.SetIP(testIp, false)
	peer.SetTCPPort(123)
	peer.SetUDPPort(321)

	kBucket := newKBucket()
	kBucket.AddPeer(peer)

	tip := &net.TCPAddr{IP: testIp, Port: 123}
	gPeer, err := kBucket.GetPeerByAddr(tip)
	if err != nil {
		t.Errorf(err.Error())
	} else if gPeer == nil {
		t.Errorf("Peer not found in K-Bucket")
	} else if !gPeer.Equal(peer) {
		t.Errorf("Peer retrieved is distinct than the original")
	}

	uip := &net.UDPAddr{IP: testIp, Port: 321}
	gPeer, err = kBucket.GetPeerByAddr(uip)
	if err != nil {
		t.Errorf(err.Error())
	} else if gPeer == nil {
		t.Errorf("Peer not found in K-Bucket")
	} else if !gPeer.Equal(peer) {
		t.Errorf("Peer retrieved is distinct than the original")
	}
}

func TestKBucket_GetPeers(t *testing.T) {
	peer0 := types2.NewPeer(types.NewUInt128FromInt(1))
	peer1 := types2.NewPeer(types.NewUInt128FromInt(2))
	peer2 := types2.NewPeer(types.NewUInt128FromInt(3))

	kBucket := newKBucket()
	kBucket.AddPeer(peer0)
	kBucket.AddPeer(peer1)
	kBucket.AddPeer(peer2)

	if len(kBucket.Peers()) != 3 {
		t.Errorf("K-Bucket must contains 3 unique peers")
	} else if !kBucket.peers[0].Equal(peer0) || !kBucket.peers[1].Equal(peer1) || !kBucket.peers[2].Equal(peer2) {
		t.Errorf("Peers are not in correct order")
	}
}

func TestKBucket_RemovePeer(t *testing.T) {
	peer := types2.NewPeer(types.NewUInt128FromInt(1))
	kBucket := newKBucket()
	kBucket.AddPeer(peer)
	kBucket.RemovePeer(peer)

	if kBucket.CountPeers() != 0 {
		t.Errorf("K-Bucket must be empty, but contains %v elements", kBucket.CountPeers())
	}
}

func TestKBucket_IsFull(t *testing.T) {
	kBucket := newKBucket()
	var peers [maxBucketSize]types2.Peer

	if kBucket.IsFull() {
		t.Errorf("K-Bucket must not be full")
	}

	for i := 0; i < maxBucketSize; i++ {
		newPeer := types2.NewPeer(types.NewUInt128FromInt(i))
		newPeer.SetIP(net.ParseIP("100.101.102."+strconv.Itoa(i)), false) // Limit of peers with the same ip
		peers[i] = newPeer
		kBucket.AddPeer(newPeer)
	}

	if !kBucket.IsFull() {
		t.Errorf("K-Bucket must be full (%v) but contains %v", maxBucketSize, kBucket.CountPeers())
	}
}

func TestKBucket_SetAlive(t *testing.T) {
	peer0 := types2.NewPeer(types.NewUInt128FromInt(1))
	peer1 := types2.NewPeer(types.NewUInt128FromInt(2))

	kBucket := newKBucket()
	kBucket.AddPeer(peer0)
	kBucket.AddPeer(peer1)

	kBucket.SetPeerAlive(peer0.GetID())

	if len(kBucket.Peers()) != 2 {
		t.Errorf("K-Bucket must contains 2 unique peers")
	} else if !kBucket.peers[1].Equal(peer0) {
		t.Errorf("Peers are not in correct order")
	}
}

func TestKBucket_GetRandomPeer(t *testing.T) {
	peer0 := types2.NewPeer(types.NewUInt128FromInt(1))
	peer1 := types2.NewPeer(types.NewUInt128FromInt(2))

	kBucket := newKBucket()
	kBucket.AddPeer(peer0)
	kBucket.AddPeer(peer1)

	count := 0
	for i := 0; i < 100; i++ {
		rpeer, _ := kBucket.GetRandomPeer()
		if rpeer.Equal(peer1) {
			count++
		}
	}

	if count == 0 || count == 100 {
		t.Errorf("Probably not random, but hey!, nothing is impossible")
	}
}
