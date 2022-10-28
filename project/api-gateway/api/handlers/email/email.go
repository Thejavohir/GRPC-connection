package v1

import (
	"net/smtp"
	"os"
)

func SendMail(to []string, message []byte) error {
	from := os.Getenv("javohir@gmail.com")
	password := os.Getenv("safwefadf")
	host := "smtp.gmail.com"
	port := "587"
	auth := smtp.PlainAuth("", from, password, host)
	errSMTP := smtp.SendMail(host+":"+port, auth, from, to, message)
	return errSMTP
}
