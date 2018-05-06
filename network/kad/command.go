package kad

const (
	CommUnknown                 = 0xff

	// TODO: Emule commands

	CommKadBootstrapReq         = 0x00
	CommKadBootstrapRes         = 0x08

	CommKadHelloReq             = 0x10
	CommKadHelloRes             = 0x18

	CommKadFirewalledReq        = 0x50
	CommKadFirewalled2Req       = 0x53
	CommKadFirewalledRes        = 0x58

	CommKadCallbackReq          = 0x52

	CommKadReq                  = 0x20
	CommKadRes                  = 0x28

	CommKadPublishReq           = 0x40
	CommKadPublishRes           = 0x48

	CommKadSearchReq            = 0x30
	CommKadSearchRes            = 0x38

	CommKadSearchNotesReq       = 0x32
	CommKadSearchNotesRes       = 0x3A

	CommKadFindbuddyReq         = 0x51
	CommKadFindbuddyRes         = 0x5A

	CommKadPublishNotesReq      = 0x42
	CommKadPublishNotesRes      = 0x4A

	CommKad2BootstrapReq        = 0x01
	CommKad2BootstrapRes        = 0x09

	CommKad2Req                 = 0x21
	CommKad2Res                 = 0x29

	CommKad2HelloReq            = 0x11
	CommKad2HelloRes            = 0x19

	CommKad2HelloResAck         = 0x22

	CommKad2FirewallUDP         = 0x62

	CommKad2SearchKeyReq        = 0x33
	CommKad2SearchSourceReq     = 0x34
	CommKad2SearchNotesReq      = 0x35

	CommKad2SearchRes           = 0x3B

	CommKad2PublichKeyReq       = 0x43
	CommKad2PublishSourceReq    = 0x44
	CommKad2PublishNotesReq     = 0x45

	CommKad2PublishRes          = 0x4B

	CommKad2PublishResAck       = 0x4C

	CommKad2Ping                = 0x60
	CommKad2Pong                = 0x61
)
