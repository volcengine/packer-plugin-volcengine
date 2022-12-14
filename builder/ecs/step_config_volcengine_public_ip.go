package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepConfigVolcenginePublicIp struct {
	eipId               string
	isCreate            bool
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepConfigVolcenginePublicIp) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	instanceId := stateBag.Get("instanceId").(string)
	if s.VolcengineEcsConfig.AssociatePublicIpAddress {
		if s.VolcengineEcsConfig.PublicIpId != "" {
			//valid
			input := vpc.DescribeEipAddressesInput{
				AllocationIds: volcengine.StringSlice([]string{s.VolcengineEcsConfig.PublicIpId}),
			}
			out, err := client.VpcClient.DescribeEipAddressesWithContext(ctx, &input)
			if err != nil || len(out.EipAddresses) == 0 {
				return Halt(stateBag, err, fmt.Sprintf("Error query Eip with id %s", s.VolcengineEcsConfig.PublicIpId))
			}
			s.eipId = s.VolcengineEcsConfig.PublicIpId
			stateBag.Put("PublicIp", *out.EipAddresses[0].EipAddress)
			ui.Say(fmt.Sprintf("Using existing Public IP id is %s", s.VolcengineEcsConfig.PublicIpId))
		} else {
			ui.Say("Creating new Eip...")
			//create new eip
			if s.VolcengineEcsConfig.PublicIpBandWidth < 1 {
				s.VolcengineEcsConfig.PublicIpBandWidth = 1
			}
			input := vpc.AllocateEipAddressInput{
				Bandwidth:   volcengine.Int64(s.VolcengineEcsConfig.PublicIpBandWidth),
				BillingType: volcengine.Int64(3),
			}
			output, err := client.VpcClient.AllocateEipAddressWithContext(ctx, &input)
			if err != nil {
				return Halt(stateBag, err, "Error creating new Eip")
			}
			s.eipId = *output.AllocationId
			s.isCreate = true
			out, err := WaitEipStatus(stateBag, s.eipId, "Available")
			if err != nil {
				return Halt(stateBag, err, "Error creating new eip")
			}
			stateBag.Put("PublicIp", *out.EipAddresses[0].EipAddress)
		}
		//set sg rule
		ui.Say(fmt.Sprintf("Authorize SecurityGroup %s Rule", s.VolcengineEcsConfig.SecurityGroupId))
		input2 := vpc.AuthorizeSecurityGroupIngressInput{
			SecurityGroupId: volcengine.String(s.VolcengineEcsConfig.SecurityGroupId),
			Protocol:        volcengine.String("tcp"),
			PortStart:       volcengine.Int64(22),
			PortEnd:         volcengine.Int64(22),
			CidrIp:          volcengine.String("0.0.0.0/0"),
		}
		_, err := client.VpcClient.AuthorizeSecurityGroupIngressWithContext(ctx, &input2)
		if err != nil {
			return Halt(stateBag, err, "Error Authorize SecurityGroup Rule")
		}

		ui.Say(fmt.Sprintf("Associate  Eip %s to ecs %s", s.eipId, instanceId))
		//bind
		input1 := vpc.AssociateEipAddressInput{
			InstanceId:   volcengine.String(instanceId),
			AllocationId: volcengine.String(s.eipId),
			InstanceType: volcengine.String("EcsInstance"),
		}
		_, err = client.VpcClient.AssociateEipAddress(&input1)
		if err != nil {
			return Halt(stateBag, err, "Error binding new eip")
		}
		_, err = WaitEipStatus(stateBag, s.eipId, "Attached")
		if err != nil {
			return Halt(stateBag, err, "Error binding new eip")
		}
		return multistep.ActionContinue
	}

	return multistep.ActionContinue
}

func (s *stepConfigVolcenginePublicIp) Cleanup(stateBag multistep.StateBag) {
	if s.VolcengineEcsConfig.AssociatePublicIpAddress {
		ui := stateBag.Get("ui").(packer.Ui)
		client := stateBag.Get("client").(*VolcengineClientWrapper)
		ui.Say(fmt.Sprintf("Disassociate Eip with Id %s ", s.eipId))

		//unbind
		_, err := WaitEipStatus(stateBag, s.eipId, "Attached")
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Eip %s", err))
			return
		}
		unbind := vpc.DisassociateEipAddressInput{
			AllocationId: volcengine.String(s.eipId),
		}

		_, err = client.VpcClient.DisassociateEipAddress(&unbind)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Eip %s", err))
			return
		}
		_, err = WaitEipStatus(stateBag, s.eipId, "Available")
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Eip %s", err))
			return
		}
		if s.isCreate {
			ui.Say(fmt.Sprintf("Delete Eip with Id %s ", s.eipId))
			input := vpc.ReleaseEipAddressInput{
				AllocationId: volcengine.String(s.eipId),
			}
			_, err = client.VpcClient.ReleaseEipAddress(&input)
			if err != nil {
				ui.Error(fmt.Sprintf("Error delete Eip %s", err))
			}
		}
	}
}
