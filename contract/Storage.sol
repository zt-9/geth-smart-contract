// SPDX-License-Identifier: MIT

pragma solidity 0.8.26;

contract Storage {
    mapping(address => uint256) records;

    function store(address user, uint256 num) public {
        records[user] = num;
    }

    function retrieve(address user) public view returns (uint256) {
        return records[user];
    }
}