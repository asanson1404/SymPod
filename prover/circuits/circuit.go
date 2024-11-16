package circuits

import (
	//"context"
	"fmt"
	"os"

	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/joho/godotenv"
)

// Context
// var ctx = context.Background()

// Load environment variables from .env file
var _ = godotenv.Load()

// Get the environment variables
var beaconNodeUrl = os.Getenv("BEACON_NODE_URL")
var ethNodeUrl = os.Getenv("ETH_NODE_URL")

var USDCTokenAddr = sdk.ConstUint248("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
var minimumVolume = sdk.ConstUint248(500000000) // minimum 500 USDC

type AppCircuit struct {
	//SymPodAddr sdk.Uint248
}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxSlots, maxTransactions int) {
	return 32, 0, 0
}

func (c *AppCircuit) Define(api *sdk.CircuitAPI, in sdk.DataInput) error {
	fmt.Println("Define")
	fmt.Println("Beacon Node URL:", beaconNodeUrl)
	fmt.Println("Eth Node URL:", ethNodeUrl)
	fmt.Println("USDCTokenAddr:", USDCTokenAddr)

	receipts := sdk.NewDataStream(api, in.Receipts)
	receipt := sdk.GetUnderlying(receipts, 0)

	// Check logic
	// The first field exports `from` parameter from Transfer Event
	// It should use the second topic in Transfer Event log
	api.Uint248.AssertIsEqual(receipt.Fields[0].Contract, USDCTokenAddr)
	api.Uint248.AssertIsEqual(receipt.Fields[0].IsTopic, sdk.ConstUint248(1))
	api.Uint248.AssertIsEqual(receipt.Fields[0].Index, sdk.ConstUint248(1))

	// Make sure two fields uses the same log to make sure account address linking with correct volume
	api.Uint32.AssertIsEqual(receipt.Fields[0].LogPos, receipt.Fields[1].LogPos)

	// The second field exports `Volume` parameter from Transfer Event
	// It should use Data in Transfer Event log
	api.Uint248.AssertIsEqual(receipt.Fields[1].IsTopic, sdk.ConstUint248(0))
	api.Uint248.AssertIsEqual(receipt.Fields[1].Index, sdk.ConstUint248(0))

	// Make sure this transfer has minimum 500 USDC volume
	api.Uint248.AssertIsLessOrEqual(minimumVolume, api.ToUint248(receipt.Fields[1].Value))

	api.OutputUint(64, api.ToUint248(receipt.BlockNum))
	api.OutputAddress(api.ToUint248(receipt.Fields[0].Value))
	api.OutputBytes32(receipt.Fields[1].Value)
	fmt.Println("End of Defiiiiiine")
	return nil
}
