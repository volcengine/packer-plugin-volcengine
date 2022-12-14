package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepValidVolcengineImage struct {
	SourceImageId string
}

func (s *stepValidVolcengineImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	req := ecs.DescribeImagesInput{
		ImageIds: volcengine.StringSlice([]string{s.SourceImageId}),
	}
	_, err := client.EcsClient.DescribeImagesWithContext(ctx, &req)
	if err != nil {
		return Halt(stateBag, err, "Error querying volcengine image")
	}

	ui.Message(fmt.Sprintf("Found volcengine image ID: %s", s.SourceImageId))
	return multistep.ActionContinue
}

func (s *stepValidVolcengineImage) Cleanup(stateBag multistep.StateBag) {
}
