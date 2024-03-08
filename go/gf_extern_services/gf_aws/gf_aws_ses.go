/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_aws

import (
    "fmt"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
    "github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func SESsendMessage(pRecipiendAddressStr string,
	pSenderAddressStr string,
    pSubjectStr       string,
	pHTMLbodyStr      string,
	pTextBodyStr      string,
    pRuntimeSys       *gf_core.RuntimeSys) *gf_core.GFerror {

    svc := ses.New(session.New())

	input := &ses.SendEmailInput{

		// sender
		Source: aws.String(pSenderAddressStr),

		// recipient
        Destination: &ses.Destination{
            ToAddresses: []*string{
                aws.String(pRecipiendAddressStr),
            },
        },

		// message
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Charset: aws.String("UTF-8"),
                    Data:    aws.String(pHTMLbodyStr),
                },
                Text: &ses.Content{
                    Charset: aws.String("UTF-8"),
                    Data:    aws.String(pTextBodyStr),
                },
            },
            Subject: &ses.Content{
                Charset: aws.String("UTF-8"),
                Data:    aws.String(pSubjectStr),
            },
        },
    }

	_, err := svc.SendEmail(input)

	if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {

			// message_rejected
            case ses.ErrCodeMessageRejected:
                fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())

			// mail_from_domain_not_verified
            case ses.ErrCodeMailFromDomainNotVerifiedException:
                fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())

			// config_does_not_exist
            case ses.ErrCodeConfigurationSetDoesNotExistException:
                fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        }

		gfErr := gf_core.ErrorCreate("failed to send AWS SES email message",
			"aws_ses_service_send_message_error",
			map[string]interface{}{
				"recipiend_address_str": pRecipiendAddressStr,
				"sender_address_str":    pSenderAddressStr,
				"subject_str":           pSubjectStr,
			},
			err, "gf_aws", pRuntimeSys)
        return gfErr
    }
    return nil
}

//---------------------------------------------------
// verifies an email address with SES so that it can be 
// used for sending emails from that address.
// not used frequently

func SESverifyAddress(pRecipiendAddressStr string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	svc := ses.New(session.New())

	_, err := svc.VerifyEmailAddress(&ses.VerifyEmailAddressInput{
		EmailAddress: aws.String(pRecipiendAddressStr),
	})

	if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case ses.ErrCodeMessageRejected:
                fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())

            case ses.ErrCodeMailFromDomainNotVerifiedException:
                fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())

            default:
                fmt.Println(aerr.Error())
            }
        }

		gfErr := gf_core.ErrorCreate("failed to verify AWS SES address",
			"aws_ses_service_verify_address_error",
			map[string]interface{}{"recipiend_address_str": pRecipiendAddressStr,},
			err, "gf_aws", pRuntimeSys)
        return gfErr
    }

	return nil
}