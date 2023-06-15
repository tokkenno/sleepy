package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sleepy/network"
	"sleepy/network/kad"
	"sleepy/types"
)

func main() {
	networkManager := network.NewManager()

	kadClient := kad.NewClient(kad.Config{
		UdpPort:  4662,
		ClientID: types.NewUInt128(rand.Uint64(), rand.Uint64()),
	}, networkManager)
	kadClient.Start()

	fmt.Println("Listening KAD")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Println("Closing KAD")

	kadClient.Stop()
}
