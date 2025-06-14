// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

abstract contract ERC1155 {
    event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
    event TransferBatch(
        address indexed operator,
        address indexed from,
        address indexed to,
        uint256[] ids,
        uint256[] values
    );
    event ApprovalForAll(address indexed account, address indexed operator, bool approved);
    event URI(string value, uint256 indexed id);

    function name() public view virtual returns (string memory);
    function symbol() public view virtual returns (string memory);
    function uri(uint256) public view virtual returns (string memory);
    function balanceOf(address account, uint256 id) public view virtual returns (uint256);
    function balanceOfBatch(
        address[] memory accounts,
        uint256[] memory ids
    ) public view virtual returns (uint256[] memory);
    function setApprovalForAll(address operator, bool approved) public virtual;
    function isApprovedForAll(address account, address operator) public view virtual returns (bool);
    function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes memory data) public virtual;
    function safeBatchTransferFrom(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory values,
        bytes memory data
    ) public virtual;
}