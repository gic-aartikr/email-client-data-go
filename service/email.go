package service

import (
	"context"
	"emaildata/model"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Email struct {
	Server     string
	Database   string
	Collection string
}

var Collection *mongo.Collection
var ctx = context.TODO()
var sendEmail = "Message Sent successfully!"

func (e *Email) Connect() {
	clientOptions := options.Client().ApplyURI(e.Server)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database(e.Database).Collection(e.Collection)
}

func (e *Email) Insert(emailBody model.Email) (string, error) {
	emailSend(emailBody)
	_, err := Collection.InsertOne(ctx, emailBody)

	if err != nil {
		return sendEmail, errors.New("Unable To Insert New Record")

	}

	return sendEmail, nil
}

func emailSend(emailData model.Email) {

	username := "ranveer.singh@gridinfocom.com"
	passwd := "xsmtpsib-2af236a5040e4b54343f4fe5b59826e9e7588b2e33d160249ab60a3060fbf348-LAnNRJT0sxHcPqYr"
	smtpHost := "smtp-relay.sendinblue.com"
	smtpPort := "587"

	to := emailData.EmailTo
	msg := ComposeMsg(emailData)
	// auth := smtp.PlainAuth("", username, passwd, smtpHost)
	auth := smtp.PlainAuth("", username, passwd, smtpHost)
	// err := smtp.SendMail(host+":2500", auth, sender, to, []byte(msg))
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, username, to, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Message Sent successfully!")

}

// construct message based on the cc
func ComposeMsg(emailData model.Email) string {
	// empty string
	msg := ""
	// set sender
	msg += fmt.Sprintf("From: %s\r\n", emailData.EmailTo)
	// if more than 1 recipient
	if len(emailData.EmailTo) > 0 {
		msg += fmt.Sprintf("Cc: %s\r\n", strings.Join(emailData.EmailCC, ";"))
	}
	// add subject
	msg += fmt.Sprintf("Subject: %s\r\n", emailData.Subject)
	// add mail body
	msg += fmt.Sprintf("\r\n%s\r\n", emailData.EmailBody)
	return msg
}
