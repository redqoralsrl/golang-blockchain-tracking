// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

abstract contract Contract {
    function name() public view virtual returns (string memory);
    function symbol() public view virtual returns (string memory);

    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner);
}