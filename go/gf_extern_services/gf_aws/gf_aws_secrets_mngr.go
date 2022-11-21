/*
MIT License

Copyright (c) 2021 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_aws

import (
	"fmt"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------

func AWSsecretsMngrCreateSecret(pSecretNameStr string,
	pDataMap        map[string]interface{},
	pDescriptionStr string,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	dataStr, err := json.Marshal(pDataMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to JSON encode AWS secret",
			"json_encode_error",
			map[string]interface{}{"secret_name_str": pSecretNameStr,},
			err, "gf_aws", pRuntimeSys)
		return gfErr
	}
	
	svc   := secretsmanager.New(session.New())
	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String(pDescriptionStr),
		Name:         aws.String(pSecretNameStr),
		SecretString: aws.String(string(dataStr)),
	}
	
	_, err = svc.CreateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		}

		gfErr := gf_core.ErrorCreate("failed to create AWS secret",
			"aws_secretsmngr_create_secret_value_error",
			map[string]interface{}{"secrets_name_str": pSecretNameStr,},
			err, "gf_aws", pRuntimeSys)
		return gfErr
	}

	return nil
}

//-------------------------------------------------------------

func AWSsecretsMngrGetSecret(pSecretNameStr string,
	pRuntimeSys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {

	svc   := secretsmanager.New(session.New())
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(pSecretNameStr),
		// VersionStage: aws.String("AWSPREVIOUS"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		}
		
		gfErr := gf_core.ErrorCreate("failed to get AWS secret",
			"aws_secretsmngr_get_secret_value_error",
			map[string]interface{}{"secrets_name_str": pSecretNameStr,},
			err, "gf_aws", pRuntimeSys)
		return nil, gfErr
	}

	valueStr := *result.SecretString

	//--------------
	var sMap map[string]interface{}
	err = json.Unmarshal([]byte(valueStr), &sMap)

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to JSON parse AWS secret",
			"json_decode_error",
			map[string]interface{}{"secret_name_str": pSecretNameStr,},
			err, "gf_aws", pRuntimeSys)
		return nil, gfErr
	}

	//--------------

	return sMap, nil
}