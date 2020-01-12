package kad

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"sleepy/network/ed2k"
	"sleepy/network/kad/router"
	"strconv"
	"time"
)

type Client struct {
	router     *router.Router
	listenPort uint16
	clientAddr *net.UDPAddr
	clientConn *net.UDPConn
	serverAddr *net.UDPAddr
	serverConn *net.UDPConn
}

func NewClient(port uint16) *Client {
	client := new(Client)
	client.listenPort = port
	return client
}

func (client *Client) Start() error {
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(int(client.listenPort)))
	if err != nil {
		return err
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}

	client.serverAddr = serverAddr
	client.serverConn = serverConn

	go client.listenUDP()
	return nil
}

func (client *Client) Stop() {
	client.serverConn.SetDeadline(time.Now())
	client.clientConn.Close()
}

func (client *Client) listenUDP() {
	defer client.serverConn.Close()

	buf := make([]byte, 8192)

	for {
		n, addr, err := client.serverConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println(err)
		} else {
			go func() {
				err := client.handleUDP(buf[0:n], addr)
				if err != nil {
					log.Printf("Datagram handle error: %s", err)
				}
			}()
		}
	}
}

func (client *Client) handleUDP(data []byte, from *net.UDPAddr) error {
	if from.Port == 53 {
		return errors.New("Dropping incoming ping from port 53. Possible DNS attack.")
	}

	fmt.Println("Received ", hex.EncodeToString(data), "text", string(data), " from ", from)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = context.WithValue(ctx, "Protocol", "UDP")
	ctx = context.WithValue(ctx, "IP", from.IP)
	ctx = context.WithValue(ctx, "Port", from.Port)

	request := &UDPRequest{
		from: from,
		Request: Request{
			body: Reader{data: data, offset: 0},
			ctx:  ctx,
		},
	}

	protocolCode, err := request.body.ReadByte()
	fmt.Printf("Protocol: %s, error: %s", hex.EncodeToString([]byte{protocolCode}), err)
	if err != nil {
		return errors.New("datagram read error")
	}

	switch protocolCode {
	case ed2k.ProtKadUDPCompress:
		log.Println("Compressed KAD datagram. Trying decompression...")
		return client.decompressKad(data, from)
	case ed2k.ProtKadUDP:
		log.Println("Handling valid Kad UDP packet...")
		return client.handleKadDatagram(request)
	default:
		return errors.New("unknown packet " + hex.EncodeToString([]byte{protocolCode}) + " to parse")
	}
}

func (client *Client) decompressKad(data []byte, from *net.UDPAddr) error {
	return errors.New("decompress kad not implemented yet")
}

func (client *Client) handleKadDatagram(request *UDPRequest) error {
	command, err := request.body.ReadByte()

	if err != nil {
		return errors.New("datagram read error")
	}

	response := Response{}

	switch command {
	case CommKad2BootstrapReq:
		HandleBootstrapRequest(client, request, response)
		return nil
	case CommKad2BootstrapRes:
		HandleBootstrapResponse(client, request, response)
		return nil
	case CommKad2HelloReq:
		HandleHelloRequest(client, request, response)
		return nil
	case CommKad2HelloRes:
		HandleHelloResponse(client, request, response)
		return nil
	case CommKad2HelloResAck:
		HandleHelloResponseAck(client, request, response)
		return nil
	case CommKadFirewalled2Req:
		HandleFirewallRequest(client, request, response)
		return nil
	case CommKad2Ping:
		HandlePingRequest(client, request, response)
		return nil
	case CommKad2Pong:
		HandlePongResponse(client, request, response)
		return nil
	default:
		return errors.New("unknown kad command")
	}
}
