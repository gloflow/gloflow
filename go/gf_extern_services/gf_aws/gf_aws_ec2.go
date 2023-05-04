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
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------------------
// EC2_SCALE_AUTSCALING_GROUP

func EC2scaleAutoscalingGroup(pAutscalingGroupNameStr string,
	pDesiredCapacityInt int,
	pRuntimeSys         *gf_core.RuntimeSys) *gf_core.GFerror {
	
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-east-1")},
		Profile: "default",
	})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create an AWS session object",
			"aws_session_create",
			map[string]interface{}{},
			err, "gf_aws", pRuntimeSys)
		return gfErr
	}

	svc := autoscaling.New(sess)

	//-----------------------
	// UPDATE
	_, err = svc.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(pAutscalingGroupNameStr),
		DesiredCapacity:      aws.Int64(int64(pDesiredCapacityInt)),
	})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update an EC2 Autoscaling Group capacity",
			"aws_ec2_autoscaling_scale_error",
			map[string]interface{}{
				"autoscaling_group_name_str": pAutscalingGroupNameStr,
				"desired_capacity_int":       pDesiredCapacityInt,
			},
			err, "gf_aws", pRuntimeSys)
		return gfErr
	}

	//-----------------------
	return nil
}

//-------------------------------------------------------------

func EC2getInfoOnAutoscalingGroup(pAutscalingGroupNameStr string,
	pRuntimeSys *gf_core.RuntimeSys) (int, *gf_core.GFerror) {
	
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-east-1")},
		Profile: "default",
	})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create an AWS session object",
			"aws_session_create",
			map[string]interface{}{},
			err, "gf_aws", pRuntimeSys)
		return 0, gfErr
	}

	svc := autoscaling.New(sess)

	//-----------------------
	// GET_INFO
	result, err := svc.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(pAutscalingGroupNameStr)},
	})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to update an EC2 Autoscaling Group capacity",
			"aws_ec2_autoscaling_describe_error",
			map[string]interface{}{
				"autoscaling_group_name_str": pAutscalingGroupNameStr,
			},
			err, "gf_aws", pRuntimeSys)
		return 0, gfErr
	}

	instancesNumInt := len(result.AutoScalingGroups[0].Instances)

	//-----------------------
	return instancesNumInt, nil
}

//-------------------------------------------------------------

func EC2describeInstancesByTags(p_tags_lst []map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) ([]*ec2.Instance, *gf_core.GFerror) {


	svc := ec2.New(session.New())



	filters_lst := []*ec2.Filter{}
	for _, t_map := range p_tags_lst {
		
		for k, v := range t_map {
			filters_lst = append(filters_lst, &ec2.Filter{
				Name: aws.String(fmt.Sprintf("tag:%s", k)),
				Values: []*string{
					aws.String(v),
				},
			})
		}
	}
	
	input := &ec2.DescribeInstancesInput{
		Filters: filters_lst,
	}


	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

		gf_err := gf_core.ErrorCreate("failed to describe ec2 instances with specified tags",
			"aws_ec2_instances_describe_error",
			map[string]interface{}{"tags_lst": p_tags_lst,},
			err, "gf_aws", pRuntimeSys)
		return nil, gf_err
	}




	instances_lst := []*ec2.Instance{}
	for _, r := range result.Reservations {


		instances_lst = append(instances_lst, r.Instances...)
	}


	return instances_lst, nil
}