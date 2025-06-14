package smtp

import (
	"log"
	"net/smtp"
)

type Smtp struct {
	Username string
	Password string
}

type SmtpAdapter interface {
	SendEmail(to, subject, body string) error
}

var _ SmtpAdapter = (*Smtp)(nil)

func NewSmtp(username, password string) *Smtp {
	return &Smtp{
		Username: username,
		Password: password,
	}
}

func (s *Smtp) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, "smtp.gmail.com")
	addr := "smtp.gmail.com:587"

	from := s.Username
	toSlice := []string{to}

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body + "\r\n")

	err := smtp.SendMail(addr, auth, from, toSlice, msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
