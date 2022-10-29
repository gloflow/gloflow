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

package gf_aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------
func AWS_ECR__update_service(p_service_name_str string,
	p_cluster_name_str         string,
	p_healthy_percent__min_int int,
	p_runtime_sys              *gf_core.RuntimeSys) *gf_core.GFerror {

	fmt.Printf("AWS ECS UPDATE_SERVICE - %s\n", p_service_name_str)

	svc := ecs.New(session.New())
	
	input := &ecs.UpdateServiceInput{
		Service: aws.String(p_service_name_str),
		Cluster: aws.String(p_cluster_name_str),
		
		DeploymentConfiguration: &ecs.DeploymentConfiguration{
			MinimumHealthyPercent: aws.Int64(int64(p_healthy_percent__min_int)),
		},

		// IMPORTANT!! - for "dev" cluster only, when the same tagged container image
		//               is deployed again.
		// Whether to force a new deployment of the service. Deployments are not forced
    	// by default. You can use this option to trigger a new deployment with no service
    	// definition changes.
		ForceNewDeployment: aws.Bool(true),

		// TaskDefinition: aws.String(p_task_def_str),
	}

	// UPDATE_SERVICE
	result, err := svc.UpdateService(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())

			// CLUSTER_NOT_FOUND
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())

			// SERVICE_NOT_FOUND
			case ecs.ErrCodeServiceNotFoundException:
				fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotActiveException:
				fmt.Println(ecs.ErrCodeServiceNotActiveException, aerr.Error())
			case ecs.ErrCodePlatformUnknownException:
				fmt.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())

			case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
				fmt.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())

			// ACCESS_DENIED
			case ecs.ErrCodeAccessDeniedException:
				fmt.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

		gf_err := gf_core.ErrorCreate("failed to update AWS ECS service",
			"aws_ecs_service_update_error",
			map[string]interface{}{"service_name_str": p_service_name_str,},
			err, "gf_aws", p_runtime_sys)
		return gf_err
	}

	fmt.Println(result)


	return nil
}