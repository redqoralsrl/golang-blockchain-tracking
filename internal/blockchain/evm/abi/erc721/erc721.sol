// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

abstract contract ERC721 {
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
    event ApprovalForAll(address indexed owner, address indexed operator, bool approved);

    function balanceOf(address owner) public view virtual returns (uint256);
    function ownerOf(uint256 tokenId) public view virtual returns (address);
    function name() public view virtual returns (string memory);
    function symbol() public view virtual returns (string memory);
    function tokenURI(uint256 tokenId) public view virtual returns (string memory);
    function approve(address to, uint256 tokenId) public virtual;
    function getApproved(uint256 tokenId) public view virtual returns (address);
    function setApprovalForAll(address operator, bool approved) public virtual;
    function isApprovedForAll(address owner, address operator) public view virtual returns (bool);
    function transferFrom(address from, address to, uint256 tokenId) public virtual;
    function safeTransferFrom(address from, address to, uint256 tokenId) public virtual;
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public virtual;
}