package scripts

import (
	"gopkg.in/gomail.v2"
	"log"
)

func SendMail(to string, subject string, body string) error {
	// Create a new email message
	m := gomail.NewMessage()

	// Set the sender
	m.SetHeader("From", "bbeka544@mail.ru")

	// Set the recipient
	m.SetHeader("To", to)

	// Set the subject
	m.SetHeader("Subject", subject)

	// Set the email body (plain text)
	m.SetBody("text/plain", body)

	// Create a new dialer with Mail.ru's SMTP settings
	d := gomail.NewDialer("smtp.mail.ru", 465, "bbeka544@mail.ru", "TMeJXXr7y6dt3njgTkE1")

	// Enable SSL (required for Mail.ru)
	d.SSL = true

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
		return err
	}
	return nil
}
