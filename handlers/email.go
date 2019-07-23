package handlers

import (
	"asira/asira"
	"fmt"
	"net/smtp"
	"strings"
)

func sendMail(to []string, subject, message string) error {
	Config := asira.App.Config.GetStringMap(fmt.Sprintf("%s.mailer", asira.App.ENV))
	body := "From: " + Config["SMTP_HOST"].(string) + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", Config["EMAIL"].(string), Config["PASSWORD"].(string), Config["SMTP_HOST"].(string))
	smtpAddr := fmt.Sprintf("%s:%d", Config["SMTP_HOST"].(string), Config["SMTP_PORT"].(int))

	err := smtp.SendMail(smtpAddr, auth, Config["EMAIL"].(string), append(to), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
