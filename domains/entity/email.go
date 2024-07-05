package entity

// EmailPayload represents the data required to send an email.
type EmailPayload struct {
	To      string // To Recipient email address
	Subject string // Subject of the email
	Body    string // Body of the email
}
