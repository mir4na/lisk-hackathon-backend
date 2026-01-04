// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

import "./InvoiceNFT.sol";

/**
 * @title InvoicePool
 * @dev Manages funding pools for invoice NFTs
 * Investors can fund invoices and receive returns when buyers repay
 */
contract InvoicePool is AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");

    // Contracts
    InvoiceNFT public invoiceNFT;
    IERC20 public stablecoin; // USDC or similar

    // Platform fee in basis points (e.g., 200 = 2%)
    uint256 public platformFeeBps = 200;
    address public platformWallet;

    // Pool status enum
    enum PoolStatus {
        Open,
        Filled,
        Disbursed,
        Closed,
        Defaulted
    }

    // Pool structure
    struct Pool {
        uint256 tokenId;
        uint256 targetAmount;
        uint256 fundedAmount;
        uint256 investorCount;
        uint256 interestRate;
        uint256 dueDate;
        address exporter;
        PoolStatus status;
        uint256 openedAt;
        uint256 filledAt;
        uint256 disbursedAt;
        uint256 closedAt;
    }

    // Investment structure
    struct Investment {
        address investor;
        uint256 amount;
        uint256 expectedReturn;
        uint256 actualReturn;
        bool claimed;
        uint256 investedAt;
    }

    // Mappings
    mapping(uint256 => Pool) public pools; // tokenId => Pool
    mapping(uint256 => Investment[]) public poolInvestments; // tokenId => investments
    mapping(address => uint256[]) public investorPools; // investor => tokenIds they invested in

    // Events
    event PoolCreated(uint256 indexed tokenId, uint256 targetAmount, uint256 interestRate);
    event InvestmentMade(uint256 indexed tokenId, address indexed investor, uint256 amount, uint256 expectedReturn);
    event PoolFilled(uint256 indexed tokenId, uint256 totalAmount, uint256 investorCount);
    event FundsDisbursed(uint256 indexed tokenId, address indexed exporter, uint256 amount);
    event RepaymentReceived(uint256 indexed tokenId, uint256 amount);
    event InvestorPaid(uint256 indexed tokenId, address indexed investor, uint256 amount);
    event PoolClosed(uint256 indexed tokenId);
    event PoolDefaulted(uint256 indexed tokenId);

    constructor(address _invoiceNFT, address _stablecoin, address _platformWallet) {
        invoiceNFT = InvoiceNFT(_invoiceNFT);
        stablecoin = IERC20(_stablecoin);
        platformWallet = _platformWallet;

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(OPERATOR_ROLE, msg.sender);
    }

    /**
     * @dev Create a funding pool for an invoice
     */
    function createPool(uint256 tokenId) external onlyRole(OPERATOR_ROLE) whenNotPaused {
        require(invoiceNFT.isFundable(tokenId), "Invoice not fundable");
        require(pools[tokenId].targetAmount == 0, "Pool already exists");

        InvoiceNFT.Invoice memory invoice = invoiceNFT.getInvoice(tokenId);

        pools[tokenId] = Pool({
            tokenId: tokenId,
            targetAmount: invoice.advanceAmount,
            fundedAmount: 0,
            investorCount: 0,
            interestRate: invoice.interestRate,
            dueDate: invoice.dueDate,
            exporter: invoice.exporter,
            status: PoolStatus.Open,
            openedAt: block.timestamp,
            filledAt: 0,
            disbursedAt: 0,
            closedAt: 0
        });

        // Update NFT status
        invoiceNFT.updateStatus(tokenId, InvoiceNFT.InvoiceStatus.Funded);

        emit PoolCreated(tokenId, invoice.advanceAmount, invoice.interestRate);
    }

    /**
     * @dev Invest in a pool
     */
    function invest(uint256 tokenId, uint256 amount) external nonReentrant whenNotPaused {
        Pool storage pool = pools[tokenId];
        require(pool.targetAmount > 0, "Pool does not exist");
        require(pool.status == PoolStatus.Open, "Pool not open");
        require(amount > 0, "Amount must be positive");

        uint256 remaining = pool.targetAmount - pool.fundedAmount;
        require(amount <= remaining, "Amount exceeds remaining capacity");

        // Calculate expected return based on interest rate and time
        uint256 daysToMaturity = (pool.dueDate - block.timestamp) / 1 days;
        if (daysToMaturity == 0) daysToMaturity = 1;
        uint256 expectedReturn = amount + (amount * pool.interestRate * daysToMaturity) / (365 * 10000);

        // Transfer stablecoin from investor
        stablecoin.safeTransferFrom(msg.sender, address(this), amount);

        // Record investment
        poolInvestments[tokenId].push(Investment({
            investor: msg.sender,
            amount: amount,
            expectedReturn: expectedReturn,
            actualReturn: 0,
            claimed: false,
            investedAt: block.timestamp
        }));

        pool.fundedAmount += amount;
        pool.investorCount++;
        investorPools[msg.sender].push(tokenId);

        emit InvestmentMade(tokenId, msg.sender, amount, expectedReturn);

        // Check if pool is filled
        if (pool.fundedAmount >= pool.targetAmount) {
            pool.status = PoolStatus.Filled;
            pool.filledAt = block.timestamp;
            emit PoolFilled(tokenId, pool.fundedAmount, pool.investorCount);
        }
    }

    /**
     * @dev Disburse funds to exporter
     */
    function disburse(uint256 tokenId) external onlyRole(OPERATOR_ROLE) nonReentrant {
        Pool storage pool = pools[tokenId];
        require(pool.status == PoolStatus.Filled, "Pool not filled");

        pool.status = PoolStatus.Disbursed;
        pool.disbursedAt = block.timestamp;

        // Transfer funds to exporter
        stablecoin.safeTransfer(pool.exporter, pool.fundedAmount);

        emit FundsDisbursed(tokenId, pool.exporter, pool.fundedAmount);
    }

    /**
     * @dev Process repayment from buyer
     */
    function processRepayment(uint256 tokenId, uint256 amount) external onlyRole(OPERATOR_ROLE) nonReentrant {
        Pool storage pool = pools[tokenId];
        require(pool.status == PoolStatus.Disbursed, "Pool not disbursed");

        // Calculate platform fee
        uint256 platformFee = (amount * platformFeeBps) / 10000;
        uint256 distributableAmount = amount - platformFee;

        // Transfer platform fee
        if (platformFee > 0) {
            stablecoin.safeTransferFrom(msg.sender, platformWallet, platformFee);
        }

        // Transfer distributable amount to contract
        stablecoin.safeTransferFrom(msg.sender, address(this), distributableAmount);

        // Distribute to investors proportionally
        Investment[] storage investments = poolInvestments[tokenId];
        for (uint256 i = 0; i < investments.length; i++) {
            Investment storage inv = investments[i];
            uint256 proportion = (inv.amount * 1e18) / pool.fundedAmount;
            uint256 payout = (distributableAmount * proportion) / 1e18;

            inv.actualReturn = payout;
            stablecoin.safeTransfer(inv.investor, payout);
            inv.claimed = true;

            emit InvestorPaid(tokenId, inv.investor, payout);
        }

        pool.status = PoolStatus.Closed;
        pool.closedAt = block.timestamp;

        // Update NFT status
        invoiceNFT.updateStatus(tokenId, InvoiceNFT.InvoiceStatus.Repaid);

        emit RepaymentReceived(tokenId, amount);
        emit PoolClosed(tokenId);
    }

    /**
     * @dev Mark pool as defaulted
     */
    function markDefaulted(uint256 tokenId) external onlyRole(OPERATOR_ROLE) {
        Pool storage pool = pools[tokenId];
        require(pool.status == PoolStatus.Disbursed, "Pool not disbursed");
        require(block.timestamp > pool.dueDate + 30 days, "Grace period not passed");

        pool.status = PoolStatus.Defaulted;

        // Update NFT status
        invoiceNFT.updateStatus(tokenId, InvoiceNFT.InvoiceStatus.Defaulted);

        emit PoolDefaulted(tokenId);
    }

    /**
     * @dev Get pool details
     */
    function getPool(uint256 tokenId) external view returns (Pool memory) {
        return pools[tokenId];
    }

    /**
     * @dev Get investments for a pool
     */
    function getPoolInvestments(uint256 tokenId) external view returns (Investment[] memory) {
        return poolInvestments[tokenId];
    }

    /**
     * @dev Get pools an investor has invested in
     */
    function getInvestorPools(address investor) external view returns (uint256[] memory) {
        return investorPools[investor];
    }

    /**
     * @dev Get remaining capacity in pool
     */
    function getRemainingCapacity(uint256 tokenId) external view returns (uint256) {
        Pool memory pool = pools[tokenId];
        if (pool.status != PoolStatus.Open) return 0;
        return pool.targetAmount - pool.fundedAmount;
    }

    /**
     * @dev Update platform fee
     */
    function setPlatformFee(uint256 newFeeBps) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newFeeBps <= 1000, "Fee too high"); // Max 10%
        platformFeeBps = newFeeBps;
    }

    /**
     * @dev Update platform wallet
     */
    function setPlatformWallet(address newWallet) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newWallet != address(0), "Invalid address");
        platformWallet = newWallet;
    }

    /**
     * @dev Pause/Unpause
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @dev Emergency withdraw (admin only)
     */
    function emergencyWithdraw(address token, uint256 amount) external onlyRole(DEFAULT_ADMIN_ROLE) {
        IERC20(token).safeTransfer(msg.sender, amount);
    }
}
