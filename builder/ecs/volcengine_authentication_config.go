//go:generate packer-sdc struct-markdown
package ecs

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type VolcengineAuthenticationConfig struct {
	VolcengineAccessKey  string `mapstructure:"access_key" required:"true"`
	VolcengineSecretKey  string `mapstructure:"secret_key" required:"true"`
	VolcengineSessionKey string `mapstructure:"session_key" required:"false"`
	VolcengineEndpoint   string `mapstructure:"endpoint" required:"false"`
	VolcengineDisableSSL *bool  `mapstructure:"disable_ssl" required:"false"`
	VolcengineRegion     string `mapstructure:"region" required:"true"`
}

func (v *VolcengineAuthenticationConfig) Config() error {
	if v.VolcengineAccessKey == "" {
		v.VolcengineAccessKey = os.Getenv("VOLCENGINE_ACCESS_KEY")
	}
	if v.VolcengineSecretKey == "" {
		v.VolcengineSecretKey = os.Getenv("VOLCENGINE_SECRET_KEY")
	}
	if v.VolcengineSessionKey == "" {
		v.VolcengineSessionKey = os.Getenv("VOLCENGINE_SESSION_KEY")
	}
	if v.VolcengineAccessKey == "" || v.VolcengineSecretKey == "" {
		return fmt.Errorf("VOLCENGINE_ACCESS_KEY and VOLCENGINE_SECRET_KEY must be set in template file or environment variables")
	}
	return nil
}

func (v *VolcengineAuthenticationConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if err := v.Config(); err != nil {
		errs = append(errs, err)
	}

	if v.VolcengineRegion == "" {
		v.VolcengineRegion = os.Getenv("VOLCENGINE_REGION")
	}

	if v.VolcengineEndpoint == "" {
		v.VolcengineEndpoint = os.Getenv("VOLCENGINE_ENDPOINT")
	}

	if v.VolcengineDisableSSL == nil {
		temp := os.Getenv("VOLCENGINE_DISABLE_SSL")
		if temp == "true" {
			v.VolcengineDisableSSL = volcengine.Bool(true)
		} else {
			v.VolcengineDisableSSL = volcengine.Bool(false)
		}
	}

	if v.VolcengineRegion == "" {
		errs = append(errs, fmt.Errorf("region option or VOLCENGINE_REGION must be provided in template file or environment variables"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
