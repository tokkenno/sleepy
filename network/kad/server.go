package kad

import (
	"errors"
	"strconv"
	"net"
	"fmt"
	"time"
	"bytes"
	"com/github/reimashi/sleepy/network/kad/ed2k"
)

type Server struct {
	port uint16
	client *Client
	serverAddr *net.UDPAddr
	serverConn *net.UDPConn
	stopListen bool
}

func newServer(port uint16, client *Client) *Server {
	server := new(Server)
	server.port = port
	return server
}

func (this *Server) Start () {
	this.stopListen = false

	serverAddr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(int(this.port)))
	checkError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)

	this.serverAddr = serverAddr
	this.serverConn = serverConn

	go this.listen()
}

/**
 * Listen new UDP connections
 */
func (this *Server) listen() {

	defer this.serverConn.Close()

	buf := make([]byte, 8192)

	for {
		n, addr, err := this.serverConn.ReadFromUDP(buf)

		if (this.stopListen) {
			return
		}

		if err != nil {
			fmt.Println(err)
		} else {
			go this.handleDatagram(buf[0:n], addr)
		}
	}
}

func (this *Server) Stop() {
	this.stopListen = true
	this.serverConn.SetDeadline(time.Now())
}

func (this *Server) handleDatagram(data []byte, from *net.UDPAddr) error {
	if from.Port == 53 {
		return errors.New("Dropping incoming ping from port 53. Possible DNS attack.")
	}

	fmt.Println("Received ",string(data), " from ", from)

	dataReader := bytes.NewReader(data)

	typeCode, err := dataReader.ReadByte()
	if err != nil { return err }
	typeCode = typeCode

	datagram := InDatagram{Datagram{typeCode, data}, dataReader, from}

	switch datagram.typeCode {
	case ed2k.ProtKadUDPCompress:
		return this.decompressKad(data, from)
		break
	case ed2k.ProtKadUDP:
		return this.handleKadDatagram(&datagram)
	}

	return errors.New("Unknow packet to parse")
}

func (this *Server) decompressKad(data []byte, from *net.UDPAddr) error {
	return errors.New("uncompressKad not implemented yet")
}

func (this *Server) handleKadDatagram(datagram *InDatagram) error {
	command, err := datagram.reader.ReadByte()
	if err != nil { return err }

	kadDatagram := KadInDatagram{datagram, command}

	switch kadDatagram.command {
	case ed2k.CommKad2Ping:
		this.handleKad2Ping(&kadDatagram)
		break
	}

	return errors.New("Unknow KAD packet to parse")
}

func (this *Server) handleKad2BootstrapReq(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2BootstrapRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2HelloReq(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2HelloRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2HelloResAck(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2Req(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2Res(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKadFirewalledReq(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKadFirewalled2Req(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKadFirewalledRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKadFirewalledAckRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2Ping(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2Pong(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (this *Server) handleKad2FirewallUdp(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}