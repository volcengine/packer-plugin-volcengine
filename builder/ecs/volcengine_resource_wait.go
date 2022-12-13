package ecs

import (
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func WaitEipStatus(stateBag multistep.StateBag, eipId, status string) (*vpc.DescribeEipAddressesOutput, error) {
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	out, err := WaitFor(&WaitForParam{
		Call: func() (interface{}, error) {
			input := vpc.DescribeEipAddressesInput{
				AllocationIds: volcengine.StringSlice([]string{eipId}),
			}
			return client.VpcClient.DescribeEipAddresses(&input)
		},
		Process: func(i interface{}, err error) ProcessResult {
			if err != nil {
				return WaitForRetry
			}
			output := i.(vpc.DescribeEipAddressesOutput)
			if len(output.EipAddresses) < 1 {
				return WaitForRetry
			}
			if *output.EipAddresses[0].Status == status {
				return WaitForSuccess
			} else {
				return WaitForRetry
			}
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    50,
	})
	if out != nil {
		return out.(*vpc.DescribeEipAddressesOutput), err
	}
	return nil, err
}

func WaitVpcStatus(stateBag multistep.StateBag, vpcId, status string) (*vpc.DescribeVpcsOutput, error) {
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	out, err := WaitFor(&WaitForParam{
		Call: func() (interface{}, error) {
			input := vpc.DescribeVpcsInput{
				VpcIds: volcengine.StringSlice([]string{vpcId}),
			}
			return client.VpcClient.DescribeVpcs(&input)
		},
		Process: func(i interface{}, err error) ProcessResult {
			if err != nil {
				return WaitForRetry
			}
			output := i.(vpc.DescribeVpcsOutput)
			if len(output.Vpcs) < 1 {
				return WaitForRetry
			}
			if *output.Vpcs[0].Status == status {
				return WaitForSuccess
			} else {
				return WaitForRetry
			}
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    50,
	})
	if out != nil {
		return out.(*vpc.DescribeVpcsOutput), err
	}
	return nil, err
}

func WaitSubnetStatus(stateBag multistep.StateBag, subnetId, status string) (*vpc.DescribeSubnetsOutput, error) {
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	out, err := WaitFor(&WaitForParam{
		Call: func() (interface{}, error) {
			input := vpc.DescribeSubnetsInput{
				SubnetIds: volcengine.StringSlice([]string{subnetId}),
			}
			return client.VpcClient.DescribeSubnets(&input)
		},
		Process: func(i interface{}, err error) ProcessResult {
			if err != nil {
				return WaitForRetry
			}
			output := i.(vpc.DescribeSubnetsOutput)
			if len(output.Subnets) < 1 {
				return WaitForRetry
			}
			if *output.Subnets[0].Status == status {
				return WaitForSuccess
			} else {
				return WaitForRetry
			}
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    50,
	})
	if out != nil {
		return out.(*vpc.DescribeSubnetsOutput), err
	}
	return nil, err
}

func WaitImageStatus(stateBag multistep.StateBag, imageId, status string) (*ecs.DescribeImagesOutput, error) {
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	out, err := WaitFor(&WaitForParam{
		Call: func() (interface{}, error) {
			input := ecs.DescribeImagesInput{
				ImageIds: volcengine.StringSlice([]string{imageId}),
			}
			return client.EcsClient.DescribeImages(&input)
		},
		Process: func(i interface{}, err error) ProcessResult {
			if err != nil {
				return WaitForRetry
			}
			output := i.(ecs.DescribeImagesOutput)
			if len(output.Images) < 1 {
				return WaitForRetry
			}
			if *output.Images[0].Status == status {
				return WaitForSuccess
			} else {
				return WaitForRetry
			}
		},
		RetryInterval: 30 * time.Second,
		RetryTimes:    100,
	})
	if out != nil {
		return out.(*ecs.DescribeImagesOutput), err
	}
	return nil, err
}

func WaitEcsStatus(stateBag multistep.StateBag, ecsId, status string) (*ecs.DescribeInstancesOutput, error) {
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	out, err := WaitFor(&WaitForParam{
		Call: func() (interface{}, error) {
			input := ecs.DescribeInstancesInput{
				InstanceIds: volcengine.StringSlice([]string{ecsId}),
			}
			return client.EcsClient.DescribeInstances(&input)
		},
		Process: func(i interface{}, err error) ProcessResult {
			if err != nil {
				return WaitForRetry
			}
			output := i.(ecs.DescribeInstancesOutput)
			if len(output.Instances) < 1 {
				return WaitForRetry
			}
			if *output.Instances[0].Status == status {
				return WaitForSuccess
			} else {
				return WaitForRetry
			}
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    50,
	})
	if out != nil {
		return out.(*ecs.DescribeInstancesOutput), err
	}
	return nil, err
}

func WaitSgClean(stateBag multistep.StateBag, sgId string) error {
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	_, err := WaitFor(&WaitForParam{
		Call: func() (interface{}, error) {
			input := vpc.DescribeNetworkInterfacesInput{
				SecurityGroupId: volcengine.String(sgId),
			}
			return client.VpcClient.DescribeNetworkInterfaces(&input)
		},
		Process: func(i interface{}, err error) ProcessResult {
			if err != nil {
				return WaitForRetry
			}
			output := i.(vpc.DescribeNetworkInterfacesOutput)
			if len(output.NetworkInterfaceSets) == 0 {
				return WaitForSuccess
			}
			return WaitForRetry
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    50,
	})
	return err
}
