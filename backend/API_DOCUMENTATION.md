# VESSEL Backend API Documentation

Base URL: `http://localhost:8080`

## Important Notes

> 🏦 **Bank Account Required**: Users register their bank account at registration for fund disbursement. No crypto wallet needed - all transactions are in IDR through escrow.

> 💰 **Account Balance** (Saldo):
> - **Investor**: Total dana yang sedang dipinjamkan ke mitra-mitra (outstanding investments)
> - **Mitra**: Total saldo hutang yang dimiliki ke investor-investor (outstanding debt)
>
> This is NOT a virtual wallet. Funds flow through bank transfers via escrow.

> 🔄 **Fund Flow**:
> - **Investor → Mitra**: Investor transfers to escrow → escrow disburses to mitra's bank account
> - **Mitra → Investor**: Mitra pays to VA (escrow) → escrow distributes to investors' bank accounts

> 🔑 **Role-Based Access Control (RBAC)**: Users select their role at registration:
> - **Investor (Pendana)**: Can invest in invoice pools
> - **Mitra (Eksportir)**: Can submit invoices for funding
> - **Guest (Tamu)**: Unregistered users who can only view marketplace

> 💵 **Currency**: All transactions use **IDR (Indonesian Rupiah)** for MVP phase.

> 📊 **Interest Calculation**: Flat rate formula: `Principal + (Principal × Rate/100)`
> - No time-based (pro-rata) calculation for MVP
> - Example: Rp 10,000,000 at 12% = Rp 10,000,000 + Rp 1,200,000 = Rp 11,200,000

---

## Role Permissions Matrix

| Feature | Guest | Investor | Mitra | Admin |
|---------|-------|----------|-------|-------|
| View Marketplace | ✅ | ✅ | ✅ | ✅ |
| View Pool Details | ✅ | ✅ | ✅ | ✅ |
| Invest in Pools | ❌ | ✅ | ❌ | ✅ |
| Complete Risk Questionnaire | ❌ | ✅ | ❌ | ❌ |
| Submit Invoices | ❌ | ❌ | ✅* | ✅ |
| View Mitra Dashboard | ❌ | ❌ | ✅ | ✅ |
| Manage Bank Account | ❌ | ✅ | ✅ | ✅ |
| Admin Panel Access | ❌ | ❌ | ❌ | ✅ |

*Mitra must be verified first

---

## Health Check

```bash
curl http://localhost:8080/health
```

---

## Importer Payment (PUBLIC - No Auth Required)

These endpoints are for importers (buyers) who are NOT users of the application. They receive a payment ID via email and use it to pay.

### Get Payment Info

```bash
curl -X GET http://localhost:8080/api/v1/public/payments/<payment_id>
```

### Pay Invoice (Importer)

```bash
curl -X POST http://localhost:8080/api/v1/public/payments/<payment_id>/pay \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 55000000
  }'
```

Response:
```json
{
  "payment_id": "uuid",
  "status": "paid",
  "amount_paid": 55000000,
  "tx_hash": "0x...",
  "message": "Payment processed successfully. Funds have been distributed to investors.",
  "paid_at": "2024-01-15T12:00:00Z"
}
```

---

## Authentication

### 1. Send OTP for Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "purpose": "registration"
  }'
```

### 2. Verify OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456",
    "purpose": "registration"
  }'
```

### 3. Register User

Users choose their role at registration. Roles are mutually exclusive. **Guest role is not available** - unregistered users viewing the app are guests.

**Available Roles:**
- `investor` - Can invest in invoice pools
- `mitra` - Can submit invoices for funding (requires verification)

**KYC Verification at Registration:**
Registration now requires identity verification:
- KTP Photo upload
- Selfie with KTP
- NIK (16 digits)
- Full Name (must match KTP and bank account)
- Bank Account (for disbursement)

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

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | User email (verified via OTP) |
| username | string | Yes | Unique username |
| password | string | Yes | Min 8 characters |
| confirm_password | string | Yes | Must match password |
| role | string | Yes | One of: `investor`, `mitra` |
| cooperative_agreement | bool | Yes | Must be `true` |
| otp_token | string | Yes | Token from verify-otp step |

> **Note**: Registration is now simplified. Identity verification (KYC) and Bank Account details are collected in the mandatory **Complete Profile** step.

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "testuser",
      "role": "investor",
      "is_verified": false
    },
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc..."
  }
}
```

### 4. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "user@example.com",
    "password": "password123"
  }'
```

### 5. Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

---

## User Profile

### 1. Complete Profile (Mandatory)

After registration/login, user must complete their profile to transact.

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

**Response:**
```json
{
  "message": "Profile completed successfully. You can now access all features."
}
```

### 2. Update Wallet (For On-Chain Transparency)
Connect a crypto wallet (e.g. MetaMask) to associate IDR transactions with on-chain identity.

```bash
curl -X PUT http://localhost:8080/api/v1/user/wallet \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x123..."
  }'
```

### 3. Get Profile

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>"
```

### Update Profile

```bash
curl -X PUT http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe Updated",
    "phone": "081234567890",
    "country": "Indonesia"
  }'
```

---

## Profile Management (Flow: MANAGEMENT PROFIL USER)

### 1. Data Diri (Read-Only)

Get personal data from KYC - read only section.

```bash
curl -X GET http://localhost:8080/api/v1/user/profile/data \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "full_name": "John Doe",
  "nik_masked": "320101******7890",
  "email": "user@example.com",
  "username": "johndoe",
  "member_status": "calon_anggota_pendana",
  "role": "investor",
  "is_verified": true,
  "joined_at": "15 January 2024"
}
```

### 2. Rekening Bank

#### Get Current Bank Account

```bash
curl -X GET http://localhost:8080/api/v1/user/profile/bank-account \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "bank_code": "bca",
  "bank_name": "Bank Central Asia (BCA)",
  "account_number": "****7890",
  "account_name": "John Doe",
  "is_primary": true,
  "is_verified": true,
  "microcopy": "Rekening ini akan menjadi satu-satunya tujuan pencairan dana demi keamanan."
}
```

#### Change Bank Account (Requires OTP)

> ⚠️ **Security**: Changing bank account requires OTP verification for security.

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

#### Get Supported Banks

```bash
curl -X GET http://localhost:8080/api/v1/user/profile/banks \
  -H "Authorization: Bearer <access_token>"
```

### 3. Keamanan (Change Password)

```bash
curl -X PUT http://localhost:8080/api/v1/user/profile/password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "oldpassword123",
    "new_password": "newpassword456",
    "confirm_password": "newpassword456"
  }'
```

### 4. Log Out

Logout is handled client-side by removing stored tokens.

### Get Saldo/Balance

Get user's outstanding balance:
- **Investor**: Total dana yang sedang dipinjamkan ke mitra (outstanding investments)
- **Mitra**: Total hutang yang dimiliki ke investor-investor (outstanding debt)

> ⚠️ This is NOT a virtual wallet. It tracks active investments/debts only.

```bash
curl -X GET http://localhost:8080/api/v1/user/balance \
  -H "Authorization: Bearer <access_token>"
```

Response for Investor:
```json
{
  "balance": 5000000,
  "currency": "IDR",
  "description": "Total dana yang sedang dipinjamkan"
}
```

Response for Mitra:
```json
{
  "balance": 5000000,
  "currency": "IDR",
  "description": "Total saldo hutang ke investor"
}
```
```

---

## Fund Transfer via Escrow (MVP PROTOTYPE)

> 💡 **Note**: In MVP, fund transfers are simulated. In production:
> - **Investor deposits**: Transfer to escrow bank account
> - **Investor withdrawals**: Escrow transfers to investor's registered bank account
> - **Mitra receives funding**: Escrow transfers to mitra's registered bank account
> - **Mitra repays**: Mitra pays to VA, escrow distributes to investors' bank accounts

### Simulate Deposit (Admin/Testing)

```bash
curl -X POST http://localhost:8080/api/v1/payments/deposit \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000000
  }'
```

### Simulate Withdraw (Admin/Testing)

```bash
curl -X POST http://localhost:8080/api/v1/payments/withdraw \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000000
  }'
```

### Check Balance (Same as /user/balance)

```bash
curl -X GET http://localhost:8080/api/v1/payments/balance \
  -H "Authorization: Bearer <access_token>"
```

---

## MITRA Application (Flow 2)

> 📋 **Note**: Users who registered with `role: "mitra"` must complete MITRA verification before submitting invoices.

### Apply as MITRA

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

### Get MITRA Application Status

```bash
curl -X GET http://localhost:8080/api/v1/user/mitra/status \
  -H "Authorization: Bearer <access_token>"
```

### Upload MITRA Document

```bash
curl -X POST http://localhost:8080/api/v1/user/mitra/documents \
  -H "Authorization: Bearer <access_token>" \
  -F "document_type=nib" \
  -F "file=@/path/to/nib.pdf"
```

Document types: `nib`, `akta_pendirian`, `ktp_direktur`

---



## Invoices (Flow 4)

### Create Invoice

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

### List My Invoices (Exporter)

```bash
curl -X GET "http://localhost:8080/api/v1/invoices?page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

### Get Invoice by ID

```bash
curl -X GET http://localhost:8080/api/v1/invoices/<invoice_id> \
  -H "Authorization: Bearer <access_token>"
```

### Upload Invoice Document

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/documents \
  -H "Authorization: Bearer <access_token>" \
  -F "document_type=commercial_invoice" \
  -F "file=@/path/to/invoice.pdf"
```

Document types: `invoice_pdf`, `bill_of_lading`, `purchase_order`, `commercial_invoice`, `packing_list`, `certificate_of_origin`, `insurance`, `customs`, `other`

### Submit Invoice for Review

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/submit \
  -H "Authorization: Bearer <access_token>"
```

---

## Marketplace - Funding Pools (Flow 6)

> 👀 **Guest Access**: Guest users can view marketplace and pool details but cannot invest.

### List Open Pools

```bash
curl -X GET "http://localhost:8080/api/v1/pools?page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

### Get Pool Details

```bash
curl -X GET http://localhost:8080/api/v1/pools/<pool_id> \
  -H "Authorization: Bearer <access_token>"
```

Response includes tranche information:
```json
{
  "pool": {
    "priority_target": 40000000,
    "priority_funded": 0,
    "catalyst_target": 10000000,
    "catalyst_funded": 0,
    "priority_interest_rate": 10,
    "catalyst_interest_rate": 15,
    "deadline": "2024-02-01T00:00:00Z"
  },
  "priority_remaining": 40000000,
  "catalyst_remaining": 10000000,
  "priority_percentage_funded": 0,
  "catalyst_percentage_funded": 0
}
```

### Marketplace with Filters

```bash
curl -X GET "http://localhost:8080/api/v1/marketplace?grade=A&is_insured=true&page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

Query Parameters:
- `grade` - Filter by grade (A, B, or C)
- `is_insured` - Filter by insurance status (true/false)
- `min_amount` - Minimum pool amount
- `max_amount` - Maximum pool amount
- `page` - Page number
- `per_page` - Items per page

Response:
```json
{
  "pools": [
    {
      "pool": { ... },
      "remaining_amount": 25000000,
      "percentage_funded": 50,
      "priority_remaining": 20000000,
      "catalyst_remaining": 5000000,
      "priority_percentage_funded": 50,
      "catalyst_percentage_funded": 50,
      "grade": "A",
      "grade_score": 85,
      "is_insured": true,
      "buyer_country_risk": "Tier 1",
      "funding_progress": 50,
      "remaining_time": "5 hari lagi",
      "remaining_hours": 120,
      "priority_progress": 50,
      "catalyst_progress": 50
    }
  ],
  "total": 10,
  "page": 1,
  "per_page": 10,
  "total_pages": 1
}
```

#### Grading System

Invoice grades (A, B, C) are determined by:
1. **Buyer Country Risk** (40% weight)
   - Tier 1 (USA, Germany, Japan, etc.): Low Risk
   - Tier 2 (China, India, Indonesia, etc.): Medium Risk
   - Tier 3 (Others): High Risk

2. **Exporter History** (30% weight)
   - Repeat buyers get higher scores
   - MITRA-verified exporters get bonus points

3. **Document Completeness** (30% weight)
   - Commercial invoice, Bill of Lading, Certificate of Origin, etc.

---

## Risk Questionnaire (Investor Only)

> ⚠️ **Important**: Investors MUST complete the risk questionnaire to unlock the Catalyst (Junior) tranche. Without completing the questionnaire correctly, investors can only invest in Priority (Senior) tranche.

### Get Questions

```bash
curl -X GET http://localhost:8080/api/v1/risk-questionnaire/questions \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "questions": [
    {
      "id": "investment_purpose",
      "question": "Apa tujuan investasi Anda?",
      "options": [
        {"value": 1, "label": "Dana darurat"},
        {"value": 2, "label": "Tabungan jangka panjang"},
        {"value": 3, "label": "Spekulasi/mencari keuntungan tinggi"}
      ]
    },
    {
      "id": "loss_tolerance",
      "question": "Apakah Anda siap kehilangan 100% modal pada investasi berisiko tinggi?",
      "options": [
        {"value": 1, "label": "Tidak, saya tidak siap"},
        {"value": 2, "label": "Ya, saya siap dan memahami risikonya"}
      ]
    },
    {
      "id": "tranche_understanding",
      "question": "Apakah Anda memahami bahwa Junior Tranche (Catalyst) dibayar SETELAH Senior Tranche (Priority)?",
      "options": [
        {"value": 1, "label": "Tidak, saya belum mengerti"},
        {"value": 2, "label": "Ya, saya mengerti dan menerima risiko ini"}
      ]
    }
  ]
}
```

### Submit Risk Questionnaire

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

> 💡 **Note**: Risk questionnaire is now **informational only**. Catalyst eligibility is determined by inline consent at each investment, not by pre-completing this questionnaire.

Response:
```json
{
  "completed": true,
  "message": "Questionnaire saved. You can invest in any tranche by accepting the required consents at investment time."
}
```

### Get Questionnaire Status

```bash
curl -X GET http://localhost:8080/api/v1/risk-questionnaire/status \
  -H "Authorization: Bearer <access_token>"
```

---

## Investments (Flow 6)

> 📊 **Interest Calculation**: Flat rate - `Expected Return = Principal + (Principal × Rate/100)`
>
> Example: Rp 5,000,000 at 10% = Rp 5,000,000 + Rp 500,000 = Rp 5,500,000
>
> 🏦 **Bank Account Required**: Investor must have a verified bank account for receiving returns.

### Investment UI Flow

**1. Header Info**: Detail Buyer, Dokumen Legal (Download/Preview), & Skema Asuransi

**2. Tab Selector**:
- Tab Kiri (Default): [ 🛡️ Pendanaan Prioritas ]
- Tab Kanan: [ ⚡ Pendanaan Katalis ]

**3. Info Box** (berubah sesuai Tab):
- **Prioritas**: "Dana Anda berada di antrean pertama. Saat Importir membayar, Anda akan menerima pengembalian modal & hasil paling awal. Risikonya lebih rendah karena dilindungi oleh dana Katalis." (Warna: Biru/Hijau)
- **Katalis**: "Dana Anda berfungsi sebagai penopang risiko. Anda akan dibayar setelah Pendana Prioritas lunas sepenuhnya. Sebagai ganti risiko ini, Anda mendapatkan imbal hasil lebih tinggi." (Warna: Oranye/Ungu)

**4. Calculator Input**:
- Field: "Nominal Pendanaan (Rp)"
- Text: "Estimasi Imbal Hasil (..% p.a)" - auto-update sesuai Tab
- Text: "Estimasi Total Diterima" - Rumus: Modal + (Modal × Rate)
- Button: "Danai Sekarang"

### Invest in Pool (with Inline Consent)

**Priority Tranche Flow:**
1. Muncul Bottom Sheet ringkasan
2. Checkbox: "Saya menyetujui Syarat & Ketentuan"
3. Input PIN/Biometric
4. Submit

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

**Catalyst Tranche Flow:**
1. Muncul Full Screen Warning Modal (Warna Merah/Oranye)
2. Wajib Scroll sampai bawah
3. Wajib Centang 3 Poin:
   - [ ] "Saya sadar dana ini menjadi jaminan pertama jika gagal bayar."
   - [ ] "Saya siap menanggung risiko kehilangan modal."
   - [ ] "Saya paham ini bukan produk bank."
4. Tombol "Lanjut Danai" baru menyala
5. Input PIN/Biometric
6. Submit

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
    }'
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| pool_id | uuid | Yes | Target funding pool ID |
| amount | number | Yes | Investment amount (Rp) |
| tranche | string | Yes | `priority` or `catalyst` |
| tnc_accepted | bool | Yes | Terms & Conditions accepted |
| catalyst_consents | object | Catalyst only | Required for catalyst tranche |
| catalyst_consents.first_loss_consent | bool | Catalyst only | "Saya sadar dana ini menjadi jaminan pertama jika gagal bayar" |
| catalyst_consents.risk_loss_consent | bool | Catalyst only | "Saya siap menanggung risiko kehilangan modal" |
| catalyst_consents.not_bank_consent | bool | Catalyst only | "Saya paham ini bukan produk bank" |

### List My Investments

```bash
curl -X GET "http://localhost:8080/api/v1/investments?page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

### Get Portfolio Summary (Flow 9)

```bash
curl -X GET http://localhost:8080/api/v1/investments/portfolio \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "total_funding": 6000000,
  "total_expected_gain": 150000,
  "total_realized_gain": 0,
  "priority_allocation": 5000000,
  "catalyst_allocation": 1000000,
  "active_investments": 2,
  "completed_deals": 0
}
```

---

## Mitra Dashboard (Flow 8)

### Get Mitra Dashboard

```bash
curl -X GET http://localhost:8080/api/v1/mitra/dashboard \
  -H "Authorization: Bearer <access_token>"
```

---

## Mitra Repayment (Flow: MITRA MEMBAYAR HUTANG)

### UI Flow Overview

1. **Active Invoice Card**: Shows "Bayar Pelunasan" primary button
2. **Modal Pop-up**: Displays repayment breakdown (Principal + Interest by tranche)
3. **Payment Method Selection**: Choose VA bank (BCA/Mandiri/BNI)
4. **VA Payment Page**:
   - VA Number (copy button)
   - Total Amount (closed amount, cannot change)
   - Timer: "Selesaikan dalam 24:00:00"
5. **Auto-refresh**: Status becomes "LUNAS" when backend detects payment

### Get Active Invoices for Repayment

```bash
curl -X GET http://localhost:8080/api/v1/mitra/invoices/active \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "invoices": [
    {
      "invoice_id": "uuid",
      "pool_id": "uuid",
      "invoice_number": "INV-2024-001",
      "buyer_name": "German Import GmbH",
      "amount": 50000000,
      "total_due": 55500000,
      "due_date": "2024-03-01T00:00:00Z",
      "status": "funded",
      "days_until_due": 30,
      "can_pay": true
    }
  ]
}
```

### Get Repayment Breakdown

Shows detailed breakdown by tranche (Priority + Catalyst).

```bash
curl -X GET http://localhost:8080/api/v1/mitra/pools/<pool_id>/breakdown \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "pool_id": "uuid",
  "invoice_id": "uuid",
  "invoice_number": "INV-2024-001",
  "buyer_name": "German Import GmbH",
  "due_date": "2024-03-01T00:00:00Z",
  
  "priority_principal": 40000000,
  "catalyst_principal": 10000000,
  "total_principal": 50000000,
  
  "priority_interest_rate": 10,
  "catalyst_interest_rate": 15,
  "priority_interest": 4000000,
  "catalyst_interest": 1500000,
  "total_interest": 5500000,
  
  "priority_total": 44000000,
  "catalyst_total": 11500000,
  "grand_total": 55500000,
  
  "currency": "IDR"
}
```

### Get Available VA Payment Methods

```bash
curl -X GET http://localhost:8080/api/v1/mitra/payment-methods \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
[
  {"bank_code": "bca", "bank_name": "Bank Central Asia (BCA)"},
  {"bank_code": "mandiri", "bank_name": "Bank Mandiri"},
  {"bank_code": "bni", "bank_name": "Bank Negara Indonesia (BNI)"}
]
```

### Create VA for Payment

```bash
curl -X POST http://localhost:8080/api/v1/mitra/repayment/va \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "bank_code": "bca"
  }'
```

Response:
```json
{
  "virtual_account": {
    "id": "uuid",
    "va_number": "8bca123456789",
    "bank_code": "bca",
    "bank_name": "Bank Central Asia (BCA)",
    "amount": 55500000,
    "status": "pending",
    "expires_at": "2024-01-16T10:00:00Z"
  },
  "breakdown": { ... },
  "remaining_time": "23:59:59",
  "remaining_hours": 24,
  "microcopy": "Selesaikan pembayaran dalam waktu 24 jam. VA akan otomatis kadaluarsa setelah batas waktu."
}
```

### Get VA Payment Page Details

```bash
curl -X GET http://localhost:8080/api/v1/mitra/repayment/va/<va_id> \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "va_number": "88001234567890",
  "bank_code": "bca",
  "bank_name": "Bank Central Asia (BCA)",
  "amount": 55500000,
  "amount_formatted": "Rp 55.500.000",
  "status": "pending",
  "expires_at": "2024-01-16T10:00:00Z",
  "remaining_time": "23:45:30",
  "breakdown": { ... },
  "microcopy": "Nominal pembayaran bersifat tetap dan tidak dapat diubah."
}
```

### Simulate VA Payment (MVP/Testing)

For testing purposes only - simulates receiving VA payment.

```bash
curl -X POST http://localhost:8080/api/v1/mitra/repayment/va/<va_id>/simulate-pay \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "status": "paid",
  "message": "Pembayaran berhasil! Dana sedang didistribusikan ke investor.",
  "va_id": "uuid",
  "paid_at": "2024-01-15T12:00:00Z"
}
```

---

## Exporter Disbursement to Investors (Flow 11)

When the exporter receives payment from the importer, they pay back to investors through escrow.

### Disburse to Investors

```bash
curl -X POST http://localhost:8080/api/v1/exporter/disbursement \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "total_amount": 55000000
  }'
```

Response:
```json
{
  "pool_id": "<pool_uuid>",
  "total_disbursed": 55000000,
  "priority_disbursed": 44000000,
  "catalyst_disbursed": 11000000,
  "investor_results": [
    {
      "investor_id": "<uuid>",
      "tranche": "priority",
      "principal": 5000000,
      "expected_return": 5500000,
      "actual_return": 5500000,
      "status": "full"
    },
    {
      "investor_id": "<uuid>",
      "tranche": "catalyst",
      "principal": 1000000,
      "expected_return": 1150000,
      "actual_return": 1150000,
      "status": "full"
    }
  ],
  "status": "completed",
  "message": "Disbursement completed. Priority investors paid first, then Catalyst investors."
}
```

#### Payment Priority Logic

1. **Priority Tranche First**: All Priority investors are paid their full expected return first
2. **Catalyst Tranche Second**: Remaining funds go to Catalyst investors
3. **Partial Payment**: If funds are insufficient:
   - Priority investors may receive partial payment (status: "partial")
   - Catalyst investors only receive funds if Priority is fully paid
   - Catalyst investors may receive partial or zero payment if funds are depleted

Status values:
- `full` - Investor received full expected return
- `partial` - Investor received partial payment
- `zero` - Investor received nothing (in case of severe shortfall)

---

## Admin Endpoints

### Get Pending KYC

```bash
curl -X GET "http://localhost:8080/api/v1/admin/kyc/pending?page=1&per_page=10" \
  -H "Authorization: Bearer <admin_access_token>"
```

### Approve KYC

```bash
curl -X POST http://localhost:8080/api/v1/admin/kyc/<kyc_id>/approve \
  -H "Authorization: Bearer <admin_access_token>"
```

### Reject KYC

```bash
curl -X POST http://localhost:8080/api/v1/admin/kyc/<kyc_id>/reject \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Dokumen tidak valid"
  }'
```

### Approve Invoice (Flow 5)

```bash
curl -X POST http://localhost:8080/api/v1/admin/invoices/<invoice_id>/approve \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "interest_rate": 12.0
  }'
```

### Reject Invoice

```bash
curl -X POST http://localhost:8080/api/v1/admin/invoices/<invoice_id>/reject \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Dokumen invoice tidak lengkap"
  }'
```

### Tokenize Invoice

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/tokenize \
  -H "Authorization: Bearer <admin_access_token>"
```

### Create Funding Pool

```bash
curl -X POST http://localhost:8080/api/v1/invoices/<invoice_id>/pool \
  -H "Authorization: Bearer <admin_access_token>"
```

### Disburse to Exporter (Flow 7)

```bash
curl -X POST http://localhost:8080/api/v1/admin/pools/<pool_id>/disburse \
  -H "Authorization: Bearer <admin_access_token>"
```

### Close Pool and Notify Exporter

When pool funding deadline ends, admin closes the pool and sends payment notification to exporter via email.

```bash
curl -X POST http://localhost:8080/api/v1/admin/pools/<pool_id>/close \
  -H "Authorization: Bearer <admin_access_token>"
```

Response:
```json
{
  "message": "Pool closed and payment notification sent to exporter",
  "notification_data": {
    "invoice_id": "<uuid>",
    "invoice_number": "INV-2024-001",
    "exporter_name": "PT Eksportir Indonesia",
    "buyer_name": "German Import GmbH",
    "buyer_email": "buyer@german-import.de",
    "principal_amount": 50000000,
    "total_interest": 5000000,
    "total_amount_due": 55000000,
    "currency": "IDR",
    "due_date": "2024-03-01T00:00:00Z",
    "investor_details": [
      {
        "tranche": "priority",
        "count": 5,
        "total_principal": 40000000,
        "total_interest": 4000000
      },
      {
        "tranche": "catalyst",
        "count": 3,
        "total_principal": 10000000,
        "total_interest": 1500000
      }
    ],
    "payment_id": "<uuid>",
    "payment_link": "https://vessel.app/pay/<payment_id>"
  }
}
```

Note: This endpoint:
1. Changes pool status from "open" to "closed"
2. Generates payment information for the importer
3. Sends email notification to exporter with payment details
4. Creates importer payment record for tracking

### Process Repayment (Flow 11)

```bash
curl -X POST http://localhost:8080/api/v1/admin/invoices/<invoice_id>/repay \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 55000000
  }'
```

Note: Repayment follows priority-first rule - Priority investors are paid first, then Catalyst investors.

### Get Pending MITRA Applications

```bash
curl -X GET "http://localhost:8080/api/v1/admin/mitra/pending?page=1&per_page=10" \
  -H "Authorization: Bearer <admin_access_token>"
```

### Approve MITRA Application

```bash
curl -X POST http://localhost:8080/api/v1/admin/mitra/<application_id>/approve \
  -H "Authorization: Bearer <admin_access_token>"
```

### Reject MITRA Application

```bash
curl -X POST http://localhost:8080/api/v1/admin/mitra/<application_id>/reject \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Dokumen tidak valid atau tidak lengkap"
  }'
```

### Grant Balance to User (MVP)

```bash
curl -X POST http://localhost:8080/api/v1/admin/balance/grant \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<user_uuid>",
    "amount": 10000000
  }'
```

---

## Tranche System Explained

### Priority Tranche (Pendanaan Prioritas)
- Lower risk, lower return
- Interest rate: Typically 10-12%
- Paid FIRST during repayment
- Default: 80% of total pool

### Catalyst Tranche (Pendanaan Katalis)
- Higher risk, higher return
- Interest rate: Typically 15-18%
- Paid AFTER Priority is fully paid
- Default: 20% of total pool
- First-loss capital in case of default

### Interest Calculation (Flat Rate)

For MVP, interest is calculated as a **flat rate** (no time factor):

```
Expected Return = Principal + (Principal × Rate / 100)
```

**Example - Priority Tranche (10% rate):**
- Investment: Rp 5,000,000
- Expected Return: Rp 5,000,000 + (Rp 5,000,000 × 10/100) = Rp 5,500,000
- Profit: Rp 500,000

**Example - Catalyst Tranche (15% rate):**
- Investment: Rp 2,000,000
- Expected Return: Rp 2,000,000 + (Rp 2,000,000 × 15/100) = Rp 2,300,000
- Profit: Rp 300,000

---

## On-Chain Transparency

The following data is recorded on-chain for transparency (Lisk blockchain):

### What's Recorded On-Chain

| Data Type | Smart Contract | Event | Description |
|-----------|----------------|-------|-------------|
| Invoice Minting | InvoiceNFT | `InvoiceMinted` | When invoice is tokenized as NFT |
| Invoice Status | InvoiceNFT | `InvoiceStatusChanged` | Status changes: Active, Funded, Repaid, Defaulted |
| Pool Creation | InvoicePool | `PoolCreated` | When funding pool is opened |
| Investment | InvoicePool | `InvestmentRecorded` | Each investor's contribution |
| Pool Filled | InvoicePool | `PoolFilled` | When pool reaches target |
| Disbursement | InvoicePool | `DisbursementRecorded` | Funds sent to mitra |
| Repayment | InvoicePool | `RepaymentRecorded` | Importer payment received |
| Investor Returns | InvoicePool | `InvestorReturnRecorded` | Returns distributed to each investor |
| Mitra Credit | InvoicePool | `MitraBalanceCredited` | Excess funds credited to mitra (partial funding) |

### What Remains Off-Chain

- Payment input (amount) - entered manually
- Bank account details
- KYC/Identity verification
- Document storage (files on IPFS, references on-chain)

---

## Partial Funding Scenario

When an invoice is not fully funded but importer pays the full invoice amount:

### Example

| Item | Amount |
|------|--------|
| Invoice Amount | Rp 100,000,000 |
| Funded Amount | Rp 10,000,000 (only 10% funded) |
| Expected Investor Returns | Rp 11,000,000 (10% + 10% interest) |
| Importer Pays | Rp 100,000,000 (full invoice) |
| Platform Fee (2%) | Rp 2,000,000 |
| **Remaining After Fee** | Rp 98,000,000 |
| **Investor Returns** | Rp 11,000,000 |
| **Excess to Mitra Balance** | Rp 87,000,000 |

### Flow

1. Importer pays full invoice amount via payment link
2. Platform fee is deducted
3. Investors receive their returns (priority-first)
4. Excess amount is credited to mitra's balance
5. All transactions recorded on-chain for transparency

### Transaction Types

- `investor_return` - Returns paid to investors
- `repayment_excess` - Excess funds credited to mitra balance

---

## MVP Abstractions

For the MVP phase, the following features are abstracted:

### 1. Payment Gateway (Abstracted)
- No real payment gateway integration
- Balance system is simulated via admin grants
- Deposits/withdrawals are prototype endpoints

### 2. Currency
- All transactions use **IDR** (Indonesian Rupiah)
- No blockchain IDRX stablecoin integration for MVP
- Currency conversion from USD/EUR to IDR is supported

### 3. Escrow System
- Escrow is simulated in the backend
- No actual smart contract escrow for MVP
- Fund flow is tracked in database

### 4. Blockchain Integration
- On-chain transparency is implemented via events
- Smart contract addresses configured in environment
- Transaction hashes returned for verification

---

## Error Responses

All errors return JSON with the following format:

```json
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Error description"
  }
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error
