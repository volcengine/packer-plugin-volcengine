package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepConfigVolcengineEcsStop struct {
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepConfigVolcengineEcsStop) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	instanceId := stateBag.Get("instanceId").(string)
	ui.Say(fmt.Sprintf("stoping ecs instance"))
	input := ecs.StopInstanceInput{
		InstanceId: volcengine.String(instanceId),
	}
	_, err := client.EcsClient.StopInstanceWithContext(ctx, &input)
	if err != nil {
		return Halt(stateBag, err, "Error stop ecs instance")
	}
	_, err = WaitEcsStatus(stateBag, instanceId, "STOPPED")
	if err != nil {
		return Halt(stateBag, err, "Error stop ecs instance")
	}
	return multistep.ActionContinue
}

func (stepConfigVolcengineEcsStop) Cleanup(bag multistep.StateBag) {
}
