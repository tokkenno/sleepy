package router

import (
	"errors"
	"math/rand"
	kadTypes "sleepy/network/kad/types"
	"sleepy/types"
	"sleepy/utils/event"
	"time"
)

type Router interface {
	// AddPeer adds a peer to the router
	AddPeer(peer kadTypes.Peer) error
	// GetBootstrapPeers returns a list of peers of [max] size prepared to do a bootstrap
	GetBootstrapPeers(max int, closedTo types.UInt128) []kadTypes.Peer
}

// The router is the special zone in the root of a zone tree
type routerImp struct {
	zone
	randomGenerator        *rand.Rand
	peerUpdateRequestEvent *event.Emitter
	peerLookupRequestEvent *event.Emitter
}

var _ Router = &routerImp{}

// Load a router zone tree from file
func LoadRouterFromFile(localId types.UInt128, path string) (*routerImp, error) {
	return nil, errors.New("not implemented yet")
}

// Create a new zone tree (routerImp) from the local peer GetID
func NewRouter(id types.UInt128) Router {
	rz := &routerImp{
		zone: zone{
			localId:           id.Clone(),
			zoneIndex:         types.NewUInt128FromInt(0),
			parent:            nil,
			leftChild:         nil,
			rightChild:        nil,
			level:             0,
			bucket:            newKBucket(),
			randomLookupTimer: nil,
		},
		randomGenerator:        rand.New(rand.NewSource(time.Now().UnixNano())),
		peerUpdateRequestEvent: event.NewEvent(),
		peerLookupRequestEvent: event.NewEvent(),
	}

	rz.zone.root = rz

	rz.startChecks()

	return rz
}

// Event fired when the router need update a peer information
func (router *routerImp) PeerUpdateRequestEvent() *event.Handler {
	return router.peerUpdateRequestEvent.GetHandler()
}

// Event fired when the router need lookup a peer
func (router *routerImp) PeerLookupRequestEvent() *event.Handler {
	return router.peerLookupRequestEvent.GetHandler()
}

func (router *routerImp) SaveFile(path string) error {
	return errors.New("not implemented yet")
}

func (router *routerImp) GetBootstrapPeers(max int, _ types.UInt128) []kadTypes.Peer {
	const BootstrapDepth = 5 // Defined as LOG_BASE_EXPONENT constant in protocol/defines.h
	return router.GetTopPeers(max, BootstrapDepth)
}
