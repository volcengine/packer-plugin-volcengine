package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepConfigVolcengineKeyPair struct {
	VolcengineEcsConfig *VolcengineEcsConfig
	isCreate            bool
}

func (s *stepConfigVolcengineKeyPair) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	//ssh private key file
	if s.VolcengineEcsConfig.Comm.SSHPrivateKeyFile != "" {
		if s.VolcengineEcsConfig.Comm.SSHKeyPairName == "" {
			return Halt(stateBag,
				fmt.Errorf(fmt.Sprintf("ssh_keypair_name is empty")), "")
		}
		ui.Say("Using existing SSH private key")
		privateKeyBytes, err := s.VolcengineEcsConfig.Comm.ReadSSHPrivateKeyFile()
		if err != nil {
			stateBag.Put("error", err)
			return multistep.ActionHalt
		}
		s.VolcengineEcsConfig.Comm.SSHPrivateKey = privateKeyBytes
		return multistep.ActionContinue
	}

	if s.VolcengineEcsConfig.Comm.SSHAgentAuth && s.VolcengineEcsConfig.Comm.SSHKeyPairName == "" {
		ui.Say("Using SSH Agent with key pair in source image")
		return multistep.ActionContinue
	}

	if s.VolcengineEcsConfig.Comm.SSHAgentAuth && s.VolcengineEcsConfig.Comm.SSHKeyPairName != "" {
		ui.Say(fmt.Sprintf("Using SSH Agent for existing key pair %s", s.VolcengineEcsConfig.Comm.SSHKeyPairName))
		return multistep.ActionContinue
	}

	if s.VolcengineEcsConfig.Comm.SSHTemporaryKeyPairName == "" {
		ui.Say("Not using temporary keypair")
		s.VolcengineEcsConfig.Comm.SSHKeyPairName = ""
		return multistep.ActionContinue
	}
	//create new key_pair
	ui.Say(fmt.Sprintf("Using SSH Agent for create new key pair %s", s.VolcengineEcsConfig.Comm.SSHTemporaryKeyPairName))
	input := ecs.CreateKeyPairInput{
		KeyPairName: volcengine.String(s.VolcengineEcsConfig.Comm.SSHTemporaryKeyPairName),
	}
	output, err := client.EcsClient.CreateKeyPairWithContext(ctx, &input)
	if err != nil {
		return Halt(stateBag, err, "Error creating new keypair")
	}
	s.VolcengineEcsConfig.Comm.SSHKeyPairName = *output.KeyPairName
	s.VolcengineEcsConfig.Comm.SSHPrivateKey = []byte(*output.PrivateKey)
	s.isCreate = true
	return multistep.ActionContinue
}

func (s *stepConfigVolcengineKeyPair) Cleanup(stateBag multistep.StateBag) {
	if s.isCreate {
		ui := stateBag.Get("ui").(packer.Ui)
		client := stateBag.Get("client").(*VolcengineClientWrapper)
		ui.Say(fmt.Sprintf("Deleting temporary keypair %s ", s.VolcengineEcsConfig.Comm.SSHKeyPairName))
		input := ecs.DeleteKeyPairsInput{
			KeyPairNames: volcengine.StringSlice([]string{s.VolcengineEcsConfig.Comm.SSHKeyPairName}),
		}
		_, err := client.EcsClient.DeleteKeyPairs(&input)
		if err != nil {
			ui.Error(fmt.Sprintf("Error cleaning up keypair. Please delete the key manually: name = %s",
				s.VolcengineEcsConfig.Comm.SSHKeyPairName))
		}
	}
}
