﻿package email

import "github.com/wisle25/be-template/domains/entity"

// Email defines the interface for an email service.
// Any struct that implements this interface can be used to send emails.
type Email interface {
	// SendEmail sends an email based on the provided payload.
	// The payload should contain necessary email details like recipient, subject, and body.
	SendEmail(payload *entity.EmailPayload)
}
