package custommodule

import (
	"gopkg.in/gomail.v2"
)

type Emailer struct {
	Host     string
	Port     int
	Email    string
	Password string
}

func (emailer *Emailer) SendMail(to string, subject, message string) error {

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailer.Email)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", message)

	dialer := gomail.NewPlainDialer(emailer.Host,
		emailer.Port,
		emailer.Email,
		emailer.Password,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}
	//

	return nil
}
