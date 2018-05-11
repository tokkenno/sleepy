package kad

import (
	"fmt"
	"os"
	"net"
	"github.com/reimashi/sleepy/types"
)

type Client struct {
	port uint16
	server *Server
	clientAddr *net.UDPAddr
	clientConn *net.UDPConn
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
	client.server = newServer(port, client)
	return client
}

func (this *Client) Start () {
	this.server.Start()
}

func (this *Client) Stop () {
	this.server.Stop()

	this.clientConn.Close()
}

func (this *Client) SendPacket(command byte, body []byte, to *net.UDPAddr, cryptId *types.UInt128) {
	_, err := this.server.serverConn.WriteToUDP(body, to)
	fmt.Println(err)
}

func (this *Client) bootstrap(peer Peer) {

}

func (this *Client) sendMyDetails() {

}

func (this *Client) SendPong(peer *Peer) error {
	return nil
}