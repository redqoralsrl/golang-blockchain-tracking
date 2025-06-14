// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

abstract contract ERC20 {
    function name() public view virtual returns (string memory);
    function symbol() public view virtual returns (string memory);
    function decimals() public view virtual returns (uint8);
    function totalSupply() public view virtual returns (uint256);
    function balanceOf(address account) public view virtual returns (uint256);
    function transfer(address to, uint256 value) public virtual returns (bool);
    function allowance(address owner, address spender) public view virtual returns (uint256);
    function approve(address spender, uint256 value) public virtual returns (bool);
    function transferFrom(address from, address to, uint256 value) public virtual returns (bool);

    event Transfer(address indexed from, address indexed to, uint tokens);
    event Approval(address indexed tokenOwner, address indexed spender, uint tokens);
}