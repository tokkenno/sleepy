package router

import (
	"testing"
	"github.com/tokkenno/sleepy/network/kad"
	"github.com/tokkenno/sleepy/types"
	"net"
	"strconv"
	"math/rand"
)

func TestKBucket_AddPeer(t *testing.T) {
	peer := kad.NewPeer(*types.NewUInt128FromInt(1))
	kbucket := &KBucket{ }
	kbucket.AddPeer(peer)

	if kbucket.CountPeers() != 1 {
		t.Errorf("K-Bucket must contains a unique peer, %d found", kbucket.CountPeers())
	}

	peerIn := &kbucket.peers[0]
	if !peerIn.Equal(*peer) {
		t.Errorf("K-Bucket peer not equal, 0x%s expected, 0x%s found", peer.GetId().ToHexString(), peerIn.GetId().ToHexString())
	}
}

func TestKBucket_AddMultiplePeer(t *testing.T) {
	kbucket := &KBucket{ }
	randGen := rand.New(rand.NewSource(0))

	for i := 0; i < maxSize; i++ {
		peer := kad.NewPeer(*types.NewUInt128FromInt(i))
		peer.SetIP(net.IPv4(byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255)), byte(randGen.Intn(255))), false)
		kbucket.AddPeer(peer)

		if kbucket.CountPeers() != i+1 {
			t.Errorf("K-Bucket must contains %d peers, %d found", i+1, kbucket.CountPeers())
		}
	}
}

func TestKBucket_Count(t *testing.T) {
	peer := kad.NewPeer(*types.NewUInt128FromInt(1))
	kbucket := &KBucket{ }
	kbucket.AddPeer(peer)

	if kbucket.CountPeers() != 1 {
		t.Errorf("K-Bucket must contains a unique peer, %d found", kbucket.CountPeers())
	}

	kbucket.AddPeer(peer)

	if kbucket.CountPeers() != 1 {
		t.Errorf("K-Bucket must contains a unique peer, %d found", kbucket.CountPeers())
	}

	otherPeer := kad.NewPeer(*types.NewUInt128FromInt(2))
	kbucket.AddPeer(otherPeer)

	if kbucket.CountPeers() != 2 {
		t.Errorf("K-Bucket must contains two unique peer, %d found", kbucket.CountPeers())
	}
}

func TestKBucket_GetPeer(t *testing.T) {
	id := types.NewUInt128FromInt(1)
	peer := kad.NewPeer(*id)
	kbucket := &KBucket{ }
	kbucket.AddPeer(peer)

	gpeer, err := kbucket.GetPeer(*id)
	if err != nil {
		t.Errorf(err.Error())
	} else if gpeer == nil {
		t.Errorf("Peer not found in K-Bucket")
	} else if !gpeer.Equal(*peer) {
		t.Errorf("Peer retrieved is distinct than the original")
	}
}

func TestKBucket_GetPeerByIp(t *testing.T) {
	testIp := net.ParseIP("100.101.102.103")
	peer := kad.NewPeer(*types.NewUInt128FromInt(1))
	peer.SetIP(testIp, false)
	kbucket := &KBucket{ }
	kbucket.AddPeer(peer)

	gpeer, err := kbucket.GetPeerByIp(testIp)
	if err != nil {
		t.Errorf(err.Error())
	} else if gpeer == nil {
		t.Errorf("Peer not found in K-Bucket")
	} else if !gpeer.Equal(*peer) {
		t.Errorf("Peer retrieved is distinct than the original")
	}
}

func TestKBucket_GetPeers(t *testing.T) {
	peer0 := kad.NewPeer(*types.NewUInt128FromInt(1))
	peer1 := kad.NewPeer(*types.NewUInt128FromInt(2))
	peer2 := kad.NewPeer(*types.NewUInt128FromInt(3))

	kbucket := &KBucket{ }
	kbucket.AddPeer(peer0)
	kbucket.AddPeer(peer1)
	kbucket.AddPeer(peer2)

	if len(kbucket.GetPeers()) != 3 {
		t.Errorf("K-Bucket must contains 3 unique peers")
	} else if !kbucket.peers[0].Equal(*peer0) || !kbucket.peers[1].Equal(*peer1) || !kbucket.peers[2].Equal(*peer2) {
		t.Errorf("Peers are not in correct order")
	}
}

func TestKBucket_RemovePeer(t *testing.T) {
	peer := kad.NewPeer(*types.NewUInt128FromInt(1))
	kbucket := &KBucket{ }
	kbucket.AddPeer(peer)
	kbucket.RemovePeer(peer)

	if kbucket.CountPeers() != 0 {
		t.Errorf("K-Bucket must be empty, but contains %v elements", kbucket.CountPeers())
	}
}

func TestKBucket_IsFull(t *testing.T) {
	kbucket := &KBucket{ }
	var peers [maxSize]kad.Peer

	if kbucket.IsFull() {
		t.Errorf("K-Bucket must not be full")
	}

	for i := 0; i < maxSize; i++ {
		newPeer := kad.NewPeer(*types.NewUInt128FromInt(i))
		newPeer.SetIP(net.ParseIP("100.101.102." + strconv.Itoa(i)), false) // Limit of peers with the same ip
		peers[i] = *newPeer
		kbucket.AddPeer(newPeer)
	}

	if !kbucket.IsFull() {
		t.Errorf("K-Bucket must be full (%v) but contains %v", maxSize, kbucket.CountPeers())
	}
}

func TestKBucket_SetAlive(t *testing.T) {
	peer0 := kad.NewPeer(*types.NewUInt128FromInt(1))
	peer1 := kad.NewPeer(*types.NewUInt128FromInt(2))

	kbucket := &KBucket{ }
	kbucket.AddPeer(peer0)
	kbucket.AddPeer(peer1)

	kbucket.SetPeerAlive(*peer0.GetId())

	if len(kbucket.GetPeers()) != 2 {
		t.Errorf("K-Bucket must contains 2 unique peers")
	} else if !kbucket.peers[1].Equal(*peer0) {
		t.Errorf("Peers are not in correct order")
	}
}

func TestKBucket_GetRandomPeer(t *testing.T) {
	peer0 := kad.NewPeer(*types.NewUInt128FromInt(1))
	peer1 := kad.NewPeer(*types.NewUInt128FromInt(2))

	kbucket := &KBucket{ }
	kbucket.AddPeer(peer0)
	kbucket.AddPeer(peer1)

	count := 0
	for i := 0; i < 100; i++ {
		rpeer, _ := kbucket.GetRandomPeer()
		if rpeer.Equal(*peer1) { count++ }
	}

	if count == 0 || count == 100 {
		t.Errorf("Probably not random, but hey!, nothing is impossible")
	}
}