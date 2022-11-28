/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	// "fmt"
    "context"
    "encoding/json"
    log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gloflow/gloflow/go/gf_core"
    // "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------------------

type GF_SQS_queue struct {
	Name_str    string
    AWS_url_str string
}

//-------------------------------------------------------------
// INIT

func SQSinit(pRuntimeSys *gf_core.RuntimeSys) (*sqs.Client, *gf_core.GFerror) {


    cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create AWS API session",
			"aws_client_v2_create",
			map[string]interface{}{},
			err, "gf_aws", pRuntimeSys)
        return nil, gfErr
	}

	client := sqs.NewFromConfig(cfg)
    return client, nil

    /*sess, err := session.NewSession()
    if err != nil {
        gfErr := gf_core.ErrorCreate("failed to create AWS API session",
			"aws_session_create",
			map[string]interface{}{},
			err, "gf_aws", pRuntimeSys)
        return nil, gfErr
    }
    svc := sqs.New(sess)
    return svc, nil*/
}

//-------------------------------------------------------------

func SQSgetQueueInfo(pSQSqueueNameStr string,
    pSQSclient *sqs.Client,
    pCtx        context.Context,
    pRuntimeSys *gf_core.RuntimeSys) (*GF_SQS_queue, *gf_core.GFerror) {


    sqs_queue_url_str, gfErr := SQSqueueGetURL(pSQSqueueNameStr,
        pSQSclient,
        pCtx,
        pRuntimeSys)
    if gfErr != nil {
        return nil, gfErr
    }
    queue_info := &GF_SQS_queue{
        Name_str:    pSQSqueueNameStr,
        AWS_url_str: sqs_queue_url_str,
    }
    return queue_info, nil
}

//-------------------------------------------------------------
// QUEUE_GET_URL

func SQSqueueGetURL(pSQSqueueNameStr string,
    pSQSclient  *sqs.Client,
    pCtx        context.Context,
    pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

    result, err := pSQSclient.GetQueueUrl(pCtx, &sqs.GetQueueUrlInput{
        QueueName: aws.String(pSQSqueueNameStr),
    })
    if err != nil {
        gfErr := gf_core.ErrorCreate("failed to get AWS SQS queue URL",
			"aws_sqs_queue_get_url_error",
			map[string]interface{}{
                "sqs_queue_name_str": pSQSqueueNameStr,
            },
			err, "gf_aws", pRuntimeSys)
        return "", gfErr
    }
    sqs_queue_url_str := *result.QueueUrl
    return sqs_queue_url_str, nil
}

//-------------------------------------------------------------
// QUEUE_CREATE

func SQSqueueCreate(pSQSqueueNameStr string,
    pSQSclient  *sqs.Client,
    pCtx         context.Context,
    pRuntimeSys *gf_core.RuntimeSys) (*GF_SQS_queue, *gf_core.GFerror) {

    

    log.WithFields(log.Fields{"name": pSQSqueueNameStr,}).Info("AWS_SQS - creating new queue")

    // CREATE
    _, err := pSQSclient.CreateQueue(pCtx, &sqs.CreateQueueInput{
        QueueName: aws.String(pSQSqueueNameStr),
    })
    if err != nil {
        gfErr := gf_core.ErrorCreate("failed to create AWS SQS queue",
			"aws_sqs_queue_create_error",
			map[string]interface{}{
                "sqs_queue_name_str": pSQSqueueNameStr,
            },
			err, "gf_aws", pRuntimeSys)
        return nil, gfErr
    }


    // URL
    queue_url_str, gfErr := SQSqueueGetURL(pSQSqueueNameStr, pSQSclient, pCtx, pRuntimeSys)
    if gfErr != nil {
        return nil, gfErr
    }


    queue := &GF_SQS_queue{
        Name_str:    pSQSqueueNameStr,
        AWS_url_str: queue_url_str,
    }

	return queue, nil
}

//-------------------------------------------------------------
// QUEUE_DELETE

func SQSqueueDelete(pSQSqueueNameStr string,
    pAWSclient  *sqs.Client,
    pCtx        context.Context,
    pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {



	return nil
}

//-------------------------------------------------------------

func SQSmsgPull(pQueueInfo *GF_SQS_queue,
    pAWSclient  *sqs.Client,
    pCtx        context.Context,
    pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {

    log.WithFields(log.Fields{"name": pQueueInfo.Name_str,}).Info("AWS_SQS - pull message from queue")

    // Must be >= 0 and <= 20, if provided
    timeout_sec_int := 20

	// SQS_RECEIVE_MESSAGE
	result, err := pAWSclient.ReceiveMessage(pCtx, &sqs.ReceiveMessageInput{
		QueueUrl:       aws.String(pQueueInfo.AWS_url_str),
		AttributeNames: []types.QueueAttributeName{
			"SentTimestamp",
		},
		MaxNumberOfMessages:   1,
		MessageAttributeNames: []string { // []types.QueueAttributeName{
			"All",
		},

		// The duration (in seconds) for which the call waits for a message to arrive
    	// in the queue before returning. If a message is available, the call returns
    	// sooner than WaitTimeSeconds. If no messages are available and the wait time
    	// expires, the call returns successfully with an empty list of messages.
		WaitTimeSeconds: int32(timeout_sec_int),
	})
	if err != nil {
        gfErr := gf_core.ErrorCreate("failed to receive a message from SQS queue durring msg_pull function",
			"aws_sqs_queue_receive_msg_error",
			map[string]interface{}{
                "sqs_queue_name_str": pQueueInfo.Name_str,
                "sqs_queue_url_str":  pQueueInfo.AWS_url_str,
            },
			err, "gf_aws", pRuntimeSys)
		return nil, gfErr
	}
    
    // fmt.Println("*******************************************")
    // spew.Dump(result.Messages)

    if result.Messages != nil {

        msg          := result.Messages[0]
        msg_body_str := *msg.Body

        log.WithFields(log.Fields{"name": pQueueInfo.Name_str, "msg_id": *msg.MessageId, "msg_body": msg_body_str}).Info("AWS_SQS - pull message from queue - OK")

        //--------------------------
        // DECODE_MESSAGE
        msg_map      := map[string]interface{}{}
        if err := json.Unmarshal([]byte(msg_body_str), &msg_map); err != nil {
            gfErr := gf_core.ErrorCreate("failed to JSON decode a message body pulled from SQS",
                "json_decode_error",
                map[string]interface{}{
                    "sqs_queue_name_str": pQueueInfo.Name_str,
                    "sqs_queue_url_str":  pQueueInfo.AWS_url_str,
                },
                err, "gf_aws", pRuntimeSys)
            return nil, gfErr
        }

        //--------------------------
        // DELETE_MESSAGE

        log.WithFields(log.Fields{"name": pQueueInfo.Name_str, "msg_id": *msg.MessageId}).Info("AWS_SQS - delete message from queue")

        msg_receipt_handle := msg.ReceiptHandle

        // https://docs.aws.amazon.com/sdk-for-go/api/service/sqs/#SQS.DeleteMessage
        _, err = pAWSclient.DeleteMessage(pCtx, &sqs.DeleteMessageInput{
            QueueUrl:      aws.String(pQueueInfo.AWS_url_str),
            ReceiptHandle: msg_receipt_handle,
        })
        if err != nil {
            gfErr := gf_core.ErrorCreate("failed to delete a message from SQS queue durring msg_pull function, after receiving it",
                "aws_sqs_queue_delete_msg_error",
                map[string]interface{}{
                    "sqs_queue_name_str": pQueueInfo.Name_str,
                    "sqs_queue_url_str":  pQueueInfo.AWS_url_str,
                },
                err, "gf_aws", pRuntimeSys)
            return nil, gfErr
        }

        //--------------------------

        return msg_map, nil
    }

	return nil, nil
}

//-------------------------------------------------------------
// MSG_PUSH

func SQSmsgPush(p_msg interface{},
    pQueueInfo  *GF_SQS_queue,
    pAWSclient  *sqs.Client,
    pCtx         context.Context,
    pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

    
    log.WithFields(log.Fields{"name": pQueueInfo.Name_str,}).Info("AWS_SQS - push message to queue")

    msg_data_JSON_encoded, err := json.Marshal(p_msg)
    if err != nil {
        gfErr := gf_core.ErrorCreate("failed to JSON encode a message to send to SQS queue",
			"json_encode_error",
			map[string]interface{}{
                "sqs_queue_name_str": pQueueInfo.Name_str,
                "sqs_queue_url_str":  pQueueInfo.AWS_url_str,
            },
			err, "gf_aws", pRuntimeSys)
        return gfErr
    }

	_, err = pAWSclient.SendMessage(pCtx, &sqs.SendMessageInput{
        MessageBody: aws.String(string(msg_data_JSON_encoded)),
        QueueUrl:    aws.String(pQueueInfo.AWS_url_str),

        // DelaySeconds: 10,
        /*MessageAttributes: map[string]types.MessageAttributeValue{

			"time_sec": {
                DataType:    aws.String("String"),
                StringValue: aws.String(fmt.Sprintf("%f", p_msg.TimeSec)),
            },
			"module": {
                DataType:    aws.String("String"),
                StringValue: aws.String(p_msg.Module),
            },
            "type": {
                DataType:    aws.String("String"),
                StringValue: aws.String(p_msg.Type),
            },
            "msg": {
                DataType:    aws.String("String"),
                StringValue: aws.String(p_msg.Msg),
            },
        },*/
    })

    if err != nil {
        gfErr := gf_core.ErrorCreate("failed to send a message to SQS queue",
			"aws_sqs_queue_send_msg_error",
			map[string]interface{}{
                "sqs_queue_name_str": pQueueInfo.Name_str,
                "sqs_queue_url_str":  pQueueInfo.AWS_url_str,
            },
			err, "gf_aws", pRuntimeSys)
        return gfErr
    }

    // msg_id_str := *result.MessageId
    // fmt.Printf("Success - msg ID - %s\n", msg_id_str)

	return nil
}