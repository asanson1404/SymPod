package circuits

import (
	"context"
	"fmt"
	"log"
	"os"
	"proof-generator/clients"
	"proof-generator/credentials"

	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

// Context
var ctx = context.Background()

// Load environment variables from .env file
var _ = godotenv.Load()

// Get the environment variables
var beaconNodeUrl = os.Getenv("BEACON_NODE_URL")
var ethNodeUrl = os.Getenv("ETH_NODE_URL")

// Get ethereum and beacon clients
var eth, beaconClient, chainId, _ = clients.GetClients(ctx, ethNodeUrl, beaconNodeUrl, true)

var symPod = sdk.ConstUint248("0xEF1E121CB094023a453552B6226714ef275EA193")

type AppCircuit struct{}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxSlots, maxTransactions int) {
	return 32, 0, 0
}

func (c *AppCircuit) Define(api *sdk.CircuitAPI, in sdk.DataInput) error {

	receipts := sdk.NewDataStream(api, in.Receipts)
	receipt := sdk.GetUnderlying(receipts, 0)

	fmt.Println("SymPod Address:", receipt.Fields[0].Contract)

	// Fetch validator credentials
	credentials, err := credentials.GenerateValidatorProof(ctx, "0x21E2a892DDc9BD3c0466299172F8b1D8026925ED", eth, chainId, beaconClient, nil, true)
	if err != nil {
		log.Fatal("Error generating validator proof:", err)
	}
	fmt.Printf("Credentials: %+v\n", credentials)
	withdrawalAddr := common.BytesToAddress(credentials[12:])
	fmt.Printf("withdrawalAddr: %+v\n", withdrawalAddr)

	// Check logic
	// The first field exports `symPod addr` parameter from the Event
	api.Uint248.AssertIsEqual(receipt.Fields[0].Contract, symPod)
	api.Uint248.AssertIsEqual(receipt.Fields[0].IsTopic, sdk.ConstUint248(1))
	api.Uint248.AssertIsEqual(receipt.Fields[0].Index, sdk.ConstUint248(1))
	api.Uint32.AssertIsEqual(receipt.Fields[0].LogPos, receipt.Fields[1].LogPos)

	// Check that the withdrawal address is correct
	api.Uint248.AssertIsEqual(sdk.ConstUint248("0x21E2a892DDc9BD3c0466299172F8b1D8026925ED"), sdk.ConstUint248(withdrawalAddr))

	// // The second field exports `pubKey` of the validator
	// api.Uint248.AssertIsEqual(receipt.Fields[1].IsTopic, sdk.ConstUint248(0))
	// api.Uint248.AssertIsEqual(receipt.Fields[1].Index, sdk.ConstUint248(0))

	// // Make sure this transfer has minimum 500 USDC volume
	// api.Uint248.AssertIsLessOrEqual(minimumVolume, api.ToUint248(receipt.Fields[1].Value))

	// Output the withdrawal credentials
	api.OutputBytes32(sdk.ConstFromBigEndianBytes(credentials))

	return nil
}
