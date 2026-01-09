package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/models"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendOTPEmail sends an OTP code to the specified email address
func (s *EmailService) SendOTPEmail(email, code, purpose string) error {
	subject := s.getOTPSubject(purpose)
	body, err := s.getOTPEmailBody(code, purpose)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(email, subject, body)
}

// SendMitraApprovalEmail sends approval notification to MITRA applicant
func (s *EmailService) SendMitraApprovalEmail(email, companyName string) error {
	subject := "Congratulations! Your MITRA Application is Approved - VESSEL"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #2563eb;">Congratulations!</h2>
				<p>Your MITRA application for <strong>%s</strong> has been approved.</p>
				<p>You can now:</p>
				<ul>
					<li>Create invoice funding requests</li>
					<li>Access MITRA features on your dashboard</li>
					<li>Receive funding from investors</li>
				</ul>
				<p>Please login to your account to get started.</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #666; font-size: 12px;">
					This email was sent by VESSEL Platform.<br>
					If you have any questions, contact support@vessel.id
				</p>
			</div>
		</body>
		</html>
	`, companyName)

	return s.sendEmail(email, subject, body)
}

// SendMitraRejectionEmail sends rejection notification with reason
func (s *EmailService) SendMitraRejectionEmail(email, companyName, reason string) error {
	subject := "Your MITRA Application Status - VESSEL"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #dc2626;">MITRA Application Rejected</h2>
				<p>We regret to inform you that the MITRA application for <strong>%s</strong> could not be approved at this time.</p>
				<div style="background-color: #fef2f2; border-left: 4px solid #dc2626; padding: 15px; margin: 20px 0;">
					<strong>Rejection Reason:</strong>
					<p style="margin: 10px 0 0 0;">%s</p>
				</div>
				<p>You may reapply after addressing the issues mentioned above.</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #666; font-size: 12px;">
					This email was sent by VESSEL Platform.<br>
					If you have any questions, contact support@vessel.id
				</p>
			</div>
		</body>
		</html>
	`, companyName, reason)

	return s.sendEmail(email, subject, body)
}

// SendInvoiceApprovalEmail sends notification when invoice is approved
func (s *EmailService) SendInvoiceApprovalEmail(email, invoiceNumber, grade string, priorityRate, catalystRate float64) error {
	subject := fmt.Sprintf("Invoice %s Approved - VESSEL", invoiceNumber)
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #16a34a;">Invoice Approved</h2>
				<p>Invoice <strong>%s</strong> has been approved and is ready for funding.</p>
				<div style="background-color: #f0fdf4; border-left: 4px solid #16a34a; padding: 15px; margin: 20px 0;">
					<p><strong>Grade:</strong> %s</p>
					<p><strong>Priority Rate:</strong> %.2f%% p.a.</p>
					<p><strong>Catalyst Rate:</strong> %.2f%% p.a.</p>
				</div>
				<p>Your invoice is now available on the marketplace for investors.</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #666; font-size: 12px;">
					This email was sent by VESSEL Platform.
				</p>
			</div>
		</body>
		</html>
	`, invoiceNumber, grade, priorityRate, catalystRate)

	return s.sendEmail(email, subject, body)
}

// SendImporterPaymentNotification sends payment notification to importer (buyer)
// This is for non-users to pay invoice via payment ID
func (s *EmailService) SendImporterPaymentNotification(email string, data *models.PaymentNotificationData) error {
	subject := fmt.Sprintf("Invoice Payment Request - %s - VESSEL", data.InvoiceNumber)
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #2563eb;">VESSEL Invoice Payment Request</h2>
				<p>Dear %s,</p>
				<p>You have received an invoice payment request from <strong>%s</strong>.</p>
				
				<div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin: 20px 0;">
					<p><strong>Invoice Number:</strong> %s</p>
					<p><strong>Amount Due:</strong> %s %.2f</p>
					<p><strong>Due Date:</strong> %s</p>
				</div>
				
				<div style="background-color: #eff6ff; border-left: 4px solid #2563eb; padding: 15px; margin: 20px 0;">
					<p><strong>Payment ID:</strong></p>
					<p style="font-size: 18px; font-family: monospace; background-color: #fff; padding: 10px; border-radius: 4px;">%s</p>
				</div>
				
				<p>To complete the payment, please use the following link:</p>
				<p><a href="%s" style="display: inline-block; background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Pay Now</a></p>
				
				<p style="color: #666; font-size: 14px; margin-top: 20px;">
					If you have any questions, please contact the exporter directly.
				</p>
				
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #666; font-size: 12px;">
					This email was sent by VESSEL Platform - Invoice Factoring for Trade Finance.<br>
					This is an automated message, please do not reply.
				</p>
			</div>
		</body>
		</html>
	`, data.BuyerName, data.ExporterName, data.InvoiceNumber, data.Currency, data.AmountDue,
		data.DueDate.Format("02 January 2006"), data.PaymentID, data.PaymentLink)

	return s.sendEmail(email, subject, body)
}

// SendExporterPaymentNotification sends payment notification to exporter
// when funding pool ends - exporter will forward this to importer
func (s *EmailService) SendExporterPaymentNotification(email string, data *models.ExporterPaymentNotificationData) error {
	subject := fmt.Sprintf("Invoice Payment Details Generated - %s - VESSEL", data.InvoiceNumber)

	// Build investor details table
	investorRows := ""
	for i, inv := range data.InvestorDetails {
		investorRows += fmt.Sprintf(`
			<tr>
				<td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">%d</td>
				<td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">%s</td>
				<td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">%s %.2f</td>
				<td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">%.2f%%</td>
				<td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">%s %.2f</td>
			</tr>
		`, i+1, inv.Tranche, data.Currency, inv.Amount, inv.InterestRate, data.Currency, inv.ExpectedReturn)
	}

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 700px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #2563eb;">VESSEL - Invoice Payment Details</h2>
				<p>Dear %s,</p>
				<p>The funding pool for your invoice has ended. Below are the payment details that need to be paid by the importer (<strong>%s</strong>).</p>

				<div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin: 20px 0;">
					<h3 style="margin-top: 0; color: #1f2937;">Invoice Summary</h3>
					<table style="width: 100%%; border-collapse: collapse;">
						<tr>
							<td style="padding: 8px 0;"><strong>Invoice Number:</strong></td>
							<td>%s</td>
						</tr>
						<tr>
							<td style="padding: 8px 0;"><strong>Principal Amount:</strong></td>
							<td>%s %.2f</td>
						</tr>
						<tr>
							<td style="padding: 8px 0;"><strong>Total Interest (Investor):</strong></td>
							<td>%s %.2f</td>
						</tr>
						<tr>
							<td style="padding: 8px 0;"><strong>Platform Fee (2%%):</strong></td>
							<td>%s %.2f</td>
						</tr>
						<tr style="background-color: #dbeafe;">
							<td style="padding: 8px;"><strong>Total Amount Due:</strong></td>
							<td style="padding: 8px;"><strong>%s %.2f</strong></td>
						</tr>
						<tr>
							<td style="padding: 8px 0;"><strong>Due Date:</strong></td>
							<td>%s</td>
						</tr>
					</table>
				</div>

				<div style="background-color: #fef3c7; border-left: 4px solid #f59e0b; padding: 15px; margin: 20px 0;">
					<p style="margin: 0;"><strong>Payment ID:</strong></p>
					<p style="font-size: 18px; font-family: monospace; background-color: #fff; padding: 10px; border-radius: 4px; margin: 10px 0;">%s</p>
				</div>

				<h3 style="color: #1f2937;">Investor Breakdown</h3>
				<table style="width: 100%%; border-collapse: collapse; background-color: #fff; border: 1px solid #e5e7eb;">
					<thead style="background-color: #f9fafb;">
						<tr>
							<th style="padding: 12px 8px; text-align: left; border-bottom: 2px solid #e5e7eb;">#</th>
							<th style="padding: 12px 8px; text-align: left; border-bottom: 2px solid #e5e7eb;">Tranche</th>
							<th style="padding: 12px 8px; text-align: left; border-bottom: 2px solid #e5e7eb;">Investment</th>
							<th style="padding: 12px 8px; text-align: left; border-bottom: 2px solid #e5e7eb;">Interest Rate</th>
							<th style="padding: 12px 8px; text-align: left; border-bottom: 2px solid #e5e7eb;">Expected Return</th>
						</tr>
					</thead>
					<tbody>
						%s
					</tbody>
				</table>

				<div style="background-color: #eff6ff; border-radius: 8px; padding: 20px; margin: 20px 0;">
					<h4 style="margin-top: 0;">Importer Information</h4>
					<p><strong>Company:</strong> %s</p>
					<p><strong>Email:</strong> %s</p>
					<p style="color: #666; font-size: 14px;">
						Please forward this payment information to the importer for settlement before the due date.
					</p>
				</div>

				<p>Payment Link: <a href="%s" style="color: #2563eb;">%s</a></p>

				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #666; font-size: 12px;">
					This email was sent by VESSEL Platform - Invoice Factoring for Trade Finance.<br>
					This is an automated message, please do not reply.
				</p>
			</div>
		</body>
		</html>
	`, data.ExporterName, data.BuyerName,
		data.InvoiceNumber,
		data.Currency, data.PrincipalAmount,
		data.Currency, data.TotalInterest,
		data.Currency, data.PlatformFee,
		data.Currency, data.TotalAmountDue,
		data.DueDate.Format("02 January 2006"),
		data.PaymentID,
		investorRows,
		data.BuyerName, data.BuyerEmail,
		data.PaymentLink, data.PaymentLink)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) getOTPSubject(purpose string) string {
	switch purpose {
	case "registration":
		return "Registration Verification Code - VESSEL"
	case "login":
		return "Login Verification Code - VESSEL"
	case "password_reset":
		return "Password Reset Code - VESSEL"
	default:
		return "Verification Code - VESSEL"
	}
}

func (s *EmailService) getOTPEmailBody(code, purpose string) (string, error) {
	tmpl := `
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #2563eb;">VESSEL</h2>
				<p>{{.Message}}</p>
				<div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; text-align: center; margin: 20px 0;">
					<span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #1f2937;">{{.Code}}</span>
				</div>
				<p style="color: #666; font-size: 14px;">
					This code is valid for {{.ExpiryMinutes}} minutes.<br>
					Do not share this code with anyone.
				</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #666; font-size: 12px;">
					If you did not request this code, please ignore this email.<br>
					This email was sent automatically by VESSEL Platform.
				</p>
			</div>
		</body>
		</html>
	`

	var message string
	switch purpose {
	case "registration":
		message = "Thank you for registering with VESSEL. Use the following code to verify your email:"
	case "login":
		message = "Use the following code to continue logging into your VESSEL account:"
	case "password_reset":
		message = "Use the following code to reset your VESSEL account password:"
	default:
		message = "Use the following code for verification:"
	}

	t, err := template.New("otp").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := struct {
		Code          string
		Message       string
		ExpiryMinutes int
	}{
		Code:          code,
		Message:       message,
		ExpiryMinutes: s.config.OTPExpiryMinutes,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	// For hackathon MVP, if SMTP is not configured, just log and return
	if s.config.SMTPUsername == "" || s.config.SMTPPassword == "" {
		fmt.Printf("[EMAIL] Would send to: %s\nSubject: %s\n", to, subject)
		return nil
	}

	from := s.config.SMTPFrom
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	// Compose email
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Gmail SMTP auth
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Send email
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
