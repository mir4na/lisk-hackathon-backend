const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MockUSDC", function () {
  let mockUSDC;
  let owner;
  let addr1;
  let addr2;

  const USDC_DECIMALS = 6;

  beforeEach(async function () {
    [owner, addr1, addr2] = await ethers.getSigners();

    const MockUSDC = await ethers.getContractFactory("MockUSDC");
    mockUSDC = await MockUSDC.deploy();
    await mockUSDC.waitForDeployment();
  });

  describe("Deployment", function () {
    it("Should set the correct name and symbol", async function () {
      expect(await mockUSDC.name()).to.equal("USD Coin (Mock)");
      expect(await mockUSDC.symbol()).to.equal("USDC");
    });

    it("Should return correct decimals", async function () {
      expect(await mockUSDC.decimals()).to.equal(6);
    });

    it("Should mint initial supply to deployer", async function () {
      const balance = await mockUSDC.balanceOf(owner.address);
      expect(balance).to.equal(ethers.parseUnits("1000000", USDC_DECIMALS));
    });

    it("Should set deployer as owner", async function () {
      expect(await mockUSDC.owner()).to.equal(owner.address);
    });
  });

  describe("Minting", function () {
    it("Should allow owner to mint tokens", async function () {
      await mockUSDC.mint(addr1.address, ethers.parseUnits("5000", USDC_DECIMALS));
      expect(await mockUSDC.balanceOf(addr1.address)).to.equal(ethers.parseUnits("5000", USDC_DECIMALS));
    });

    it("Should not allow non-owner to mint tokens", async function () {
      await expect(
        mockUSDC.connect(addr1).mint(addr2.address, ethers.parseUnits("5000", USDC_DECIMALS))
      ).to.be.reverted;
    });
  });

  describe("Faucet", function () {
    it("Should allow anyone to use faucet", async function () {
      await mockUSDC.connect(addr1).faucet(ethers.parseUnits("1000", USDC_DECIMALS));
      expect(await mockUSDC.balanceOf(addr1.address)).to.equal(ethers.parseUnits("1000", USDC_DECIMALS));
    });

    it("Should limit faucet to 10000 USDC per request", async function () {
      await expect(
        mockUSDC.connect(addr1).faucet(ethers.parseUnits("10001", USDC_DECIMALS))
      ).to.be.revertedWith("Max 10000 USDC per request");
    });

    it("Should allow max 10000 USDC from faucet", async function () {
      await mockUSDC.connect(addr1).faucet(ethers.parseUnits("10000", USDC_DECIMALS));
      expect(await mockUSDC.balanceOf(addr1.address)).to.equal(ethers.parseUnits("10000", USDC_DECIMALS));
    });
  });

  describe("Transfers", function () {
    it("Should allow transfers between accounts", async function () {
      await mockUSDC.transfer(addr1.address, ethers.parseUnits("1000", USDC_DECIMALS));
      expect(await mockUSDC.balanceOf(addr1.address)).to.equal(ethers.parseUnits("1000", USDC_DECIMALS));
    });

    it("Should allow approved transfers", async function () {
      await mockUSDC.approve(addr1.address, ethers.parseUnits("1000", USDC_DECIMALS));
      await mockUSDC.connect(addr1).transferFrom(owner.address, addr2.address, ethers.parseUnits("1000", USDC_DECIMALS));
      expect(await mockUSDC.balanceOf(addr2.address)).to.equal(ethers.parseUnits("1000", USDC_DECIMALS));
    });
  });
});
