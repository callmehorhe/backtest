package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"

	"github.com/spf13/viper"
)

type EmailService struct {
	Client  *smtp.Client
	Headers map[string]string
}

func NewEmailService() *EmailService {
	host := viper.GetString("email.host")
	port := viper.GetString("email.port")
	servername := host + ":" + port
	username := viper.GetString("email.username")
	password := viper.GetString("email.password")
	auth := smtp.PlainAuth("", username, password, host)
	from := mail.Address{Name: "Delivery", Address: "reg_deliveryhouse@mail.ru"}
	headers := make(map[string]string)
	headers["From"] = from.String()

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	if err = client.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	return &EmailService{
		Client:  client,
		Headers: headers,
	}
}

func (e *EmailService) SendEmail(email, subject, text string) error {
	if err := e.Client.Rcpt(email); err != nil {
		log.Panic(err)
	}
	e.Headers["To"] = email
	e.Headers["Subject"] = subject

	message := ""

	for k, v := range e.Headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + text

	send, err := e.Client.Data()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer send.Close()

	_, err = send.Write([]byte(message))
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
