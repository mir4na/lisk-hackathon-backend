const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("InvoicePool", function () {
  let invoiceNFT;
  let mockUSDC;
  let invoicePool;
  let owner;
  let exporter;
  let investor1;
  let investor2;

  const USDC_DECIMALS = 6;
  const sampleInvoice = {
    invoiceNumber: "INV-2024-001",
    amount: ethers.parseUnits("10000", USDC_DECIMALS),
    advanceAmount: ethers.parseUnits("8000", USDC_DECIMALS),
    interestRate: 1000, // 10%
    issueDate: Math.floor(Date.now() / 1000),
    dueDate: Math.floor(Date.now() / 1000) + 60 * 24 * 60 * 60,
    buyerCountry: "Germany",
    documentHash: "QmXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    uri: "ipfs://QmYyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
  };

  beforeEach(async function () {
    [owner, exporter, investor1, investor2] = await ethers.getSigners();

    // Deploy MockUSDC
    const MockUSDC = await ethers.getContractFactory("MockUSDC");
    mockUSDC = await MockUSDC.deploy();
    await mockUSDC.waitForDeployment();

    // Deploy InvoiceNFT
    const InvoiceNFT = await ethers.getContractFactory("InvoiceNFT");
    invoiceNFT = await InvoiceNFT.deploy();
    await invoiceNFT.waitForDeployment();

    // Deploy InvoicePool
    const InvoicePool = await ethers.getContractFactory("InvoicePool");
    invoicePool = await InvoicePool.deploy(
      await invoiceNFT.getAddress(),
      await mockUSDC.getAddress(),
      owner.address
    );
    await invoicePool.waitForDeployment();

    // Grant roles
    const MINTER_ROLE = await invoiceNFT.MINTER_ROLE();
    await invoiceNFT.grantRole(MINTER_ROLE, await invoicePool.getAddress());

    // Mint USDC to investors
    await mockUSDC.mint(investor1.address, ethers.parseUnits("50000", USDC_DECIMALS));
    await mockUSDC.mint(investor2.address, ethers.parseUnits("50000", USDC_DECIMALS));

    // Approve pool to spend USDC
    await mockUSDC.connect(investor1).approve(
      await invoicePool.getAddress(),
      ethers.parseUnits("100000", USDC_DECIMALS)
    );
    await mockUSDC.connect(investor2).approve(
      await invoicePool.getAddress(),
      ethers.parseUnits("100000", USDC_DECIMALS)
    );
  });

  async function mintAndVerifyInvoice() {
    await invoiceNFT.mintInvoice(
      exporter.address,
      sampleInvoice.invoiceNumber,
      sampleInvoice.amount,
      sampleInvoice.advanceAmount,
      sampleInvoice.interestRate,
      sampleInvoice.issueDate,
      sampleInvoice.dueDate,
      sampleInvoice.buyerCountry,
      sampleInvoice.documentHash,
      sampleInvoice.uri
    );
    await invoiceNFT.verifyShipment(1);
    return 1;
  }

  describe("Pool Creation", function () {
    it("Should create a funding pool", async function () {
      const tokenId = await mintAndVerifyInvoice();
      await invoicePool.createPool(tokenId);

      const pool = await invoicePool.getPool(tokenId);
      expect(pool.targetAmount).to.equal(sampleInvoice.advanceAmount);
      expect(pool.status).to.equal(0); // Open
    });

    it("Should not create pool for unverified invoice", async function () {
      await invoiceNFT.mintInvoice(
        exporter.address,
        sampleInvoice.invoiceNumber,
        sampleInvoice.amount,
        sampleInvoice.advanceAmount,
        sampleInvoice.interestRate,
        sampleInvoice.issueDate,
        sampleInvoice.dueDate,
        sampleInvoice.buyerCountry,
        sampleInvoice.documentHash,
        sampleInvoice.uri
      );

      await expect(invoicePool.createPool(1)).to.be.revertedWith("Invoice not fundable");
    });
  });

  describe("Investment", function () {
    beforeEach(async function () {
      await mintAndVerifyInvoice();
      await invoicePool.createPool(1);
    });

    it("Should allow investment", async function () {
      const investAmount = ethers.parseUnits("5000", USDC_DECIMALS);
      await invoicePool.connect(investor1).invest(1, investAmount);

      const pool = await invoicePool.getPool(1);
      expect(pool.fundedAmount).to.equal(investAmount);
      expect(pool.investorCount).to.equal(1);
    });

    it("Should track investments correctly", async function () {
      const investAmount = ethers.parseUnits("5000", USDC_DECIMALS);
      await invoicePool.connect(investor1).invest(1, investAmount);

      const investments = await invoicePool.getPoolInvestments(1);
      expect(investments.length).to.equal(1);
      expect(investments[0].investor).to.equal(investor1.address);
      expect(investments[0].amount).to.equal(investAmount);
    });

    it("Should fill pool when target reached", async function () {
      await invoicePool.connect(investor1).invest(1, ethers.parseUnits("5000", USDC_DECIMALS));
      await invoicePool.connect(investor2).invest(1, ethers.parseUnits("3000", USDC_DECIMALS));

      const pool = await invoicePool.getPool(1);
      expect(pool.status).to.equal(1); // Filled
    });

    it("Should not allow over-investment", async function () {
      await expect(
        invoicePool.connect(investor1).invest(1, ethers.parseUnits("10000", USDC_DECIMALS))
      ).to.be.revertedWith("Amount exceeds remaining capacity");
    });
  });

  describe("Disbursement", function () {
    beforeEach(async function () {
      await mintAndVerifyInvoice();
      await invoicePool.createPool(1);
      await invoicePool.connect(investor1).invest(1, ethers.parseUnits("8000", USDC_DECIMALS));
    });

    it("Should disburse funds to exporter", async function () {
      const exporterBalanceBefore = await mockUSDC.balanceOf(exporter.address);
      await invoicePool.disburse(1);
      const exporterBalanceAfter = await mockUSDC.balanceOf(exporter.address);

      expect(exporterBalanceAfter - exporterBalanceBefore).to.equal(
        ethers.parseUnits("8000", USDC_DECIMALS)
      );
    });

    it("Should update pool status after disbursement", async function () {
      await invoicePool.disburse(1);
      const pool = await invoicePool.getPool(1);
      expect(pool.status).to.equal(2); // Disbursed
    });

    it("Should not disburse unfilled pool", async function () {
      await invoicePool.createPool(1);
      await expect(invoicePool.disburse(1)).to.be.reverted;
    });
  });

  describe("Repayment", function () {
    beforeEach(async function () {
      await mintAndVerifyInvoice();
      await invoicePool.createPool(1);
      await invoicePool.connect(investor1).invest(1, ethers.parseUnits("4000", USDC_DECIMALS));
      await invoicePool.connect(investor2).invest(1, ethers.parseUnits("4000", USDC_DECIMALS));
      await invoicePool.disburse(1);

      // Mint USDC to owner for repayment
      await mockUSDC.mint(owner.address, ethers.parseUnits("10000", USDC_DECIMALS));
      await mockUSDC.approve(await invoicePool.getAddress(), ethers.parseUnits("10000", USDC_DECIMALS));
    });

    it("Should process repayment and distribute to investors", async function () {
      const investor1BalanceBefore = await mockUSDC.balanceOf(investor1.address);
      const investor2BalanceBefore = await mockUSDC.balanceOf(investor2.address);

      await invoicePool.processRepayment(1, ethers.parseUnits("10000", USDC_DECIMALS));

      const investor1BalanceAfter = await mockUSDC.balanceOf(investor1.address);
      const investor2BalanceAfter = await mockUSDC.balanceOf(investor2.address);

      // Both should receive roughly equal amounts (minus platform fee)
      expect(investor1BalanceAfter > investor1BalanceBefore).to.be.true;
      expect(investor2BalanceAfter > investor2BalanceBefore).to.be.true;
    });

    it("Should close pool after repayment", async function () {
      await invoicePool.processRepayment(1, ethers.parseUnits("10000", USDC_DECIMALS));
      const pool = await invoicePool.getPool(1);
      expect(pool.status).to.equal(3); // Closed
    });
  });

  describe("Admin Functions", function () {
    it("Should update platform fee", async function () {
      await invoicePool.setPlatformFee(300); // 3%
      expect(await invoicePool.platformFeeBps()).to.equal(300);
    });

    it("Should not allow fee above 10%", async function () {
      await expect(invoicePool.setPlatformFee(1001)).to.be.revertedWith("Fee too high");
    });

    it("Should pause and unpause", async function () {
      await invoicePool.pause();
      await expect(
        invoicePool.connect(investor1).invest(1, ethers.parseUnits("1000", USDC_DECIMALS))
      ).to.be.reverted;

      await invoicePool.unpause();
      // Should work again after unpause
    });
  });
});
