package models


import (
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/awserr"
"github.com/aws/aws-sdk-go/service/ses"
"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewEmail(db *sqlx.DB, awsSession *session.Session) *Email {
	Email := &Email{}
	Email.db = db
	Email.awsSession = awsSession
	return Email
}

type EmailRow struct {
	Subject string `db:"subject" json:"subject"`
	EmailTo string `db:"emailto" json:"emailto"`
	EmailFrom string `db:"emailfrom" json:"emailfrom"`
	BodyHTML string `db:"bodyhtml" json:"bodyhtml"`
	BodyText string `db:"bodytext" json:"bodytext"`

	//EmailCC []*string `db:"emailcc" json:"emailcc,omitempty"`

}

type Email struct {
	Base
	EmailRow
}

func (e *Email) SendEmail() error {

	log.Debugf("Session :v", e.awsSession)

	svc := ses.New(e.awsSession)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{

			ToAddresses: []*string{
				aws.String(e.EmailTo),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{

				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(e.BodyHTML),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(e.BodyText),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(e.Subject),
			},
		},

		Source:        aws.String(e.EmailFrom),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Debugf(ses.ErrCodeMessageRejected, aerr.Error())

			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Debugf(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Debugf(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Debugf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Debugf(err.Error())

			return err
		}
		return err
	}

	fmt.Println(result)

	return nil
}


