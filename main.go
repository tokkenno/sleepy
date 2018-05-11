package main

import (
	"github.com/reimashi/sleepy/network/kad"
	"bufio"
	"os"
	"fmt"
)

func main() {
	kadClient := kad.NewClient(4662)
	kadClient.Start()

	fmt.Println("Listening KAD")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Println("Closing KAD")

	kadClient.Stop()
}
