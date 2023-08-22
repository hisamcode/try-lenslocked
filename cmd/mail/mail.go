package main

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 2525
	username = "0ae8004190d135"
	password = "bdc050b3764e53"
)

func main() {
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
	err := dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent")
}
