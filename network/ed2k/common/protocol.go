package common

type Protocol byte

const (
	ProtocolEd2kTCP            Protocol = 0xe3
	ProtocolEd2kServerUDP      Protocol = 0xe3
	ProtocolEd2kPeerUDP        Protocol = 0xc5
	ProtocolEmuleTcp           Protocol = 0xe3
	ProtocolEmuleTcpCompressed Protocol = 0xe3

	ProtEd2kUSP          = 0xc5
	ProtEd2kUDPServer    = 0xe3
	ProtEd2k2TCP         = 0xf4
	ProtEd2k2UDP         = 0xf5
	ProtEmuleTCP         = 0xc5
	ProtEmuleTCPCompress = 0xd4
	ProtEmuleUDPR1       = 0xa3 /* For encrypted datagrams */
	ProtEmuleUDPR2       = 0xb2
	ProtKadUDP           = 0xe4
	ProtKadUDPCompress   = 0xe5
	ProtocolVersion2     = uint8(2) // eMule 0.47a
	ProtocolVersion3     = uint8(3) // eMule 0.47b
	ProtocolVersion4     = uint8(4) // eMule 0.47c
	ProtocolVersion5     = uint8(5) // eMule 0.48a
	ProtocolVersion6     = uint8(6) // eMule 0.48b
	ProtocolVersion7     = uint8(7) // eMule 0.49a
	ProtocolVersion8     = uint8(8) // eMule 0.49b
)
