package main

import (
	"flag"
	"fmt"
	"os"

	"proof-generator/circuits"

	"github.com/brevis-network/brevis-sdk/sdk/prover"
)

var port = flag.Uint("port", 33247, "the port to start the service at")

func main() {
	flag.Parse()

	proverService, err := prover.NewService(&circuits.AppCircuit{}, prover.ServiceConfig{
		SetupDir: "$HOME/circuitOut",
		SrsDir:   "$HOME/kzgsrs",
		RpcURL:   "https://sepolia.infura.io/v3/81bcd90d66aa4944a44ef530fb36f329",
		ChainId:  11155111,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proverService.Serve("", *port)
}
