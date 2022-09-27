package service

import (
	"context"
	"crypto/tls"
	"emaildata/modelData"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gomail "gopkg.in/mail.v2"
)

const fileFromPath = "D:/go-lang/email-client-data-go/demo/sample.pdf"

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

func (e *Email) Insert(emailData modelData.EmailModel) (string, error) {

	err := sendMailWithAttachment(emailData)
	if err != nil {
		return "", err
	}

	insertData := e.insertMenthod(emailData)
	if insertData != nil {
		return "", insertData
	}

	return "Email Sent Successfully", nil
}

func (*Email) insertMenthod(emailData modelData.EmailModel) error {

	var saveEmail modelData.Email
	saveEmail.EmailTo = emailData.EmailTo
	saveEmail.EmailCC = emailData.EmailCC
	saveEmail.EmailBCC = emailData.EmailBCC
	saveEmail.EmailBody = emailData.EmailBody
	saveEmail.Subject = emailData.EmailSubject
	saveEmail.Date = time.Now()
	_, err := Collection.InsertOne(ctx, saveEmail)

	if err != nil {
		return errors.New("Unable To Insert New Record")
	}
	return err
}

func sendMailWithAttachment(emailData modelData.EmailModel) error {

	m := gomail.NewMessage()

	m.SetHeaders(map[string][]string{
		"From": {m.FormatAddress("aarti.kumari@gridinfocom.com", "Aarti")},
		"To":   emailData.EmailTo,
		"Cc":   emailData.EmailCC,
		// "Subject": emailData.EmailSubject,
	})
	m.SetHeader("Subject", emailData.EmailSubject)

	for i := range emailData.FileLocation {
		m.Attach(emailData.FileLocation[i])
	}
	m.SetBody("text/plain", emailData.EmailBody)
	d := gomail.NewDialer("smtp-relay.sendinblue.com", 587, "aarti.kumari@gridinfocom.com", "f7c3OQF8UzC6pIh1")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/////////////////////search data/////////////////

func (e *Email) SearchData(search modelData.EmailSearch) ([]*modelData.Email, error) {
	var searchData []*modelData.Email

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
		var data modelData.Email
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

func (e *Email) WriteEmailDataInPDF(id string) error {
	var emailData []*modelData.Email
	dir := "data/download/"
	file := "searchData" + fmt.Sprintf("%v", time.Now().Format("3_4_5_pm"))

	emailId, err := primitive.ObjectIDFromHex(id)
	fmt.Println("emailId:", emailId)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: emailId}}

	resultData, err := Collection.Find(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("data record:", resultData)
	for resultData.Next(ctx) {
		var data modelData.Email
		err := resultData.Decode(&data)
		if err != nil {
			return err
		}
		emailData = append(emailData, &data)

	}

	if emailData == nil {
		return errors.New("Data Not Found In DB")
	}

	_, err = writeToPdf(dir, file, emailData)

	if err != nil {
		return err
	}

	fmt.Println("emailData:", emailData)
	return nil
}
func writeToPdf(dir, file string, emailData []*modelData.Email) (*creator.Creator, error) {
	c := creator.New()

	// datare:= data.
	err := license.SetMeteredKey("72c4ab06d023bbc8b2e186d089f9e052654afea32b75141f39c7dc1ab3b108ca")

	robotoFontRegular, err := model.NewPdfFontFromTTFFile("Roboto/Roboto-Regular.ttf")
	if err != nil {
		return c, err
	}

	robotoFontPro, err := model.NewPdfFontFromTTFFile("Roboto/Roboto-Bold.ttf")
	if err != nil {
		return c, err
	}

	c.SetPageMargins(50, 50, 50, 50)

	// c.NewPage()

	ch := c.NewChapter("Email Data")
	chapterFont := robotoFontPro
	chapterFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	// chapterRedColor := creator.ColorRGBFrom8bit(255, 0, 0)
	chapterFontSize := 18.0

	normalFont := robotoFontRegular
	// normalFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	normalFontColorGreen := creator.ColorRGBFrom8bit(4, 79, 3)
	normalFontSize := 10.0

	ch.GetHeading().SetFont(chapterFont)
	ch.GetHeading().SetFontSize(chapterFontSize)
	ch.GetHeading().SetColor(chapterFontColor)

	for i := range emailData {

		p := c.NewParagraph(emailData[i].Subject)
		p.SetFont(normalFont)
		p.SetFontSize(normalFontSize)
		p.SetColor(normalFontColorGreen)
		p.SetMargins(0, 0, 5, 0)
		p.SetLineHeight(2)
		ch.Add(p)

		p = c.NewParagraph("To" + ":" + "  " + convertArrayOfStringIntoString(emailData[i].EmailTo))
		p.SetFont(normalFont)
		p.SetFontSize(normalFontSize)
		p.SetColor(normalFontColorGreen)
		p.SetMargins(0, 0, 5, 0)
		p.SetLineHeight(2)
		ch.Add(p)

		p = c.NewParagraph("CC" + ":" + "  " + convertArrayOfStringIntoString(emailData[i].EmailCC))
		p.SetFont(normalFont)
		p.SetFontSize(normalFontSize)
		p.SetColor(normalFontColorGreen)
		p.SetMargins(0, 0, 5, 0)
		p.SetLineHeight(2)
		ch.Add(p)

		p = c.NewParagraph("bcc" + ":" + "  " + convertArrayOfStringIntoString(emailData[i].EmailBCC))
		p.SetFont(normalFont)
		p.SetFontSize(normalFontSize)
		p.SetColor(normalFontColorGreen)
		p.SetMargins(0, 0, 5, 0)
		p.SetLineHeight(2)
		ch.Add(p)

		p = c.NewParagraph(emailData[i].EmailBody)
		p.SetFont(normalFont)
		p.SetFontSize(normalFontSize)
		p.SetColor(normalFontColorGreen)
		p.SetMargins(0, 0, 5, 0)
		p.SetLineHeight(2)
		ch.Add(p)

	}
	// var buf bytes.Buffer
	// c.Write(&buf)
	c.Draw(ch)
	c.WriteToFile(dir + file + "report.pdf")
	return c, nil
}

func convertArrayOfStringIntoString(str []string) string {

	finalData := ""

	y := 0

	for x := range str {

		if y != 0 {

			finalData = finalData + ", "

		}

		finalData = finalData + str[x]

		y++

	}

	y = 0

	return finalData

}
