package main

import (
	"bytes"
	"net/smtp"
)

// basicMailSend will connect to a remote mail server without authentication and send a message.
func basicMailSend(host string, to string, from string, subject string, body string) error {
	var (
		message string
		buf *bytes.Buffer
	)
	// connect to the remote server
	client, err := smtp.Dial(host)
	if err != nil {
		return err
	}
	defer client.Close()

	// set sender and and recipient
	client.Mail(from)
	client.Rcpt(to)

	// send the body
	mailContent, err := client.Data()
	if err != nil {
		return err
	}

	// close connection once function is complete
	defer mailContent.Close()

	message = "From: " + from + "\n"
	message += "To: " + to + "\n"
	message += "Subject: " + subject + "\n"
	message += "MIME-Version: 1.0\n"
	message += "Content-Type: text/html; charset=UTF-8\n"
	message += "<html>\n"
	message += "<body>\n"
	message += body
	message += "</body>\n"
	message += "</html>\n"
	message += "\n"

	// set up connection
	buf = bytes.NewBufferString(message)

	// send message
	if _, err = buf.WriteTo(mailContent); err != nil {
		return err
	}

	return nil
}
