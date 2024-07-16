package mail

import "gopkg.in/gomail.v2"

type Mailer interface {
	Send(to string, content string) error
}

type Mail struct {
	mail gomail.Sender
}