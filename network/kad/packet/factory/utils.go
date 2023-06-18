package factory

import (
	"sleepy/network/kad/packet"
	kadTypes "sleepy/network/kad/types"
)

func insertContact(packet *packet.Packet, peer kadTypes.Peer) {
	packet.AppendBytes(peer.GetID().ToBytes())
	packet.AppendBytes(peer.GetIP().To4())
	packet.AppendUInt16(peer.GetUDPPort())
	packet.AppendUInt16(peer.GetTCPPort())
	packet.AppendUInt8(peer.GetProtocolVersion())
}
