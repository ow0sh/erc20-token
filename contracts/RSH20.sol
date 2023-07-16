// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

interface IERC20 {
    function totalSupply() external view returns (uint256);

    function balanceOf(address account) external view returns (uint256);

    function transfer(
        address recipient,
        uint256 amount
    ) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
}

contract RSH20 is IERC20 {
    string public constant name = "Test RSH20 tokens";
    string public constant symbol = "RSH";
    uint public constant decimals = 18;

    mapping(address => uint256) balances;

    uint256 _totalsupply = 10 ether;

    constructor() {
        balances[msg.sender] = _totalsupply;
    }

    function totalSupply() public view override returns (uint256) {
        return _totalsupply;
    }

    function balanceOf(address account) public view override returns (uint256) {
        return balances[account];
    }

    function transfer(
        address recipient,
        uint256 amount
    ) public override returns (bool) {
        require(amount <= balances[msg.sender]);
        balances[msg.sender] -= amount;
        balances[recipient] += amount;
        emit Transfer(msg.sender, recipient, amount);
        return true;
    }
}
