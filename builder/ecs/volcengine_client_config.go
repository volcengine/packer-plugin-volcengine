package ecs

import (
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

type VolcengineClientConfig struct {
	VolcengineAuthenticationConfig `mapstructure:",squash"`
	client                         *VolcengineClientWrapper
}

func (v *VolcengineClientConfig) Client(stateBag *multistep.BasicStateBag) *VolcengineClientWrapper {
	if v.client != nil {
		stateBag.Put("client", v.client)
		return v.client
	}

	config := volcengine.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(v.VolcengineAccessKey, v.VolcengineSecretKey, v.VolcengineSessionKey)).
		WithDisableSSL(*v.VolcengineDisableSSL).
		WithRegion(v.VolcengineRegion)

	if v.VolcengineEndpoint != "" {
		config.WithEndpoint(v.VolcengineEndpoint)
	}

	sess, _ := session.NewSession(config)

	v.client = &VolcengineClientWrapper{
		VpcClient: vpc.New(sess),
		EcsClient: ecs.New(sess),
	}
	stateBag.Put("client", v.client)

	return v.client
}
