package mailer

// NoopMailer discards all emails. Use in tests or when mail is disabled.
type NoopMailer struct{}

func (n *NoopMailer) SendWelcome(_, _ string) error { return nil }
