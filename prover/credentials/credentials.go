package credentials

import (
	"context"
	"fmt"
	"math/big"
	"proof-generator/clients"
	"proof-generator/credentials/onchain"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ValidatorWithIndex = struct {
	Validator *phase0.Validator
	Index     uint64
}

/**
 * Generates a .ProveValidatorContainers() proof for all eligible validators on the pod. If `validatorIndex` is set, it will only generate  a proof
 * against that validator, regardless of the validator's state.
 */
func GenerateValidatorProof(ctx context.Context, sympodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient clients.BeaconClient, validatorIndex *big.Int, verbose bool) ([]byte, error) {
	latestBlock, err := eth.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load latest block: %w", err)
	}

	symPod, err := onchain.NewEigenPod(common.HexToAddress(sympodAddress), eth)
	if err != nil {
		return nil, fmt.Errorf("failed to reach sympod: %w", err)
	}

	expectedBlockRoot, err := symPod.GetParentBlockRoot(nil, latestBlock.Time())
	if err != nil {
		return nil, fmt.Errorf("failed to load parent block root: %w", err)
	}

	header, err := beaconClient.GetBeaconHeader(ctx, "0x"+common.Bytes2Hex(expectedBlockRoot[:]))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon header: %w", err)
	}

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}
	if beaconState != nil {
		fmt.Println("Beacon state fetched!!!")
	}

	allValidators, err := FindAllValidatorsForSympod(sympodAddress, beaconState)
	if err != nil {
		return nil, fmt.Errorf("failed to find validators: %w", err)
	}

	if len(allValidators) == 0 {
		return nil, fmt.Errorf("no validators found")
	}
	// Get withdrawal credentials of first validator
	fmt.Println("Validator details:", allValidators[0].Validator)
	withdrawalCredentials := allValidators[0].Validator.WithdrawalCredentials
	return withdrawalCredentials, nil

	// proofExecutor, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	// if err != nil {
	// 	return nil, 0, fmt.Errorf("failed to initialize provider: %w", err)
	// }

	// proofs, err := GenerateValidatorProofAtState(ctx, proofExecutor, sympodAddress, beaconState, eth, chainId, header, latestBlock.Time(), validatorIndex, verbose)
	// return proofs, latestBlock.Time(), err
}

// search through beacon state for validators whose withdrawal address is set to eigenpod.
func FindAllValidatorsForSympod(sympodAddress string, beaconState *spec.VersionedBeaconState) ([]ValidatorWithIndex, error) {
	allValidators, err := beaconState.Validators()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}

	eigenpod := common.HexToAddress(sympodAddress)

	var outputValidators []ValidatorWithIndex = []ValidatorWithIndex{}
	var i uint64 = 0
	maxValidators := uint64(len(allValidators))
	for i = 0; i < maxValidators; i++ {
		validator := allValidators[i]
		if validator == nil || validator.WithdrawalCredentials[0] != 1 { // withdrawalCredentials _need_ their first byte set to 1 to withdraw to execution layer.
			continue
		}
		// we check that the last 20 bytes of expectedCredentials matches validatorCredentials.
		// // first 12 bytes are not the pubKeyHash, see (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/pods/EigenPod.sol#L663)
		validatorWithdrawalAddress := common.BytesToAddress(validator.WithdrawalCredentials[12:])

		if eigenpod.Cmp(validatorWithdrawalAddress) == 0 {
			outputValidators = append(outputValidators, ValidatorWithIndex{
				Validator: validator,
				Index:     i,
			})
		}
	}
	return outputValidators, nil
}
