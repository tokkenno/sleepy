package kad

import (
	"fmt"
	"log"
	"sleepy/network/kad/packet/factory"
)

func HandleBootstrapRequest(client *Client, r *UDPRequest) {
	log.Println("Bootstrap request")
	contacts := client.router.GetBootstrapPeers(20)
	log.Println(fmt.Sprintf("Se van a enviar %d contactos.", len(contacts)))
	packet := factory.GetBootstrap1Response(contacts)
	client.network.SendUDP(r.from.IP, uint16(r.from.Port), packet)
}

func HandleBootstrapResponse(client *Client, r *UDPRequest, w Response) {
	log.Println("Bootstrap response")
}

func HandleFirewallRequest(client *Client, r *UDPRequest, w Response) {
	log.Println("Firewall request")
}

func HandleHelloRequest(client *Client, r *UDPRequest, w Response) {
	log.Println("Hello request")
}

func HandleHelloResponse(client *Client, r *UDPRequest, w Response) {
	log.Println("Hello response")
}

func HandleHelloResponseAck(client *Client, r *UDPRequest, w Response) {
	log.Println("Hello response ack")
}

func HandlePingRequest(client *Client, r *UDPRequest, w Response) {
	log.Println("Ping request")
}

func HandlePongResponse(client *Client, r *UDPRequest, w Response) {
	log.Println("Pong response")
}
