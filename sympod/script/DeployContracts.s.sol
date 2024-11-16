// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Sympod} from "../src/Sympod.sol";
import {SymbioticVaultMock} from "../src/SymbioticVaultMock.sol";

// forge script script/DeployContracts.s.sol --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY --broadcast --etherscan-api-key $ETHERSCAN_API_KEY --verify -vvv
contract DeployContracts is Script {
    address constant BREVIS_REQUEST = 0xa082F86d9d1660C29cf3f962A31d7D20E367154F;

    function run() public {
        vm.startBroadcast();
        Sympod sympod = new Sympod(BREVIS_REQUEST);
        SymbioticVaultMock symbioticVault = new SymbioticVaultMock(address(sympod));
        sympod.setSymbioticVault(address(symbioticVault));
        console.log("Sympod deployed at:", address(sympod));
        console.log("SymbioticVault deployed at:", address(symbioticVault));
        console.log("BrevisRequest address:", BREVIS_REQUEST);
        vm.stopBroadcast();
    }
}