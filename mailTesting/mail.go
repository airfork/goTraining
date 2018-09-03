package main

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
)

/**
	Taken and modified from https://github.com/tangingw/go_smtp
**/

// GMAIL smtp server
const sMTPServer = "smtp.gmail.com"

// Sender has sender address and password
type Sender struct {
	User     string
	Password string
}

// NewSender ceates new sender
func NewSender(Username, Password string) Sender {

	return Sender{Username, Password}
}

// SendText sends the text via email
func (sender Sender) SendText(Dest []string, bodyMessage string) {

	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject:\n" + bodyMessage

	err := smtp.SendMail(sMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, sMTPServer),
		sender.User, Dest, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s", err)
		return
	}

	fmt.Println("Message sent successfully")
}

// WriteText writes the text to be sent and returns
func (sender Sender) WriteText(dest []string, bodyMessage string) string {

	header := make(map[string]string)
	header["From"] = sender.User

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	header["To"] = receipient
	header["Subject"] = ""
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("text/plain; charset=\"utf-8\"")
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}
