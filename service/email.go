package service

import (
	"context"
	"emaildata/model"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	emailBody.Date = time.Now()
	_, err := Collection.InsertOne(ctx, emailBody)

	if err != nil {
		return sendEmail, errors.New("Unable To Insert New Record")

	}

	return sendEmail, nil
}

func emailSend(emailData model.Email) {

	username := "vidhi.goel@gridinfocom.com"
	passwd := "pGL756txPrWkSBX4"
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

/////////////////////search data/////////////////

func (e *Email) SearchData(search model.EmailSearch) ([]*model.Email, error) {
	var searchData []*model.Email

	filter := bson.D{}

	if search.EmailTo != "" {
		filter = append(filter, primitive.E{Key: "email_to", Value: bson.M{"$regex": search.EmailTo}})
	}
	if search.EmailCC != "" {
		filter = append(filter, primitive.E{Key: "email_cc", Value: bson.M{"$regex": search.EmailCC}})
	}
	if search.EmailBCC != "" {
		filter = append(filter, primitive.E{Key: "email_bcc", Value: bson.M{"$regex": search.EmailBCC}})
	}
	if search.Subject != "" {
		filter = append(filter, primitive.E{Key: "subject", Value: bson.M{"$regex": search.Subject}})
	}

	t, _ := time.Parse("2006-01-02", search.Date)
	if search.Date != "" {
		filter = append(filter, primitive.E{Key: "date", Value: bson.M{
			"$gte": primitive.NewDateTimeFromTime(t)}})
	}

	result, err := Collection.Find(ctx, filter)

	if err != nil {
		return searchData, err
	}

	for result.Next(ctx) {
		var data model.Email
		err := result.Decode(&data)
		if err != nil {
			return searchData, err
		}
		searchData = append(searchData, &data)
	}

	if searchData == nil {
		return searchData, errors.New("Data Not Found In DB")
	}

	return searchData, nil
}
