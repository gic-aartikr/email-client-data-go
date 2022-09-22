package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Email struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmailTo   []string           `bson:"email_to,omitempty" json:"emailTo,omitempty"`
	EmailCC   []string           `bson:"email_cc,omitempty" json:"emailCc,omitempty"`
	EmailBCC  []string           `bson:"email_bcc, omitempty" json:"emailBcc,omitempty"`
	Subject   string             `bson:"subject,omitempty" json:"subject,omitempty"`
	EmailBody string             `bson:"email_body,omitempty" json:"email_body,omitempty"`
	Date      time.Time          `bson:"date,omitempty" json:"date,omitempty"`
}

type EmailSearch struct {
	EmailTo  string `bson:"email_to,omitempty" json:"emailTo,omitempty"`
	EmailCC  string `bson:"email_cc,omitempty" json:"emailCc,omitempty"`
	EmailBCC string `bson:"email_bcc, omitempty" json:"emailBcc,omitempty"`
	Subject  string `bson:"subject,omitempty" json:"subject,omitempty"`
	Date     string `bson:"date,omitempty" json:"date,omitempty"`
}
