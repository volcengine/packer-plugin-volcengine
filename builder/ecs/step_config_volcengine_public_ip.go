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
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepConfigVolcenginePublicIp) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	instanceId := stateBag.Get("instanceId").(string)
	ui.Say(fmt.Sprintf("Creating new Eip "))
	if s.VolcengineEcsConfig.AssociatePublicIpAddress {
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
		_, err = WaitEipStatus(stateBag, s.eipId, "Available")
		if err != nil {
			return Halt(stateBag, err, "Error creating new eip")
		}
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
		_, err = WaitEipStatus(stateBag, s.eipId, "Available")
		if err != nil {
			return Halt(stateBag, err, "Error binding new eip")
		}
	}

	return multistep.ActionContinue
}

func (s *stepConfigVolcenginePublicIp) Cleanup(stateBag multistep.StateBag) {
	if s.eipId != "" {
		ui := stateBag.Get("ui").(packer.Ui)
		client := stateBag.Get("client").(*VolcengineClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Eip with Id %s ", s.eipId))

		//unbind
		_, err := WaitVpcStatus(stateBag, s.eipId, "Available")
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
		_, err = WaitVpcStatus(stateBag, s.eipId, "Available")
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Eip %s", err))
			return
		}
		input := vpc.ReleaseEipAddressInput{
			AllocationId: volcengine.String(s.eipId),
		}
		_, err = client.VpcClient.ReleaseEipAddress(&input)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Eip %s", err))
		}
	}
}
