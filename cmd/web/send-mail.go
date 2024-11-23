package main

import (
	"fmt"
	"github.com/Joshua-Russel/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"strings"
	"time"
)

func listenForMail() {

	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()

}
func sendMsg(msg models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = false

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(msg.From).SetSubject(msg.Subject).AddTo(msg.To)
	if msg.Template == "" {

		email.SetBody(mail.TextHTML, msg.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", msg.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailtemplate := string(data)
		msgToSend := strings.Replace(mailtemplate, "[%body%]", msg.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	} else {
		log.Println("email sent")
	}
}
