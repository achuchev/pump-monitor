package mail

import (
	"crypto/tls"
	"fmt"

	c "github.com/achuchev/pump-monitor/cmd/common"
	gomail "gopkg.in/mail.v2"
)

func Notify(subject string, body string) error {
	fmt.Println("Sending mail...")

	m := gomail.NewMessage()
	m.SetHeader("From", c.FlagsConfig.EmailAddressFrom)
	m.SetHeader("To", c.FlagsConfig.EmailAddressTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(c.FlagsConfig.EmailHost, c.FlagsConfig.EmailHostPort, c.FlagsConfig.EmailAddressFrom, c.FlagsConfig.EmailPassword)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: c.FlagsConfig.EmailHost}
	//d.StartTLSPolicy = gomail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Mail sent!")
	return nil
}
