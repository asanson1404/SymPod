// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {Sympod} from "../src/Sympod.sol";

contract SympodTest is Test {
    Sympod public sympod;

    // brevisRequest address Sepolia
    address brevisRequest = address(0xa082F86d9d1660C29cf3f962A31d7D20E367154F);

    function setUp() public {
        sympod = new Sympod(brevisRequest);
    }
}
