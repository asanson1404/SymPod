// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"proof-generator/circuits"
// 	"proof-generator/clients"

// 	"github.com/fatih/color"
// 	"github.com/joho/godotenv"
// )

// // func main() {

// 	// Context
// 	ctx := context.Background()

// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	// Get the environment variables
// 	beaconNodeUrl := os.Getenv("BEACON_NODE_URL")
// 	ethNodeUrl := os.Getenv("ETH_NODE_URL")

// 	// Retrieve the sympod address
// 	if len(os.Args) != 2 {
// 		log.Fatal("Usage: go run main.go <symPodAddress>")
// 	}
// 	sympodAddr := os.Args[1]

// 	// Get ethereum and beacon clients
// 	eth, beaconClient, chainId, err := clients.GetClients(ctx, ethNodeUrl, beaconNodeUrl, true)
// 	if err != nil {
// 		log.Fatal("Error fetching clients:", err)
// 	}

// 	// Generate the proof for all the validators on the pod
// 	validatorProofs, oracleBeaconTimestamp, err := circuits.GenerateValidatorProof(ctx, sympodAddr, eth, chainId, beaconClient, nil, true)
// 	if err != nil || validatorProofs == nil {
// 		PanicOnError("Failed to generate validator proof", err)
// 	}
// 	fmt.Printf("Validator proofs: %+v\n", validatorProofs)
// 	fmt.Printf("Oracle beacon timestamp: %d\n", oracleBeaconTimestamp)
// }

// func PanicOnError(message string, err error) {
// 	if err != nil {
// 		color.Red(fmt.Sprintf("error: %s\n\n", message))

// 		info := color.New(color.FgRed, color.Italic)
// 		info.Printf(fmt.Sprintf("caused by: %s\n", err))

// 		os.Exit(1)
// 	}
// }

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
		// RpcURL:   "https://rpc.ankr.com/eth_holesky",
		// ChainId:  17000,
		// RpcURL:  "https://sepolia.infura.io/v3/81bcd90d66aa4944a44ef530fb36f329",
		// ChainId: 11155111,
		RpcURL:  "https://eth.llamarpc.com",
		ChainId: 1,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proverService.Serve("", *port)
}
