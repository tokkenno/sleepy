package ed2k

const (
	ProtEd2kTCP             = 0xe3
	ProtEd2kUSP             = 0xc5
	ProtEd2kUDPServer       = 0xe3

	ProtEd2k2TCP            = 0xf4
	ProtEd2k2UDP            = 0xf5

	ProtEmuleTCP            = 0xc5
	ProtEmuleTCPCompress    = 0xd4

	/* For encrypted datagrams */
	ProtEmuleUDPR1          = 0xa3
	ProtEmuleUDPR2          = 0xb2

	ProtKadUDP              = 0xe4
	ProtKadUDPCompress      = 0xe5

	MlDonkey                = 0x00
)