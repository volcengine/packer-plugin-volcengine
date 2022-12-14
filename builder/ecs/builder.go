//go:generate packer-sdc mapstructure-to-hcl2 -type Config,VolcengineDataDiskConfig
package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	VolcengineClientConfig `mapstructure:",squash"`
	VolcengineEcsConfig    `mapstructure:",squash"`

	ctx interpolate.Context
}

const BuilderId = "volcengine.ecs"

type Builder struct {
	runner multistep.Runner
	config Config
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"run_command",
			},
		},
	}, raws...)
	b.config.ctx.EnableEnv = true
	if err != nil {
		return nil, nil, err
	}
	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.VolcengineClientConfig.Prepare(&b.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	packer.LogSecretFilter.Set(b.config.VolcengineAccessKey, b.config.VolcengineSecretKey, b.config.VolcengineSessionKey)
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	stateBag := new(multistep.BasicStateBag)
	b.config.Client(stateBag)

	stateBag.Put("hook", hook)
	stateBag.Put("ui", ui)
	stateBag.Put("config", &b.config)

	if b.config.VolcengineEcsConfig.Comm.Type == "" {
		b.config.VolcengineEcsConfig.Comm.Type = "ssh"
		b.config.VolcengineEcsConfig.Comm.SSHPort = 22
		b.config.VolcengineEcsConfig.Comm.WinRMPort = 5895
	}

	//steps...
	var steps []multistep.Step
	steps = []multistep.Step{
		&stepValidVolcengineImage{
			SourceImageId: b.config.SourceImageId,
		},
		&stepConfigVolcengineKeyPair{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&stepConfigVolcengineVpc{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&stepConfigVolcengineSubnet{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&stepConfigVolcengineSg{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&stepConfigVolcengineEcs{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&stepConfigVolcenginePublicIp{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&communicator.StepConnect{
			Config:    &b.config.VolcengineEcsConfig.Comm,
			Host:      SSHHost(),
			SSHConfig: b.config.VolcengineEcsConfig.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&stepConfigVolcengineEcsStop{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
		&stepCreateVolcengineImage{
			VolcengineEcsConfig: &b.config.VolcengineEcsConfig,
		},
	}

	// Run...
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, stateBag)

	// If there was an error, return that
	if err, ok := stateBag.GetOk("error"); ok {
		ui.Say(fmt.Sprintf("find some error %v ", err))
		return nil, err.(error)
	}

	artifact := &Artifact{
		VolcengineImageId: stateBag.Get("TargetImageId").(string),
		BuilderIdValue:    BuilderId,
		Client:            b.config.client,
	}
	return artifact, nil
}
