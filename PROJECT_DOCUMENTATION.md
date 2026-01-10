# VESSEL: Bridging Trade Finance Gap with RWA Tokenization

## 1. DApp Use Case Impact and Real Problems Solved

### The Problem: $2.5 Trillion Trade Finance Gap
Small and Medium Enterprises (SMEs) in emerging markets struggle to get access to trade finance. Traditional banks reject 45% of SME applications due to high operational costs and lack of collateral, leaving a massive financing gap.

### The Solution: VESSEL
VESSEL is a decentralised invoice financing platform that allows SMEs ("Mitra") to turn their unpaid invoices into liquid assets. By tokenizing invoices as **Real World Assets (RWA)** on the Lisk blockchain, we allow global investors to fund these invoices in a fractionalized, transparent manner.

**Key Impacts:**
*   **Liquidity for SMEs**: Exporters get paid immediately (up to 80% advance), allowing them to fulfill more orders.
*   **Access for Investors**: Retail and institutional investors can earn attractive yields (10-15% APY) backed by real trade activities.
*   **Transparency**: Every step from invoice creation, funding, to repayment is recorded on-chain, eliminating fraud.

---

## 2. Technical Implementation

### Architecture Overview
VESSEL utilizes a hybrid architecture combining a robust **Web2 Backend** for user experience and complex business logic with a **Web3 Layer** for transparency, settlement, and asset ownership.

### Tech Stack
*   **Backend**: Go (Golang) with Gin Framework.
    *   Modular **Service-Repository Pattern** for scalability and testability.
    *   **PostgreSQL** for relational data persistence.
    *   **Redis** for caching and rate limiting.
*   **Smart Contracts**: Solidity (EVM Compatible).
    *   Deployed on **Lisk Sepolia Testnet**.
    *   Uses **OpenZeppelin** standards for security.
*   **Frontend**: Next.js 14 (App Router) with TypeScript & TailwindCSS.
    *   Modern, responsive, and "Premium" feel UI.

### Key Features
1.  **Invoice Tokenization (NFT)**:
    *   Each approved invoice is minted as a unique **ERC721 NFT**.
    *   Metadata includes invoice details, risk grade, and document hash (IPFS).
    *   Ensures single ownership and prevents double-financing.

2.  **Fractionalized Investment Pools**:
    *   Investors don't buy the whole invoice. They contribute to a **Funding Pool**.
    *   **Waterfall Model**: Supports **Priority Tranche** (lower risk, lower yield) and **Catalyst Tranche** (higher risk, higher yield).
    *   Smart contract records every investment and calculates pro-rata ownership.

3.  **Abstracted Payment & Settlement**:
    *   Payment is handled via fiat rails (IDR/USD) for easy onboarding.
    *   **On-Chain Accounting**: The smart contract acts as the "Source of Truth". Every disbursement and repayment is recorded on-chain via Operator middleware, ensuring transparency even for fiat transactions.

---

## 3. UI/UX & Usability

We prioritize a "Web2-like" experience where blockchain complexity is abstracted away.

*   **For Exporters (Mitra)**: Simple dashboard to upload invoice, check status, and receive funds directly to bank accounts.
*   **For Investors**: Marketplace view with clear risk grading (A/B/C), yield calculator, and portfolio tracking.
*   **Admin Dashboard**: Comprehensive tools for KYC verification, invoice risk assessment, and payment processing.

---

## 4. Innovation & Uniqueness

*   **Risk-Based Tranching**: Unlike simple crowdfunding, VESSEL introduces structured finance concepts (Senior/Junior tranches) to DeFi, allowing risk-averse and risk-seeking investors to participate in the same asset.
*   **Hybrid "On-Chain Accounting"**: We acknowledge that fully crypto-native payments are hard for traditional SMEs. Our **Abstracted Payment** model records fiat settlements on-chain, providing the transparency of DeFi with the usability of Fintech.
*   **Dynamic Risk Grading**: Invoices are not just valid/invalid. They are graded (A/B/C) based on Buyer history, Country risk, and Document completeness.

---

## 5. Smart Contracts Logic

### `InvoiceNFT.sol`
*   **Standard**: ERC721URIStorage.
*   **Function**: Represents the legal claim to the invoice payment.
*   **Status Management**: Tracks lifecycle `Active` -> `Tokenized` -> `Funded` -> `Repaid`.

### `InvoicePool.sol`
*   **Function**: Manages the ledger for each invoice funding.
*   **Logic**:
    *   `createPool`: Initializes funding target based on Invoice Advance Amount.
    *   `recordInvestment`: Logs investor deposits (Abstracted).
    *   `recordDisbursement`: Locks the pool and records fund transfer to Exporter.
    *   `recordRepayment`: Distributes returns to investors based on the Waterfall model (Priority first).

---

## 6. How to Run

### Prerequisities
*   Docker & Docker Compose
*   Go 1.22+
*   Node.js 18+

### Setup
1.  **Backend**:
    ```bash
    cd backend
    go mod download
    go run main.go
    ```
    *Server runs on port 8080.*

2.  **Smart Contracts**:
    ```bash
    cd smart-contract
    npm install
    npx hardhat compile
    ```

3.  **Frontend**:
    ```bash
    cd ../vessel
    npm install
    npm run dev
    ```

---

*Built with ❤️ for Lisk Builders Hackathon*
