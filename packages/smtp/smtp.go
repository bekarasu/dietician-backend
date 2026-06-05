package smtp

import (
	"fmt"
	"net/smtp"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Sender struct {
	cfg Config
}

func NewSender(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

// Send sends a plain-text email. It is a no-op when Host is empty.
func (s *Sender) Send(to, subject, body string) error {
	if s.cfg.Host == "" {
		return nil
	}
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	return smtp.SendMail(addr, auth, s.cfg.Username, []string{to}, []byte(msg))
}
