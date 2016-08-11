package main

import (
	"bufio"
	"os"

	"github.com/phito/gofus/dofus"
	"github.com/phito/gofus/network"
)

func main() {
	const executable = "/home/romain/programs/Dofus/bin/Dofus"
	const fingerprint = "./fingerprint"
	const payload = "./payload"

	dofus.RunClient(executable, fingerprint, payload)

	server := network.Server{}
	server.Open(":5555")

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
