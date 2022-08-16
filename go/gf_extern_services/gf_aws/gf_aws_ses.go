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
func AWS_SES__send_message(p_recipiend_address_str string,
	p_sender_address_str string,
    p_subject_str        string,
	p_html_body_str      string,
	p_text_body_str      string,
    p_runtime_sys        *gf_core.RuntimeSys) *gf_core.GF_error {

    svc := ses.New(session.New())

	input := &ses.SendEmailInput{

		// sender
		Source: aws.String(p_sender_address_str),

		// recipient
        Destination: &ses.Destination{
            ToAddresses: []*string{
                aws.String(p_recipiend_address_str),
            },
        },

		// message
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Charset: aws.String("UTF-8"),
                    Data:    aws.String(p_html_body_str),
                },
                Text: &ses.Content{
                    Charset: aws.String("UTF-8"),
                    Data:    aws.String(p_text_body_str),
                },
            },
            Subject: &ses.Content{
                Charset: aws.String("UTF-8"),
                Data:    aws.String(p_subject_str),
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

		gf_err := gf_core.ErrorCreate("failed to send AWS SES email message",
			"aws_ses_service_send_message_error",
			map[string]interface{}{
				"recipiend_address_str": p_recipiend_address_str,
				"sender_address_str":    p_sender_address_str,
				"subject_str":           p_subject_str,
			},
			err, "gf_aws", p_runtime_sys)
        return gf_err
    }
    return nil
}

//---------------------------------------------------
// verifies an email address with SES so that it can be 
// used for sending emails from that address.
// not used frequently
func AWS_SES__verify_address(p_recipiend_address_str string,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GF_error {

	svc := ses.New(session.New())

	_, err := svc.VerifyEmailAddress(&ses.VerifyEmailAddressInput{
		EmailAddress: aws.String(p_recipiend_address_str),
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

		gf_err := gf_core.ErrorCreate("failed to verify AWS SES address",
			"aws_ses_service_verify_address_error",
			map[string]interface{}{"recipiend_address_str": p_recipiend_address_str,},
			err, "gf_aws", p_runtime_sys)
        return gf_err
    }

	return nil
}