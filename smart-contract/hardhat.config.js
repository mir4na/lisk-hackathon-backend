require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config({ path: "../.env" });

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: "0.8.24",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  // =============================================================================
  // IMPORTANT: This project is configured for Lisk Sepolia testnet ONLY
  // Get testnet ETH from: https://sepolia-faucet.lisk.com/
  // =============================================================================
  networks: {
    hardhat: {
      chainId: 31337,
    },
    lisk_sepolia: {
      url: process.env.LISK_SEPOLIA_RPC_URL || "https://rpc.sepolia-api.lisk.com",
      chainId: 4202,
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
    },
  },
  etherscan: {
    apiKey: {
      lisk_sepolia: "placeholder",
    },
    customChains: [
      {
        network: "lisk_sepolia",
        chainId: 4202,
        urls: {
          apiURL: "https://sepolia-blockscout.lisk.com/api",
          browserURL: "https://sepolia-blockscout.lisk.com",
        },
      },
    ],
  },
};
