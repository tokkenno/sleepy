package kad

import (
	"fmt"
	"log"
	"sleepy/network/kad/packet/factory"
)

func HandleBootstrapRequest(client *Client, r *UDPRequest) {
	// Some clients send the remote ip and port, others don't
	remoteIp := r.from.IP
	remoteUdpPort := uint16(r.from.Port)

	remoteId, err := r.body.ReadUInt128()
	if err == nil {
		remoteIp, err = r.body.ReadIPv4()
		if err != nil {
			log.Fatalf("Error al leer la ip remota: %s", err.Error())
		}
		remoteUdpPort, err = r.body.ReadUInt16()
		if err != nil {
			log.Fatalf("Error al leer el puerto udp remoto: %s", err.Error())
		}
	}

	log.Println("Bootstrap request")
	contacts := client.router.GetBootstrapPeers(20, remoteId)
	log.Println(fmt.Sprintf("Se van a enviar %d contactos.", len(contacts)))
	packet := factory.GetBootstrap1Response(contacts)
	err = client.network.SendUDP(remoteIp, remoteUdpPort, packet)
	if err != nil {
		log.Println(err)
	}
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
