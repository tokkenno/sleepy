package kad

import (
	"net"
	"fmt"
	"os"
	"strconv"
	"errors"
	"bytes"
	"com/github/reimashi/sleepy/network/kad/ed2k"
	"time"
)

type Client struct {
	port uint16
	server *net.UDPConn
	stopListen bool
}

func checkError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}
}

func NewClient(port uint16) *Client {
	client := new(Client)
	client.port = port
	return client
}

func (this *Client) Start () {
	this.stopListen = false
	go this.listen()
}

/**
 * Listen new UDP connections
 */
func (this *Client) listen() {
	udpAddress, err := net.ResolveUDPAddr("udp4", ":" + strconv.Itoa(int(this.port)))
	checkError(err)

	serverConn, err := net.ListenUDP("udp", udpAddress)
	checkError(err)

	this.server = serverConn
	defer serverConn.Close()

	buf := make([]byte, 8192)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)

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

func (this *Client) Stop() {
	this.stopListen = true
	this.server.SetDeadline(time.Now())
}

func (this *Client) handleDatagram(data []byte, from *net.UDPAddr) error {
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
		return decompressKad(data, from)
		break
	case ed2k.ProtKadUDP:
		return handleKadDatagram(&datagram)
	}

	return errors.New("Unknow packet to parse")
}

func decompressKad(data []byte, from *net.UDPAddr) error {
	return errors.New("uncompressKad not implemented yet")
}

func handleKadDatagram(datagram *InDatagram) error {
	command, err := datagram.reader.ReadByte()
	if err != nil { return err }

	kadDatagram := KadInDatagram{datagram, command}

	switch kadDatagram.command {
	case ed2k.CommKad2Ping:

		break
	}

	return errors.New("Unknow KAD packet to parse")
}

func handleKadPing(datagram *KadInDatagram) error {
	return nil
}