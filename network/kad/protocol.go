package kad

import (
	"log"
)

func HandleBootstrapRequest(client *Client, r *UDPRequest, w Response) {
	log.Println("Bootstrap request")

	contacts := client.router.GetBootstrapPeers(20)
	log.Println("Se van a enviar " + string(len(contacts)) + " contactos.")
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
