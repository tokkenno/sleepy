package factory

import (
	"sleepy/network/kad/common"
	kadPacket "sleepy/network/kad/packet"
	kadTypes "sleepy/network/kad/types"
)

func GetBootstrap1Response(peers []kadTypes.Peer) *kadPacket.Packet {
	packet := kadPacket.NewFixedSizePacket(common.OperationBootstrapResponse, 2+len(peers)*(16+4+2+2+1))
	packet.AppendUInt16(uint16(len(peers)))
	for _, peer := range peers {
		insertContact(packet, peer)
	}
	return packet
}
