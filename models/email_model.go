package models


import (
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/awserr"
"github.com/aws/aws-sdk-go/service/ses"
"fmt"
	"github.com/jmoiron/sqlx"
)

func NewEmail(db *sqlx.DB, aws *ses) *Email {
	Email := &Email{}
	Email.db = db
	Email.aws = aws
	return Email
}

type EmailRow struct {
	Subject string `db:"subject" json:"subject"`
	EmailTo string `db:"emailto" json:"emailto"`
	EmailFrom string `db:"emailfrom" json:"emailfrom"`
	BodyHTML string `db:"bodyhtml" json:"bodyhtml"`
	BodyText string `db:"bodytext" json:"bodytext"`

	EmailCC[] string `db:"emailcc" json:"emailcc,omitempty"`

}

type Email struct {
	Base
	EmailRow
}

func (e *Email) SendEmail(){
	svc := ses.New(e.aws)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{
				aws.String(e.EmailCC),
			},
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
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}


