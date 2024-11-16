// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "../lib/openzeppelin-contracts/contracts/access/Ownable.sol";


contract SymbioticVaultMock is Ownable {

    // speth being deposited
    ERC20 public speth;

    // Total deposits
    uint256 public totalDeposits;

    constructor(address _speth) Ownable(msg.sender) {
        speth = ERC20(_speth);
    }

    function deposit(uint256 amount) external {
        require(amount > 0, "Amount must be greater than 0");
        require(speth.transferFrom(msg.sender, address(this), amount), "Transfer failed");
        
        totalDeposits += amount;
    }

    function setSpeth(address _speth) external onlyOwner {
        speth = ERC20(_speth);
    }

}
