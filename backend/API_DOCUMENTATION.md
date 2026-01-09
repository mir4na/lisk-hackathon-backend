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

### 5. Complete Profile (KYC & Bank Account) - MANDATORY
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

### 6. Connect Crypto Wallet (Optional)
Link a MetaMask wallet for on-chain transparency.

```bash
curl -X PUT http://localhost:8080/api/v1/user/wallet \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x123..."
  }'
```

### 7. Profile Management
View and update user details.

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

### 5. Create Invoice
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

### 6. Upload Invoice Documents
Upload Bill of Lading, Commercial Invoice, etc.

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/documents \
  -H "Authorization: Bearer <access_token>" \
  -F "document_type=bill_of_lading" \
  -F "file=@/path/to/bol.pdf"
```

### 7. Submit Invoice
Submit for admin review and grading.

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/submit \
  -H "Authorization: Bearer <access_token>"
```

### 8. Request Funding (Tokenize)
After admin approval, request funding to turn invoice into a pool.

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/funding-request \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
     "required_amount": 40000.00,
     "duration_days": 60
  }'
```

### 9. Repayment (Create VA)
When ready to repay the loan.

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

---

## Flow 4: Investor Operations

For users with `role: "investor"`.

### 1. Marketplac (List Pools)
Browse available investment opportunities.

```bash
curl -X GET "http://localhost:8080/api/v1/marketplace?page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

### 2. View Pool Detail
See specific details about an invoice and its funding tranches.

```bash
curl -X GET http://localhost:8080/api/v1/pools/<pool_id> \
  -H "Authorization: Bearer <access_token>"
```

### 3. Risk Questionnaire (Prerequisite for Catalyst)
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

### 4. Invest
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

### 5. My Portfolio
View investment summary and returns.

```bash
curl -X GET http://localhost:8080/api/v1/investments/portfolio \
  -H "Authorization: Bearer <access_token>"
```

---

## Flow 5: Admin Operations

For users with `role: "admin"`.

### 1. Approve KYC
Review and approve user identity verification.

```bash
curl -X POST http://localhost:8080/api/v1/admin/kyc/<kyc_id>/approve \
  -H "Authorization: Bearer <access_token>"
```

### 2. Approve Mitra Application
Review and approve company verification.

```bash
curl -X POST http://localhost:8080/api/v1/admin/mitra/<application_id>/approve \
  -H "Authorization: Bearer <access_token>"
```

### 3. Review Invoice (Grading)
Get grading suggestion and submit review decision.

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

### 4. Disburse Funds
When a pool is fully funded, disburse money to Mitra.

```bash
curl -X POST http://localhost:8080/api/v1/admin/pools/<pool_id>/disburse \
  -H "Authorization: Bearer <access_token>"
```

### 5. Platform Revenue
Check fees collected by the platform.

```bash
curl -X GET http://localhost:8080/api/v1/admin/platform/revenue \
  -H "Authorization: Bearer <access_token>"
```
