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
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_queue_info struct {
	name_str   string
	url_str    string
	aws_client *sqs.SQS
}

//-------------------------------------------------
// INIT_QUEUE
func Event__init_queue(p_queue_name_str string) (*GF_queue_info, error) {


	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))


	svc := sqs.New(sess)


	// QUEUE_URL
	fmt.Printf("get AWS SQS queue - %s\n", p_queue_name_str)
	result_url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(p_queue_name_str),
	})
	if err != nil {
		fmt.Println(fmt.Sprint(err))

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
func event__start_sqs_consumer(p_queue_info *GF_queue_info,
	p_metrics *GF_metrics) {

	go func() {

		for {
			Event__process_from_sqs(p_queue_info, p_metrics)
		}
	}()
}

//-------------------------------------------------
func Event__process_from_sqs(p_queue_info *GF_queue_info,
	p_metrics *GF_metrics) {

	// 20s - before this call returns if no message is present.
	// Must be >= 0 and <= 20
	timeout_sec_int := 20


	// SQS_RECEIVE_MESSAGE
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



	for _, m := range result.Messages {
		
		SQS_timestamp_str := *m.Attributes["SentTimestamp"]
		fmt.Printf("SQS_timestamp - %s\n", SQS_timestamp_str)

		// JSON_DECODE
		var event_map map[string]interface{}
		json.Unmarshal([]byte(*m.Body), &event_map)

		//---------------------------
		// EVENT__PROCESS
		event__process(event_map, p_metrics)

		//---------------------------

		// DELETE_MESSAGE
		// https://docs.aws.amazon.com/sdk-for-go/api/service/sqs/#SQS.DeleteMessage
		_, err := p_queue_info.aws_client.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(p_queue_info.url_str),
			ReceiptHandle: m.ReceiptHandle,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to delete message from queue - %s - %v", p_queue_info.name_str, err))

		}
	}
}

//-------------------------------------------------
func event__process(p_event_map map[string]interface{},
	p_metrics *GF_metrics,) {



	fmt.Println(" PROCESS EVENT ================")
	spew.Dump(p_event_map)



	if p_metrics != nil {
		p_metrics.counter__sqs_msgs_num.Inc()
	}

	
}