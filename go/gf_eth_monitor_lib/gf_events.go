/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_eth_monitor_lib

import (
    "fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

//-------------------------------------------------
type GF_queue_info struct {
	name_str   string
	url_str    string
	aws_client *sqs.SQS
}

//-------------------------------------------------
// INIT_QUEUE
func init_queue(p_queue_name_str string) (*GF_queue_info, error) {


	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))


	svc := sqs.New(sess)


	// QUEUE_URL
	result_url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(p_queue_name_str),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			panic(fmt.Sprintf("Unable to find queue - %s", p_queue_name_str))
			return nil, err
		}
		return nil, err
	}

	fmt.Println(result_url)

	queue_info := &GF_queue_info{
		name_str:   p_queue_name_str,
		url_str:    *result_url.QueueUrl,
		aws_client: svc,
	}
	return queue_info, nil
}

//-------------------------------------------------
func process(p_queue_info *GF_queue_info) {

	// 10min - before this call returns if no message is present
	timeout_sec_int := 60*10

	result, err := p_queue_info.aws_client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:       aws.String(p_queue_info.url_str),
		AttributeNames: aws.StringSlice([]string{
			"SentTimestamp",
		}),
		MaxNumberOfMessages:   aws.Int64(1),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),

		// The duration (in seconds) for which the call waits for a message to arrive
    	// in the queue before returning. If a message is available, the call returns
    	// sooner than WaitTimeSeconds. If no messages are available and the wait time
    	// expires, the call returns successfully with an empty list of messages.
		WaitTimeSeconds: aws.Int64(int64(timeout_sec_int)),
	})
	if err != nil {
		panic(fmt.Sprintf("Unable to receive message from queue - %s - %v", p_queue_info.name_str, err))
	}
	
	fmt.Printf("Received %d messages.\n", len(result.Messages))
	if len(result.Messages) > 0 {
		fmt.Println(result.Messages)
	}
}