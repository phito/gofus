package main

import (
	"bufio"
	"os"

	"github.com/phito/gofus/dofus"
)

func main() {

	const executable = "/home/romain/programs/Dofus/bin/Dofus"
	const fingerprint = "./fingerprint"
	const payload = "./payload"

	_, err := dofus.RunClient(executable, fingerprint, payload)

	if err != nil {
		println("RunClient failed:", err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	/*var conn network.Connection
	if err := conn.Open("213.248.126.39:5555"); err != nil {
		println("Open failed: ", err)
		os.Exit(1)
	}

	println("connection opened")

	conn.Close()

	println("connection closed")*/
}
