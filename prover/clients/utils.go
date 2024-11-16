package clients

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetClients(ctx context.Context, node, beaconNodeUri string, enableLogs bool) (*ethclient.Client, BeaconClient, *big.Int, error) {
	eth, err := ethclient.Dial(node)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reach eth --node: %w", err)
	}

	chainId, err := eth.ChainID(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to fetch chain id: %w", err)
	}

	if chainId == nil || (chainId.Int64() != 17000 && chainId.Int64() != 1) {
		return nil, nil, nil, errors.New("this tool only supports the Holesky and Mainnet Ethereum Networks")
	}

	beaconClient, err := GetBeaconClient(beaconNodeUri, enableLogs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reach beacon client: %w", err)
	}

	genesisForkVersion, err := beaconClient.GetGenesisForkVersion(ctx)
	expectedForkVersion := ForkVersions()[chainId.Uint64()]
	gotForkVersion := hex.EncodeToString((*genesisForkVersion)[:])
	if err != nil || expectedForkVersion != gotForkVersion {
		return nil, nil, nil, fmt.Errorf("check that both nodes correspond to the same network and try again (expected genesis_fork_version: %s, got %s)", expectedForkVersion, gotForkVersion)
	}

	return eth, beaconClient, chainId, nil
}

func GetBeaconClient(beaconUri string, verbose bool) (BeaconClient, error) {
	beaconClient, _, err := NewBeaconClient(beaconUri, verbose)
	return beaconClient, err
}

// this is a mapping from <chainId, genesis_fork_version>.
func ForkVersions() map[uint64]string {
	return map[uint64]string{
		11155111: "90000069", //sepolia (https://github.com/eth-clients/sepolia/blob/main/README.md?plain=1#L66C26-L66C36)
		17000:    "01017000", //holesky (https://github.com/eth-clients/holesky/blob/main/README.md)
		1:        "00000000", // mainnet (https://github.com/eth-clients/mainnet)
	}
}
