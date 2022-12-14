package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepConfigVolcengineSubnet struct {
	subnetId            string
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepConfigVolcengineSubnet) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	if s.VolcengineEcsConfig.SubnetId != "" {
		//valid
		input := vpc.DescribeSubnetsInput{
			SubnetIds: volcengine.StringSlice([]string{s.VolcengineEcsConfig.SubnetId}),
		}
		out, err := client.VpcClient.DescribeSubnetsWithContext(ctx, &input)
		if err != nil || len(out.Subnets) == 0 {
			return Halt(stateBag, err, fmt.Sprintf("Error query Subnet with id %s", s.VolcengineEcsConfig.SubnetId))
		}

		if *out.Subnets[0].VpcId != s.VolcengineEcsConfig.VpcId {
			return Halt(stateBag, fmt.Errorf(fmt.Sprintf("Subnet id %s vpc not match",
				s.VolcengineEcsConfig.SubnetId)), "")
		}
		s.VolcengineEcsConfig.AvailabilityZone = *out.Subnets[0].ZoneId
		ui.Say(fmt.Sprintf("Using existing Subnet id is %s", s.VolcengineEcsConfig.SubnetId))
		return multistep.ActionContinue
	}
	//create new subnet
	if s.VolcengineEcsConfig.SubnetName == "" {
		s.VolcengineEcsConfig.SubnetName = defaultSubnetName
	}
	if s.VolcengineEcsConfig.SubnetCidrBlock == "" {
		s.VolcengineEcsConfig.SubnetCidrBlock = defaultSubnetCidr
	}
	ui.Say(fmt.Sprintf("Creating new Subnet with name %s cidr %s", s.VolcengineEcsConfig.SubnetName,
		s.VolcengineEcsConfig.SubnetCidrBlock))

	input := vpc.CreateSubnetInput{
		VpcId:      volcengine.String(s.VolcengineEcsConfig.VpcId),
		SubnetName: volcengine.String(s.VolcengineEcsConfig.SubnetName),
		CidrBlock:  volcengine.String(s.VolcengineEcsConfig.SubnetCidrBlock),
	}
	if s.VolcengineEcsConfig.AvailabilityZone != "" {
		input.ZoneId = volcengine.String(s.VolcengineEcsConfig.AvailabilityZone)
	} else {
		out, err := client.EcsClient.DescribeZones(&ecs.DescribeZonesInput{})
		if err != nil {
			return Halt(stateBag, err, "Error creating new Subnet")
		}
		input.ZoneId = out.Zones[0].ZoneId
	}
	output, err := client.VpcClient.CreateSubnetWithContext(ctx, &input)
	if err != nil {
		return Halt(stateBag, err, "Error creating new Subnet")
	}
	s.subnetId = *output.SubnetId
	s.VolcengineEcsConfig.SubnetId = *output.SubnetId
	subnets, err := WaitSubnetStatus(stateBag, s.VolcengineEcsConfig.SubnetId, "Available")
	if err != nil {
		return Halt(stateBag, err, "Error creating new Subnet")
	}
	s.VolcengineEcsConfig.AvailabilityZone = *subnets.Subnets[0].ZoneId

	_, err = WaitVpcStatus(stateBag, s.VolcengineEcsConfig.VpcId, "Available")
	if err != nil {
		return Halt(stateBag, err, "Error creating new Subnet")
	}

	return multistep.ActionContinue
}

func (s *stepConfigVolcengineSubnet) Cleanup(stateBag multistep.StateBag) {
	if s.subnetId != "" {
		ui := stateBag.Get("ui").(packer.Ui)
		client := stateBag.Get("client").(*VolcengineClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Subnet with Id %s ", s.subnetId))
		_, err := WaitSubnetStatus(stateBag, s.subnetId, "Available")
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Subnet %s", err))
			return
		}
		input := vpc.DeleteSubnetInput{
			SubnetId: volcengine.String(s.subnetId),
		}
		_, err = client.VpcClient.DeleteSubnet(&input)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete subnet %s", err))
		}
	}
}
