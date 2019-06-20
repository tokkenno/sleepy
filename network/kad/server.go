package kad

import (
	"errors"
	"fmt"
	"github.com/tokkenno/sleepy/io"
	"github.com/tokkenno/sleepy/network/ed2k"
	"net"
	"strconv"
	"time"
)

type Server struct {
	port       uint16
	client     *Client
	serverAddr *net.UDPAddr
	serverConn *net.UDPConn
	stopListen bool
}

func newServer(port uint16, client *Client) *Server {
	server := new(Server)
	server.port = port
	return server
}

func (server *Server) Start() {
	server.stopListen = false

	serverAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(int(server.port)))
	checkError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)

	server.serverAddr = serverAddr
	server.serverConn = serverConn

	go server.listen()
}

// Listen new UDP connections
func (server *Server) listen() {

	defer server.serverConn.Close()

	buf := make([]byte, 8192)

	for {
		n, addr, err := server.serverConn.ReadFromUDP(buf)

		if server.stopListen {
			return
		}

		if err != nil {
			fmt.Println(err)
		} else {
			go server.handleDatagram(buf[0:n], addr)
		}
	}
}

func (server *Server) Stop() {
	server.stopListen = true
	server.serverConn.SetDeadline(time.Now())
}

func (server *Server) handleDatagram(data []byte, from *net.UDPAddr) error {
	if from.Port == 53 {
		return errors.New("Dropping incoming ping from port 53. Possible DNS attack.")
	}

	fmt.Println("Received ", string(data), " from ", from)

	dataReader := io.NewReader(data)

	protocolCode := dataReader.ReadByte()
	size := dataReader.ReadInt()

	if !dataReader.Correct() {
		return errors.New("datagram read error")
	} else if size+5 == len(data) {
		datagram := InDatagram{Datagram{protocolCode, &data}, dataReader, from}

		switch datagram.protocolCode {
		case ed2k.ProtKadUDPCompress:
			return server.decompressKad(data, from)
		case ed2k.ProtKadUDP:
			return server.handleKadDatagram(&datagram)
		}

		return errors.New("unknown packet to parse")
	} else {
		return errors.New("datagram size mismatch")
	}
}

func (server *Server) decompressKad(data []byte, from *net.UDPAddr) error {
	return errors.New("decompress kad not implemented yet")
}

func (server *Server) handleKadDatagram(datagram *InDatagram) error {
	command := datagram.reader.ReadByte()

	kadDatagram := KadInDatagram{datagram, command}

	switch kadDatagram.command {
	case CommKad2BootstrapReq:
		return server.handleKad2BootstrapReq(&kadDatagram)
	case CommKad2BootstrapRes:
		return server.handleKad2BootstrapRes(&kadDatagram)
	case CommKad2HelloReq:
		return server.handleKad2HelloReq(&kadDatagram)
	case CommKad2HelloRes:
		return server.handleKad2HelloRes(&kadDatagram)
	case CommKad2HelloResAck:
		return server.handleKad2HelloResAck(&kadDatagram)
	case CommKad2Req:
		return server.handleKad2Req(&kadDatagram)
	case CommKad2Res:
		return server.handleKad2Res(&kadDatagram)
	case CommKadFirewalled2Req:
		return server.handleKadFirewalled2Req(&kadDatagram)
	case CommKad2FirewallUDP:
		return server.handleKad2FirewallUdp(&kadDatagram)
	case CommKad2Ping:
		return server.handleKad2Ping(&kadDatagram)
	case CommKad2Pong:
		return server.handleKad2Pong(&kadDatagram)
	}

	return errors.New("unknown kad packet to parse")
}

func (server *Server) handleKad2BootstrapReq(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2BootstrapRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2HelloReq(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2HelloRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2HelloResAck(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2Req(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2Res(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKadFirewalledReq(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKadFirewalled2Req(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKadFirewalledRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKadFirewalledAckRes(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2FirewallUdp(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2Ping(datagram *KadInDatagram) error {
	//server.client.SendPong()
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}

func (server *Server) handleKad2Pong(datagram *KadInDatagram) error {
	return errors.New("Not implemented exception " + strconv.Itoa(int(datagram.command)))
}
