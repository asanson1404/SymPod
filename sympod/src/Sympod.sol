// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SymbioticVaultMock} from "./SymbioticVaultMock.sol";
import {BrevisApp} from "../lib/brevis-contracts/contracts/sdk/apps/framework/BrevisApp.sol";
import {ERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "../lib/openzeppelin-contracts/contracts/access/Ownable.sol";

contract Sympod is BrevisApp, ERC20, Ownable {

    /* ============== EVENTS ============== */

    event Restaked(address indexed staker, address indexed symbioticVaultAddr);
    event VerificationRequired(bytes indexed pubkey, bytes32 depositDataRoot, address sympodAddr);

    /* ============== STORAGE ============== */

    // Track the amount a staker has restaked
    mapping (address => uint256) public stakerToStakedAmount;

    // Restaking strategy
    address public symbioticVault;

    // VK hash of the Brevis circuit
    bytes32 public vkHash;

    /* ============== CONSTRUCTOR ============== */

    constructor(address _brevisRequest) BrevisApp(_brevisRequest) ERC20("SPETH", "SYMPODETH") Ownable(msg.sender) {}

    /* ============== EXTERNAL FUNCTIONS ============== */

    function restake(address symbioticVaultAddr) external payable {
        require(msg.value == 32 ether, "can only restake 32 ETH");
        stakerToStakedAmount[msg.sender] += msg.value;
        emit Restaked(msg.sender, symbioticVaultAddr);
    }

    function verifyWithdrawalCredential(
        bytes calldata pubkey,
        bytes32 depositDataRoot
    ) external {
        emit VerificationRequired(pubkey, depositDataRoot, address(this));
    }

    function setVkHash(bytes32 _vkHash) external onlyOwner {
        vkHash = _vkHash;
    }

    function setSymbioticVault(address _symbioticVault) external onlyOwner {
        symbioticVault = _symbioticVault;
    }

    /* ============== BREVIS FUNCTIONS ============== */

    // BrevisQuery contract will call our callback once Brevis backend submits the proof.
    // This method is called with once the proof is verified.
    function handleProofResult(bytes32 _vkHash, bytes calldata _circuitOutput) internal override {
        // We need to check if the verifying key that Brevis used to verify the proof
        // generated by our circuit is indeed our designated verifying key. This proves
        // that the _circuitOutput is authentic
        require(vkHash == _vkHash, "invalid vk");
        address withdrawalAddr = decodeOutput(_circuitOutput);
        require(withdrawalAddr == address(this), "invalid withdrawal address");
        _mint(address(this), 32 ether);
        approve(symbioticVault, 32 ether);
        SymbioticVaultMock(symbioticVault).deposit(32 ether);
    }

    /* ============== INTERNAL FUNCTIONS ============== */

    function decodeOutput(bytes calldata o) internal pure returns (address) {
        // Extract the last 20 bytes which contain the smart contract address
        address withdrawalAddr = address(bytes20(o[12:32]));
        return withdrawalAddr;
    }

}