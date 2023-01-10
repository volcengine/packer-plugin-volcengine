package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepConfigVolcengineEcs struct {
	ecsId               string
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepConfigVolcengineEcs) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	//create new ecs
	if s.VolcengineEcsConfig.InstanceName == "" {
		s.VolcengineEcsConfig.InstanceName = defaultEcsName
	}
	//volume-sys
	var volumes []*ecs.VolumeForRunInstancesInput
	volumes = append(volumes, &ecs.VolumeForRunInstancesInput{
		DeleteWithInstance: volcengine.String("true"),
		VolumeType:         volcengine.String(s.VolcengineEcsConfig.SystemDiskType),
		Size:               volcengine.Int32(s.VolcengineEcsConfig.SystemDiskSize),
	})
	//volume-data
	for _, vd := range s.VolcengineEcsConfig.DataDisks {
		volumes = append(volumes, &ecs.VolumeForRunInstancesInput{
			DeleteWithInstance: volcengine.String("true"),
			VolumeType:         volcengine.String(vd.DataDiskType),
			Size:               volcengine.Int32(vd.DataDiskSize),
		})
	}
	//net
	var networks []*ecs.NetworkInterfaceForRunInstancesInput
	networks = append(networks, &ecs.NetworkInterfaceForRunInstancesInput{
		SubnetId:         volcengine.String(s.VolcengineEcsConfig.SubnetId),
		SecurityGroupIds: volcengine.StringSlice([]string{s.VolcengineEcsConfig.SecurityGroupId}),
	})

	input := ecs.RunInstancesInput{
		InstanceTypeId:     volcengine.String(s.VolcengineEcsConfig.InstanceType),
		ImageId:            volcengine.String(s.VolcengineEcsConfig.SourceImageId),
		ZoneId:             volcengine.String(s.VolcengineEcsConfig.AvailabilityZone),
		InstanceName:       volcengine.String(s.VolcengineEcsConfig.InstanceName),
		InstanceChargeType: volcengine.String("PostPaid"),
		Volumes:            volumes,
		NetworkInterfaces:  networks,
	}
	//password/ssh/key
	if s.VolcengineEcsConfig.Comm.SSHKeyPairName != "" {
		input.KeyPairName = volcengine.String(s.VolcengineEcsConfig.Comm.SSHKeyPairName)
	} else {
		if s.VolcengineEcsConfig.Comm.SSHPassword != "" {
			input.Password = volcengine.String(s.VolcengineEcsConfig.Comm.SSHPassword)
		} else if s.VolcengineEcsConfig.Comm.WinRMPassword != "" {
			input.Password = volcengine.String(s.VolcengineEcsConfig.Comm.WinRMPassword)
		} else {
			return Halt(stateBag, fmt.Errorf("no keypair or password set"), "Error creating new Ecs")
		}
	}

	//userdata
	if s.VolcengineEcsConfig.UserData != "" {
		input.UserData = volcengine.String(s.VolcengineEcsConfig.UserData)
	}

	//hpc
	if s.VolcengineEcsConfig.HpcClusterId != "" {
		input.HpcClusterId = volcengine.String(s.VolcengineEcsConfig.HpcClusterId)
	}

	ui.Say(fmt.Sprintf("Creating new ecs with name %s", s.VolcengineEcsConfig.InstanceName))

	output, err := client.EcsClient.RunInstancesWithContext(ctx, &input)
	if err != nil {
		return Halt(stateBag, err, "Error creating new Ecs")
	}
	s.ecsId = *output.InstanceIds[0]

	out, err := WaitEcsStatus(stateBag, s.ecsId, "RUNNING")
	if err != nil {
		return Halt(stateBag, err, "Error creating new Ecs")
	}
	stateBag.Put("instanceId", s.ecsId)
	stateBag.Put("PrivateIp", *out.Instances[0].NetworkInterfaces[0].PrimaryIpAddress)
	return multistep.ActionContinue
}

func (s *stepConfigVolcengineEcs) Cleanup(stateBag multistep.StateBag) {
	if s.ecsId != "" {
		ui := stateBag.Get("ui").(packer.Ui)
		client := stateBag.Get("client").(*VolcengineClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Ecs with Id %s ", s.ecsId))
		input := ecs.DeleteInstanceInput{
			InstanceId: volcengine.String(s.ecsId),
		}
		_, err := client.EcsClient.DeleteInstance(&input)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Ecs %s", err))
		}
	}
}
