# VESSEL Backend API Documentation

Base URL: `http://localhost:8080`

## Important Notes

> 🏦 **Bank Account Required**: Users register their bank account at registration for fund disbursement. No crypto wallet needed - all transactions are in IDR through escrow.

> 💰 **Account Balance** (Saldo):
> - **Investor**: Total dana yang sedang dipinjamkan ke mitra-mitra (outstanding investments)
> - **Mitra**: Total saldo hutang yang dimiliki ke investor-investor (outstanding debt)
>
> This is NOT a virtual wallet. Funds flow through bank transfers via escrow.

> 🔑 **Role-Based Access Control (RBAC)**: Users select their role at registration:
> - **Investor (Pendana)**: Can invest in invoice pools
> - **Mitra (Eksportir)**: Can submit invoices for funding
> - **Guest (Tamu)**: Unregistered users who can only view marketplace

> 💵 **Currency**: All transactions use **IDR (Indonesian Rupiah)** for MVP phase.

---

## Flow 1: Guest & Public Access (No Auth)

Endpoints available to anyone without logging in.

### 1. Health Check
Check if the server is running correctly.

```bash
curl http://localhost:8080/health
```

### 2. Importer Payment Page
For global buyers (importers) to view their invoice payment details.

```bash
curl -X GET http://localhost:8080/api/v1/public/payments/<payment_id>
```

### 3. Pay Invoice (Importer)
For global buyers to pay an invoice.

```bash
curl -X POST http://localhost:8080/api/v1/public/payments/<payment_id>/pay \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 55000000
  }'
```

---

## Flow 2: User Onboarding & Authentication

The standard flow for new users (Investor/Mitra) to sign up, log in, and complete their profile (KYC).

### 1. Send OTP (Email Verification)
First, verify the user's email address by sending a code.

```bash
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "purpose": "registration"
  }'
```

### 2. Verify OTP
Validate the code received in email to get an `otp_token`.

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456",
    "purpose": "registration"
  }'
```

**Response** includes `"otp_token": "..."` which is required for registration.

### 3. Register User
Create a new account using the `otp_token`. Choose `role`: `"investor"` or `"mitra"`.

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123",
    "confirm_password": "password123",
    "role": "investor",
    "cooperative_agreement": true,
    "otp_token": "<token_from_verify_otp>"
  }'
```

### 4. Login
Obtain `access_token` and `refresh_token`.

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "user@example.com",
    "password": "password123"
  }'
```

### 5. Refresh Token
Get a new access token using the refresh token (when access token expires).

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token_from_login>"
  }'
```

### 6. Complete Profile (KYC & Bank Account) - MANDATORY
Before transacting, users MUST complete their profile details, including NIK, KTP, and Bank Account.

```bash
curl -X POST http://localhost:8080/api/v1/user/complete-profile \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "phone": "081234567890",
    "country": "Indonesia",
    
    "nik": "3201011234567890",
    "ktp_photo_url": "https://storage.example.com/ktp/xxx.jpg",
    "selfie_url": "https://storage.example.com/selfie/xxx.jpg",
    
    "bank_code": "bca",
    "account_number": "1234567890",
    "account_name": "John Doe"
  }'
```

### 7. Upload Documents (KTP/Selfie)
Upload identity documents for KYC verification.

```bash
curl -X POST http://localhost:8080/api/v1/user/documents \
  -H "Authorization: Bearer <access_token>" \
  -F "document_type=ktp" \
  -F "file=@/path/to/ktp.jpg"
```

Supported document types: `ktp`, `selfie`

### 8. Submit KYC
Submit KYC data for admin review.

```bash
curl -X POST http://localhost:8080/api/v1/user/kyc \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nik": "3201011234567890",
    "ktp_photo_url": "https://storage.example.com/ktp/xxx.jpg",
    "selfie_url": "https://storage.example.com/selfie/xxx.jpg"
  }'
```

### 9. Get KYC Status
Check current KYC verification status.

```bash
curl -X GET http://localhost:8080/api/v1/user/kyc \
  -H "Authorization: Bearer <access_token>"
```

### 10. Connect Crypto Wallet (Optional)
Link a MetaMask wallet for on-chain transparency.

```bash
curl -X PUT http://localhost:8080/api/v1/user/wallet \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x123..."
  }'
```

### 11. Profile Management
View and update user details.

**Get Profile:**
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>"
```

**Update Profile:**
```bash
curl -X PUT http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe Updated",
    "phone": "081234567891"
  }'
```

**Get Profile Data (Read Only):**
```bash
curl -X GET http://localhost:8080/api/v1/user/profile/data \
  -H "Authorization: Bearer <access_token>"
```

**Get Bank Account:**
```bash
curl -X GET http://localhost:8080/api/v1/user/profile/bank-account \
  -H "Authorization: Bearer <access_token>"
```

**Change Bank Account (Requires OTP):**
```bash
curl -X PUT http://localhost:8080/api/v1/user/profile/bank-account \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "otp_token": "<otp_verification_token>",
    "bank_code": "mandiri",
    "account_number": "0987654321",
    "account_name": "John Doe"
  }'
```

**Change Password:**
```bash
curl -X PUT http://localhost:8080/api/v1/user/profile/password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456",
    "confirm_password": "newpassword456"
  }'
```

**Get Supported Banks:**
```bash
curl -X GET http://localhost:8080/api/v1/user/profile/banks \
  -H "Authorization: Bearer <access_token>"
```

**Get Balance:**
```bash
curl -X GET http://localhost:8080/api/v1/user/balance \
  -H "Authorization: Bearer <access_token>"
```

---

## Flow 3: Mitra Operations (Exporter)

For users with `role: "mitra"`.

### 1. Apply as Mitra (Company Verification)
Submit company documents (NIB, NPWP) to get verified status.

```bash
curl -X POST http://localhost:8080/api/v1/user/mitra/apply \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "PT Contoh Eksportir",
    "company_type": "PT",
    "npwp": "1234567890123456",
    "annual_revenue": "1-5 Miliar"
  }'
```

### 2. Upload Mitra Document
Upload supporting docs for application (nib, akta_pendirian, etc).

```bash
curl -X POST http://localhost:8080/api/v1/user/mitra/documents \
  -H "Authorization: Bearer <access_token>" \
  -F "document_type=nib" \
  -F "file=@/path/to/nib.pdf"
```

### 3. Check Application Status

```bash
curl -X GET http://localhost:8080/api/v1/user/mitra/status \
  -H "Authorization: Bearer <access_token>"
```

### 4. Mitra Dashboard
View overview of active invoices, total loans, and repayment status.

```bash
curl -X GET http://localhost:8080/api/v1/mitra/dashboard \
  -H "Authorization: Bearer <access_token>"
```

### 5. Get Mitra Invoices
Get all invoices created by mitra.

```bash
curl -X GET http://localhost:8080/api/v1/mitra/invoices \
  -H "Authorization: Bearer <access_token>"
```

### 6. Get Active Invoices (Needing Repayment)
Get invoices that are currently funded and need repayment.

```bash
curl -X GET http://localhost:8080/api/v1/mitra/invoices/active \
  -H "Authorization: Bearer <access_token>"
```

### 7. Create Invoice
Once verified, create a new invoice draft.

```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "buyer_name": "Global Coffee Importers Ltd",
    "buyer_country": "United States",
    "invoice_number": "INV-2024-001",
    "currency": "USD",
    "amount": 50000.00,
    "issue_date": "2024-01-01",
    "due_date": "2024-03-01",
    "description": "Export Kopi Arabika Gayo 1 Container"
  }'
```

### 8. List My Invoices
Get all invoices created by current mitra.

```bash
curl -X GET http://localhost:8080/api/v1/invoices \
  -H "Authorization: Bearer <access_token>"
```

### 9. Get Invoice Detail

```bash
curl -X GET http://localhost:8080/api/v1/invoices/<invoice_id> \
  -H "Authorization: Bearer <access_token>"
```

### 10. Update Invoice (Draft Only)

```bash
curl -X PUT http://localhost:8080/api/v1/invoices/<invoice_id> \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "buyer_name": "Updated Buyer Name",
    "amount": 55000.00
  }'
```

### 11. Delete Invoice (Draft Only)

```bash
curl -X DELETE http://localhost:8080/api/v1/invoices/<invoice_id> \
  -H "Authorization: Bearer <access_token>"
```

### 12. Upload Invoice Documents
Upload Bill of Lading, Commercial Invoice, etc.

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/documents \
  -H "Authorization: Bearer <access_token>" \
  -F "document_type=bill_of_lading" \
  -F "file=@/path/to/bol.pdf"
```

Supported document types: `bill_of_lading`, `commercial_invoice`, `packing_list`, `certificate_of_origin`, `insurance`, `other`

### 13. Get Invoice Documents

```bash
curl -X GET http://localhost:8080/api/v1/invoices/<invoice_id>/documents \
  -H "Authorization: Bearer <access_token>"
```

### 14. Check Repeat Buyer
Check if buyer has previous transaction history (affects grading).

```bash
curl -X POST http://localhost:8080/api/v1/invoices/check-repeat-buyer \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "buyer_name": "Global Coffee Importers Ltd",
    "buyer_country": "United States"
  }'
```

### 15. Submit Invoice
Submit for admin review and grading.

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/submit \
  -H "Authorization: Bearer <access_token>"
```

### 16. Request Funding
After admin approval, request funding to create a pool from invoice.

```bash
curl -X POST http://localhost:8080/api/v1/invoices/funding-request \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_id": "<invoice_uuid>",
    "required_amount": 40000.00,
    "duration_days": 60
  }'
```

### 17. Get Repayment Breakdown
Get breakdown of repayment by tranche for a funded pool.

```bash
curl -X GET http://localhost:8080/api/v1/mitra/pools/<pool_id>/breakdown \
  -H "Authorization: Bearer <access_token>"
```

### 18. Get VA Payment Methods
Get available virtual account payment methods.

```bash
curl -X GET http://localhost:8080/api/v1/mitra/payment-methods \
  -H "Authorization: Bearer <access_token>"
```

### 19. Create VA for Repayment
Create virtual account for loan repayment.

```bash
curl -X POST http://localhost:8080/api/v1/mitra/repayment/va \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_id>",
    "bank_code": "bca",
    "amount": 55000000
  }'
```

### 20. Get VA Payment Status
Check virtual account payment status.

```bash
curl -X GET http://localhost:8080/api/v1/mitra/repayment/va/<va_id> \
  -H "Authorization: Bearer <access_token>"
```

### 21. Simulate VA Payment (MVP Only)
Simulate payment for testing purposes.

```bash
curl -X POST http://localhost:8080/api/v1/mitra/repayment/va/<va_id>/simulate-pay \
  -H "Authorization: Bearer <access_token>"
```

### 22. Request Disbursement (Exporter)
Request disbursement after pool is funded.

```bash
curl -X POST http://localhost:8080/api/v1/exporter/disbursement \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_id>"
  }'
```

---

## Flow 4: Currency Conversion

Endpoints for currency conversion and estimates.

### 1. Get Locked Exchange Rate
Get locked exchange rate for currency conversion (valid for limited time).

```bash
curl -X POST http://localhost:8080/api/v1/currency/convert \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USD",
    "to_currency": "IDR",
    "amount": 50000
  }'
```

### 2. Get Supported Currencies
Get list of supported currencies.

```bash
curl -X GET http://localhost:8080/api/v1/currency/supported \
  -H "Authorization: Bearer <access_token>"
```

### 3. Calculate Estimated Disbursement
Calculate estimated disbursement amount in IDR.

```bash
curl -X GET "http://localhost:8080/api/v1/currency/disbursement-estimate?currency=USD&amount=50000" \
  -H "Authorization: Bearer <access_token>"
```

---

## Flow 5: Payments (Prototype)

Payment endpoints for deposit and withdrawal.

### 1. Deposit
Deposit funds to account.

```bash
curl -X POST http://localhost:8080/api/v1/payments/deposit \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 10000000
  }'
```

### 2. Withdraw
Withdraw funds from account.

```bash
curl -X POST http://localhost:8080/api/v1/payments/withdraw \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000000
  }'
```

### 3. Get Balance
Get current account balance.

```bash
curl -X GET http://localhost:8080/api/v1/payments/balance \
  -H "Authorization: Bearer <access_token>"
```

---

## Flow 6: Investor Operations

For users with `role: "investor"`.

### 1. Marketplace (List Pools)
Browse available investment opportunities.

```bash
curl -X GET "http://localhost:8080/api/v1/marketplace?page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

### 2. View Pool Detail
See specific details about an invoice and its funding tranches.

```bash
curl -X GET http://localhost:8080/api/v1/marketplace/<pool_id>/detail \
  -H "Authorization: Bearer <access_token>"
```

### 3. Calculate Investment
Calculate potential returns for an investment amount.

```bash
curl -X POST http://localhost:8080/api/v1/marketplace/calculate \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "amount": 5000000,
    "tranche": "priority"
  }'
```

### 4. List Pools
Get list of all funding pools.

```bash
curl -X GET http://localhost:8080/api/v1/pools \
  -H "Authorization: Bearer <access_token>"
```

### 5. Get Pool by ID
Get specific pool details.

```bash
curl -X GET http://localhost:8080/api/v1/pools/<pool_id> \
  -H "Authorization: Bearer <access_token>"
```

### 6. Get Fundable Invoices
Get list of invoices available for funding.

```bash
curl -X GET http://localhost:8080/api/v1/invoices/fundable \
  -H "Authorization: Bearer <access_token>"
```

### 7. Risk Questionnaire (Prerequisite for Catalyst)
Must complete this to invest in Junior (Catalyst) Tranche.

**Get Questions:**
```bash
curl -X GET http://localhost:8080/api/v1/risk-questionnaire/questions \
  -H "Authorization: Bearer <access_token>"
```

**Submit Answers:**
```bash
curl -X POST http://localhost:8080/api/v1/risk-questionnaire \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "investment_purpose": 3,
    "loss_tolerance": 2,
    "tranche_understanding": 2
  }'
```

**Get Status:**
```bash
curl -X GET http://localhost:8080/api/v1/risk-questionnaire/status \
  -H "Authorization: Bearer <access_token>"
```

### 8. Invest
Invest in a pool. Choose `priority` (Senior) or `catalyst` (Junior) tranche.

**Priority Investment:**
```bash
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "amount": 5000000,
    "tranche": "priority",
    "tnc_accepted": true
  }'
```

**Catalyst Investment (Higher Risk/Return):**
```bash
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "amount": 1000000,
    "tranche": "catalyst",
    "tnc_accepted": true,
    "catalyst_consents": {
      "first_loss_consent": true,
      "risk_loss_consent": true,
      "not_bank_consent": true
    }
  }'
```

### 9. Confirm Investment
Confirm a pending investment.

```bash
curl -X POST http://localhost:8080/api/v1/investments/confirm \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "investment_id": "<investment_uuid>"
  }'
```

### 10. My Investments
Get list of all my investments.

```bash
curl -X GET http://localhost:8080/api/v1/investments \
  -H "Authorization: Bearer <access_token>"
```

### 11. My Portfolio
View investment summary and returns.

```bash
curl -X GET http://localhost:8080/api/v1/investments/portfolio \
  -H "Authorization: Bearer <access_token>"
```

### 12. Active Investments
Get list of currently active investments.

```bash
curl -X GET http://localhost:8080/api/v1/investments/active \
  -H "Authorization: Bearer <access_token>"
```

---

## Flow 7: Admin Operations

For users with `role: "admin"`. Default admin: `admin@vessel.com` / `adminpassword123`

### KYC Management

**Get Pending KYC:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/kyc/pending \
  -H "Authorization: Bearer <access_token>"
```

**Approve KYC:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/kyc/<kyc_id>/approve \
  -H "Authorization: Bearer <access_token>"
```

**Reject KYC:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/kyc/<kyc_id>/reject \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Document not clear"
  }'
```

### Mitra Application Management

**Get Pending Applications:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/mitra/pending \
  -H "Authorization: Bearer <access_token>"
```

**Get Application Detail:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/mitra/<application_id> \
  -H "Authorization: Bearer <access_token>"
```

**Approve Mitra Application:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/mitra/<application_id>/approve \
  -H "Authorization: Bearer <access_token>"
```

**Reject Mitra Application:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/mitra/<application_id>/reject \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Missing NIB document"
  }'
```

### Invoice Management

**Get Pending Invoices:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/invoices/pending \
  -H "Authorization: Bearer <access_token>"
```

**Get Invoice Review Data (Split-screen):**
```bash
curl -X GET http://localhost:8080/api/v1/admin/invoices/<invoice_id>/review \
  -H "Authorization: Bearer <access_token>"
```

**Get Grading Suggestion (AI/algo):**
```bash
curl -X GET http://localhost:8080/api/v1/admin/invoices/<invoice_id>/grade-suggestion \
  -H "Authorization: Bearer <access_token>"
```

**Approve & Assign Grade:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/invoices/<invoice_id>/approve \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "grade": "A",
    "advance_percentage": 80
  }'
```

**Reject Invoice:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/invoices/<invoice_id>/reject \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Insufficient documentation"
  }'
```

### Pool Management

**Disburse Funds:**
When a pool is fully funded, disburse money to Mitra.

```bash
curl -X POST http://localhost:8080/api/v1/admin/pools/<pool_id>/disburse \
  -H "Authorization: Bearer <access_token>"
```

**Close Pool:**
Close a pool and notify investors.

```bash
curl -X POST http://localhost:8080/api/v1/admin/pools/<pool_id>/close \
  -H "Authorization: Bearer <access_token>"
```

**Process Repayment:**
Process repayment for an invoice/pool.

```bash
curl -X POST http://localhost:8080/api/v1/admin/invoices/<invoice_id>/repay \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 55000000
  }'
```

### Balance Management (MVP)

**Grant Balance:**
Grant balance to a user (for testing purposes).

```bash
curl -X POST http://localhost:8080/api/v1/admin/balance/grant \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<user_uuid>",
    "amount": 10000000
  }'
```

### Platform Revenue

**Get Platform Revenue:**
Check fees collected by the platform.

```bash
curl -X GET http://localhost:8080/api/v1/admin/platform/revenue \
  -H "Authorization: Bearer <access_token>"
```

---

## API Route Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/health` | No | Health check |
| **Auth** |
| POST | `/api/v1/auth/send-otp` | No | Send OTP email |
| POST | `/api/v1/auth/verify-otp` | No | Verify OTP code |
| POST | `/api/v1/auth/register` | No | Register user |
| POST | `/api/v1/auth/login` | No | Login |
| POST | `/api/v1/auth/refresh` | No | Refresh token |
| **Public** |
| GET | `/api/v1/public/payments/:payment_id` | No | Get payment info |
| POST | `/api/v1/public/payments/:payment_id/pay` | No | Pay invoice |
| **User** |
| GET | `/api/v1/user/profile` | Yes | Get profile |
| PUT | `/api/v1/user/profile` | Yes | Update profile |
| POST | `/api/v1/user/kyc` | Yes | Submit KYC |
| GET | `/api/v1/user/kyc` | Yes | Get KYC status |
| POST | `/api/v1/user/complete-profile` | Yes | Complete profile |
| POST | `/api/v1/user/documents` | Yes | Upload document |
| GET | `/api/v1/user/balance` | Yes | Get balance |
| GET | `/api/v1/user/profile/data` | Yes | Get personal data |
| GET | `/api/v1/user/profile/bank-account` | Yes | Get bank account |
| PUT | `/api/v1/user/profile/bank-account` | Yes | Change bank account |
| PUT | `/api/v1/user/profile/password` | Yes | Change password |
| GET | `/api/v1/user/profile/banks` | Yes | Get supported banks |
| PUT | `/api/v1/user/wallet` | Yes | Update wallet |
| **Mitra Application** |
| POST | `/api/v1/user/mitra/apply` | Yes | Apply as mitra |
| GET | `/api/v1/user/mitra/status` | Yes | Get application status |
| POST | `/api/v1/user/mitra/documents` | Yes | Upload mitra document |
| **Currency** |
| POST | `/api/v1/currency/convert` | Yes | Get locked exchange rate |
| GET | `/api/v1/currency/supported` | Yes | Get supported currencies |
| GET | `/api/v1/currency/disbursement-estimate` | Yes | Calculate disbursement |
| **Payments** |
| POST | `/api/v1/payments/deposit` | Yes | Deposit |
| POST | `/api/v1/payments/withdraw` | Yes | Withdraw |
| GET | `/api/v1/payments/balance` | Yes | Get balance |
| **Invoices** |
| POST | `/api/v1/invoices` | Yes (Mitra) | Create invoice |
| POST | `/api/v1/invoices/funding-request` | Yes (Mitra) | Request funding |
| POST | `/api/v1/invoices/check-repeat-buyer` | Yes (Mitra) | Check repeat buyer |
| GET | `/api/v1/invoices` | Yes (Mitra) | List my invoices |
| GET | `/api/v1/invoices/fundable` | Yes | List fundable invoices |
| GET | `/api/v1/invoices/:id` | Yes | Get invoice |
| PUT | `/api/v1/invoices/:id` | Yes (Mitra) | Update invoice |
| DELETE | `/api/v1/invoices/:id` | Yes (Mitra) | Delete invoice |
| POST | `/api/v1/invoices/:id/submit` | Yes (Mitra) | Submit invoice |
| POST | `/api/v1/invoices/:id/documents` | Yes (Mitra) | Upload document |
| GET | `/api/v1/invoices/:id/documents` | Yes | Get documents |
| POST | `/api/v1/invoices/:id/tokenize` | Yes (Admin) | Tokenize invoice |
| POST | `/api/v1/invoices/:id/pool` | Yes (Admin) | Create pool |
| **Pools** |
| GET | `/api/v1/pools` | Yes | List pools |
| GET | `/api/v1/pools/:id` | Yes | Get pool |
| **Marketplace** |
| GET | `/api/v1/marketplace` | Yes | List marketplace |
| GET | `/api/v1/marketplace/:id/detail` | Yes | Get pool detail |
| POST | `/api/v1/marketplace/calculate` | Yes | Calculate investment |
| **Risk Questionnaire** |
| GET | `/api/v1/risk-questionnaire/questions` | Yes (Investor) | Get questions |
| POST | `/api/v1/risk-questionnaire` | Yes (Investor) | Submit answers |
| GET | `/api/v1/risk-questionnaire/status` | Yes (Investor) | Get status |
| **Investments** |
| POST | `/api/v1/investments` | Yes (Investor) | Invest |
| POST | `/api/v1/investments/confirm` | Yes (Investor) | Confirm investment |
| GET | `/api/v1/investments` | Yes (Investor) | List my investments |
| GET | `/api/v1/investments/portfolio` | Yes (Investor) | Get portfolio |
| GET | `/api/v1/investments/active` | Yes (Investor) | Get active investments |
| **Exporter** |
| POST | `/api/v1/exporter/disbursement` | Yes (Mitra) | Request disbursement |
| **Mitra Dashboard** |
| GET | `/api/v1/mitra/dashboard` | Yes (Mitra) | Get dashboard |
| GET | `/api/v1/mitra/invoices` | Yes (Mitra) | Get invoices |
| GET | `/api/v1/mitra/invoices/active` | Yes (Mitra) | Get active invoices |
| GET | `/api/v1/mitra/pools/:id/breakdown` | Yes (Mitra) | Get repayment breakdown |
| GET | `/api/v1/mitra/payment-methods` | Yes (Mitra) | Get payment methods |
| POST | `/api/v1/mitra/repayment/va` | Yes (Mitra) | Create VA payment |
| GET | `/api/v1/mitra/repayment/va/:id` | Yes (Mitra) | Get VA status |
| POST | `/api/v1/mitra/repayment/va/:id/simulate-pay` | Yes (Mitra) | Simulate payment |
| **Admin** |
| GET | `/api/v1/admin/kyc/pending` | Yes (Admin) | Get pending KYC |
| POST | `/api/v1/admin/kyc/:id/approve` | Yes (Admin) | Approve KYC |
| POST | `/api/v1/admin/kyc/:id/reject` | Yes (Admin) | Reject KYC |
| GET | `/api/v1/admin/invoices/pending` | Yes (Admin) | Get pending invoices |
| GET | `/api/v1/admin/invoices/:id/grade-suggestion` | Yes (Admin) | Get grade suggestion |
| GET | `/api/v1/admin/invoices/:id/review` | Yes (Admin) | Get review data |
| POST | `/api/v1/admin/invoices/:id/approve` | Yes (Admin) | Approve invoice |
| POST | `/api/v1/admin/invoices/:id/reject` | Yes (Admin) | Reject invoice |
| POST | `/api/v1/admin/pools/:id/disburse` | Yes (Admin) | Disburse funds |
| POST | `/api/v1/admin/pools/:id/close` | Yes (Admin) | Close pool |
| POST | `/api/v1/admin/invoices/:id/repay` | Yes (Admin) | Process repayment |
| GET | `/api/v1/admin/mitra/pending` | Yes (Admin) | Get pending applications |
| GET | `/api/v1/admin/mitra/:id` | Yes (Admin) | Get application detail |
| POST | `/api/v1/admin/mitra/:id/approve` | Yes (Admin) | Approve application |
| POST | `/api/v1/admin/mitra/:id/reject` | Yes (Admin) | Reject application |
| POST | `/api/v1/admin/balance/grant` | Yes (Admin) | Grant balance |
| GET | `/api/v1/admin/platform/revenue` | Yes (Admin) | Get platform revenue |
