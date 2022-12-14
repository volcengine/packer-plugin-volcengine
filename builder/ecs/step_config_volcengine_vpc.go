package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepConfigVolcengineVpc struct {
	vpcId               string
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepConfigVolcengineVpc) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	if s.VolcengineEcsConfig.VpcId == "" && s.VolcengineEcsConfig.SubnetId != "" {
		input := vpc.DescribeSubnetsInput{
			SubnetIds: volcengine.StringSlice([]string{s.VolcengineEcsConfig.SubnetId}),
		}
		out, err := client.VpcClient.DescribeSubnetsWithContext(ctx, &input)
		if err != nil || len(out.Subnets) == 0 {
			return Halt(stateBag, err, fmt.Sprintf("Error query Subnet with id %s", s.VolcengineEcsConfig.VpcId))
		}

		s.VolcengineEcsConfig.VpcId = *out.Subnets[0].VpcId

		ui.Say(fmt.Sprintf("Using existing Vpc id is %s", s.VolcengineEcsConfig.VpcId))
		return multistep.ActionContinue
	}

	if s.VolcengineEcsConfig.VpcId != "" {
		//valid
		input := vpc.DescribeVpcsInput{
			VpcIds: volcengine.StringSlice([]string{s.VolcengineEcsConfig.VpcId}),
		}
		out, err := client.VpcClient.DescribeVpcsWithContext(ctx, &input)
		if err != nil || len(out.Vpcs) == 0 {
			return Halt(stateBag, err, fmt.Sprintf("Error query Vpc with id %s", s.VolcengineEcsConfig.VpcId))
		}
		ui.Say(fmt.Sprintf("Using existing Vpc id is %s", s.VolcengineEcsConfig.VpcId))
		return multistep.ActionContinue
	}
	//create new vpc
	if s.VolcengineEcsConfig.VpcName == "" {
		s.VolcengineEcsConfig.VpcName = defaultVpcName
	}
	if s.VolcengineEcsConfig.VpcCidrBlock == "" {
		s.VolcengineEcsConfig.VpcCidrBlock = defaultVpcCidr
	}
	var dnses []string
	if s.VolcengineEcsConfig.DNS1 != "" {
		dnses = append(dnses, s.VolcengineEcsConfig.DNS1)
	}
	if s.VolcengineEcsConfig.DNS2 != "" {
		dnses = append(dnses, s.VolcengineEcsConfig.DNS2)
	}
	ui.Say(fmt.Sprintf("Creating new Vpc with name %s cidr %s", s.VolcengineEcsConfig.VpcName,
		s.VolcengineEcsConfig.VpcCidrBlock))

	input := vpc.CreateVpcInput{
		VpcName:   volcengine.String(s.VolcengineEcsConfig.VpcName),
		CidrBlock: volcengine.String(s.VolcengineEcsConfig.VpcCidrBlock),
	}
	if len(dnses) > 0 {
		input.DnsServers = volcengine.StringSlice(dnses)
	}
	output, err := client.VpcClient.CreateVpcWithContext(ctx, &input)
	if err != nil {
		return Halt(stateBag, err, "Error creating new Vpc")
	}
	s.vpcId = *output.VpcId
	s.VolcengineEcsConfig.VpcId = *output.VpcId

	_, err = WaitVpcStatus(stateBag, s.VolcengineEcsConfig.VpcId, "Available")
	if err != nil {
		return Halt(stateBag, err, "Error creating new Vpc")
	}
	return multistep.ActionContinue
}

func (s *stepConfigVolcengineVpc) Cleanup(stateBag multistep.StateBag) {
	if s.vpcId != "" {
		ui := stateBag.Get("ui").(packer.Ui)
		client := stateBag.Get("client").(*VolcengineClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Vpc with Id %s ", s.vpcId))
		_, err := WaitVpcStatus(stateBag, s.vpcId, "Available")
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Vpc %s", err))
			return
		}
		input := vpc.DeleteVpcInput{
			VpcId: volcengine.String(s.vpcId),
		}
		_, err = client.VpcClient.DeleteVpc(&input)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete vpc %s", err))
		}
	}
}
