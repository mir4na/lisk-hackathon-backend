const hre = require("hardhat");

async function main() {
  const [deployer] = await hre.ethers.getSigners();

  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString());

  // Deploy MockUSDC (for testnet)
  console.log("\n1. Deploying MockUSDC...");
  const MockUSDC = await hre.ethers.getContractFactory("MockUSDC");
  const mockUSDC = await MockUSDC.deploy();
  await mockUSDC.waitForDeployment();
  const mockUSDCAddress = await mockUSDC.getAddress();
  console.log("MockUSDC deployed to:", mockUSDCAddress);

  // Deploy InvoiceNFT
  console.log("\n2. Deploying InvoiceNFT...");
  const InvoiceNFT = await hre.ethers.getContractFactory("InvoiceNFT");
  const invoiceNFT = await InvoiceNFT.deploy();
  await invoiceNFT.waitForDeployment();
  const invoiceNFTAddress = await invoiceNFT.getAddress();
  console.log("InvoiceNFT deployed to:", invoiceNFTAddress);

  // Deploy InvoicePool
  console.log("\n3. Deploying InvoicePool...");
  const InvoicePool = await hre.ethers.getContractFactory("InvoicePool");
  const invoicePool = await InvoicePool.deploy(
    invoiceNFTAddress,
    mockUSDCAddress,
    deployer.address // Platform wallet
  );
  await invoicePool.waitForDeployment();
  const invoicePoolAddress = await invoicePool.getAddress();
  console.log("InvoicePool deployed to:", invoicePoolAddress);

  // Grant roles
  console.log("\n4. Setting up roles...");

  // Grant MINTER_ROLE to InvoicePool
  const MINTER_ROLE = await invoiceNFT.MINTER_ROLE();
  await invoiceNFT.grantRole(MINTER_ROLE, invoicePoolAddress);
  console.log("Granted MINTER_ROLE to InvoicePool");

  // Grant OPERATOR_ROLE to deployer
  const OPERATOR_ROLE = await invoicePool.OPERATOR_ROLE();
  await invoicePool.grantRole(OPERATOR_ROLE, deployer.address);
  console.log("Granted OPERATOR_ROLE to deployer");

  console.log("\n========================================");
  console.log("Deployment Summary:");
  console.log("========================================");
  console.log("Network:", hre.network.name);
  console.log("MockUSDC:", mockUSDCAddress);
  console.log("InvoiceNFT:", invoiceNFTAddress);
  console.log("InvoicePool:", invoicePoolAddress);
  console.log("========================================");

  // Save deployment info
  const fs = require("fs");
  const deploymentInfo = {
    network: hre.network.name,
    chainId: (await deployer.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    contracts: {
      MockUSDC: mockUSDCAddress,
      InvoiceNFT: invoiceNFTAddress,
      InvoicePool: invoicePoolAddress,
    },
    deployedAt: new Date().toISOString(),
  };

  fs.writeFileSync(
    `./deployments-${hre.network.name}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );
  console.log(`\nDeployment info saved to deployments-${hre.network.name}.json`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
