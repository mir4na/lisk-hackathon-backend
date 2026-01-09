# VESSEL Backend API Documentation

Base URL: `http://localhost:8080`

## Important Notes

> ⚠️ **Wallet Requirement**: Users (both eksportir and investor) MUST connect their wallet before creating invoices or investing. Use `PUT /api/v1/user/wallet` to connect wallet.

> 💰 **Account Balance**: Each user has an account balance (balance_idr) visible via `GET /api/v1/user/balance` or in profile response.

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

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123",
    "confirm_password": "password123",
    "role": "investor",
    "full_name": "John Doe",
    "phone_number": "081234567890",
    "cooperative_agreement": true,
    "otp_token": "<token_from_verify_otp>"
  }'
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

### Get Profile

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

### Update Wallet

```bash
curl -X PUT http://localhost:8080/api/v1/user/wallet \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x1234567890123456789012345678901234567890"
  }'
```

---

## Payment Gateway (PROTOTYPE)

### Deposit Saldo

```bash
curl -X POST http://localhost:8080/api/v1/payments/deposit \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000000
  }'
```

### Withdraw Saldo

```bash
curl -X POST http://localhost:8080/api/v1/payments/withdraw \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000000
  }'
```

### Check Balance

```bash
curl -X GET http://localhost:8080/api/v1/payments/balance \
  -H "Authorization: Bearer <access_token>"
```

---

## MITRA Application (Flow 2)

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

## Buyers (Importir)

### Create Buyer

```bash
curl -X POST http://localhost:8080/api/v1/buyers \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "German Import GmbH",
    "country": "Germany",
    "contact_email": "buyer@german-import.de",
    "contact_phone": "+49123456789"
  }'
```

### List Buyers

```bash
curl -X GET "http://localhost:8080/api/v1/buyers?page=1&per_page=10" \
  -H "Authorization: Bearer <access_token>"
```

### Get Buyer by ID

```bash
curl -X GET http://localhost:8080/api/v1/buyers/<buyer_id> \
  -H "Authorization: Bearer <access_token>"
```

---

## Invoices (Flow 4)

### Create Invoice

```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "buyer_id": "<buyer_uuid>",
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

Response (if eligible for Catalyst):
```json
{
  "catalyst_unlocked": true,
  "message": "Selamat! Anda sekarang dapat berinvestasi di Junior Tranche (Catalyst)."
}
```

Response (if not eligible):
```json
{
  "catalyst_unlocked": false,
  "message": "Berdasarkan jawaban Anda, Anda hanya dapat berinvestasi di Senior Tranche (Priority)."
}
```

### Get Questionnaire Status

```bash
curl -X GET http://localhost:8080/api/v1/risk-questionnaire/status \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "completed": true,
  "catalyst_unlocked": true,
  "completed_at": "2024-01-15T10:30:00Z"
}
```

#### Catalyst Unlock Logic

To unlock Catalyst tranche, investor must answer:
1. **Investment Purpose**: "Spekulasi" (value: 3)
2. **Loss Tolerance**: "Ya, saya siap" (value: 2)
3. **Tranche Understanding**: "Ya, saya mengerti" (value: 2)

---

## Investments (Flow 6)

### Invest in Pool (with Tranche Selection)

> ⚠️ **Wallet Required**: Investor must have a connected wallet address.
>
> ⚠️ **Catalyst Tranche**: Requires completing risk questionnaire with correct answers. If not unlocked, you will receive error: "Catalyst tranche not unlocked. Please complete risk questionnaire first."

**Priority Tranche (Lower Risk, Lower Return):**

```bash
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "amount": 5000000,
    "tranche": "priority"
  }'
```

**Catalyst Tranche (Higher Risk, Higher Return):**

```bash
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": "<pool_uuid>",
    "amount": 1000000,
    "tranche": "catalyst"
  }'
```

Note: Catalyst tranche requires completing risk questionnaire first (see Risk Questionnaire section).

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
    "currency": "IDRX",
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

---

## Tranche System Explained

### Priority Tranche (Pendanaan Prioritas)
- Lower risk, lower return
- Interest rate: Typically 10-12% p.a.
- Paid FIRST during repayment
- Default: 80% of total pool

### Catalyst Tranche (Pendanaan Katalis)
- Higher risk, higher return
- Interest rate: Typically 15-18% p.a.
- Paid AFTER Priority is fully paid
- Default: 20% of total pool
- First-loss capital in case of default

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
