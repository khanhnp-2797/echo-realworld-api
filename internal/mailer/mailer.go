package mailer

// Mailer is the interface for sending transactional emails.
type Mailer interface {
	// SendWelcome sends a welcome email to a newly registered user.
	SendWelcome(toEmail, username string) error
}
