package mailer

import (
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/khanhnp-2797/echo-realworld-api/internal/config"
)

type smtpMailer struct {
	dialer *gomail.Dialer
	from   string
}

// NewSMTPMailer creates a Mailer backed by an SMTP server via gomail.v2.
// For local development, point cfg at MailHog (default: localhost:1025).
// When MAIL_USERNAME is empty (e.g. MailHog), authentication is skipped.
func NewSMTPMailer(cfg config.MailConfig) Mailer {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	// MailHog does not support TLS — disable it for local dev.
	// In production set MAIL_USERNAME/PASSWORD and gomail will use STARTTLS.
	if cfg.Username == "" {
		d.SSL = false
	}
	return &smtpMailer{dialer: d, from: cfg.From}
}

func (m *smtpMailer) SendWelcome(toEmail, username string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Welcome to RealWorld!")
	msg.SetBody("text/html", welcomeHTML(username, toEmail))
	msg.AddAlternative("text/plain", welcomePlain(username, toEmail))

	return m.dialer.DialAndSend(msg)
}

func welcomePlain(username, email string) string {
	return fmt.Sprintf(`Hi %s,

Welcome to RealWorld! Your account has been created successfully.

Email: %s

Happy reading and writing!

— The RealWorld Team`, username, email)
}

func welcomeHTML(username, email string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:480px;margin:40px auto;color:#333">
  <h2 style="color:#5c6bc0">Welcome to RealWorld, %s!</h2>
  <p>Your account has been created successfully.</p>
  <table style="border-collapse:collapse;margin:16px 0">
    <tr><td style="padding:4px 8px;font-weight:bold">Email</td><td style="padding:4px 8px">%s</td></tr>
  </table>
  <p>Happy reading and writing!</p>
  <hr style="border:none;border-top:1px solid #eee;margin:24px 0"/>
  <p style="font-size:12px;color:#999">— The RealWorld Team</p>
</body>
</html>`, username, email)
}
