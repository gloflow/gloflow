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
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------
func AWS_EC2__describe_instances__by_tags(p_tags_lst []map[string]string,
	p_runtime_sys *gf_core.Runtime_sys) ([]*ec2.Instance, *gf_core.Gf_error) {



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

		gf_err := gf_core.Error__create("failed to describe ec2 instances with specified tags",
			"aws_ec2_instances_describe_error",
			map[string]interface{}{"tags_lst": p_tags_lst,},
			err, "gf_aws", p_runtime_sys)
		return nil, gf_err
	}


	instances_lst := []*ec2.Instance{}
	for _, r := range result.Reservations {


		instances_lst = append(instances_lst, r.Instances...)
	}


	return instances_lst, nil
}