package emails

import (
	"fmt"
	"github.com/wisle25/be-template/applications/email"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/entity"
	"net/smtp"
)

type SmtpEmailService struct /* implements Email */ {
	smtpHost  string
	smtpPort  string
	auth      smtp.Auth
	emailFrom string
}

func NewStmpEmailService(config *commons.Config) email.Email {
	auth := smtp.PlainAuth("", config.SmtpUsername, config.SmtpPassword, config.SmtpHost)

	return &SmtpEmailService{
		smtpHost:  config.SmtpHost,
		smtpPort:  config.SmtpPort,
		auth:      auth,
		emailFrom: config.SmtpUsername,
	}
}

func (s *SmtpEmailService) SendEmail(email *entity.EmailPayload) {
	// Arrange email
	msg := []byte("From: " + s.emailFrom + "\n" +
		"To: " + email.To + "\n" +
		"Subject: " + email.Subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<html><body>" + email.Body + "</body></html>",
	)

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	// Send
	err := smtp.SendMail(addr, s.auth, s.emailFrom, []string{email.To}, msg)
	if err != nil {
		commons.ThrowServerError("send_email_err: sending email", err)
	}
}
