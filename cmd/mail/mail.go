package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-mail/mail/v2"
	"github.com/hisamcode/lenslocked/models"
	"github.com/joho/godotenv"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 2525
	username = "0ae8004190d135"
	password = "bdc050b3764e53"
)

func main() {
	// env
	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	return

	ess := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err1 := ess.ForgotPassword("hisamcode@gmail.com", "http://hisamcode.com/reset-pw?token=adasdasd2131")
	if err1 != nil {
		panic(err1)
	}

	return
	email := models.Email{
		From:      "hera@gmail.com",
		To:        "hisamcode@gmail.com",
		Subject:   "this is test email",
		Plaintext: "this is the body of the email",
		HTML:      "<h1>Hi dear</h1><p>This is email the email</p><p>Hope you enjoy it</p>",
	}
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	e := es.Send(email)
	if e != nil {
		panic(e)
	}

	fmt.Println("Email sent")

	return
	from := "hera@gmail.com"
	to := "hisamcode@gmail.com"
	subject := "this is a test email"
	plainText := "this is the body of the email"
	html := "<h1>Hi dear</h1><p>This is email the email</p><p>Hope you enjoy it</p>"

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/plain", plainText)
	msg.AddAlternative("text/html", html)
	// msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)
	// dial bagus kalo misalkan mengirim multiple email
	// sender, err := dialer.Dial()
	// if err != nil {
	// 	panic(err)
	// }
	// defer sender.Close()
	// sender.Send(from, []string{to}, msg)
	err3 := dialer.DialAndSend(msg)
	if err3 != nil {
		panic(err3)
	}

	fmt.Println("Message sent")
}
