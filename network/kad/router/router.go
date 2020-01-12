package router

import (
	"errors"
	"math/rand"
	types2 "sleepy/network/kad/types"
	"sleepy/types"
	"sleepy/utils/event"
	"time"
)

// The router is the special zone in the root of a zone tree
type Router struct {
	Zone
	randomGenerator        *rand.Rand
	peerUpdateRequestEvent *event.Emitter
	peerLookupRequestEvent *event.Emitter
}

// Load a router zone tree from file
func LoadRouterFromFile(localId types.UInt128, path string) (*Router, error) {
	return nil, errors.New("not implemented yet")
}

// Create a new Zone tree (Router) from the local peer Id
func NewRouter(id *types.UInt128) *Router {
	rz := &Router{
		Zone: Zone{
			localId:           *id.Clone(),
			zoneIndex:         *types.NewUInt128FromInt(0),
			parent:            nil,
			leftChild:         nil,
			rightChild:        nil,
			level:             0,
			bucket:            new(kBucket),
			randomLookupTimer: nil,
		},
		randomGenerator:        rand.New(rand.NewSource(time.Now().UnixNano())),
		peerUpdateRequestEvent: event.NewEvent(),
		peerLookupRequestEvent: event.NewEvent(),
	}

	rz.Zone.root = rz

	rz.startChecks()

	return rz
}

// Event fired when the router need update a peer information
func (router *Router) PeerUpdateRequestEvent() *event.Handler {
	return router.peerUpdateRequestEvent.GetHandler()
}

// Event fired when the router need lookup a peer
func (router *Router) PeerLookupRequestEvent() *event.Handler {
	return router.peerLookupRequestEvent.GetHandler()
}

func (router *Router) SaveFile(path string) error {
	return errors.New("not implemented yet")
}

// Get a list of peers of [max] size prepared to do a bootstrap
func (router *Router) GetBootstrapPeers(max int) []*types2.Peer {
	const BootstrapDepth = 5 // Defined as LOG_BASE_EXPONENT constant in protocol/defines.h
	return router.GetTopPeers(max, BootstrapDepth)
}
