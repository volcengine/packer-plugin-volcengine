//go:generate packer-sdc struct-markdown
package ecs

import "github.com/hashicorp/packer-plugin-sdk/communicator"

type VolcengineDataDiskConfig struct {
	DataDiskType string `mapstructure:"data_disk_type" required:"true"`
	DataDiskSize int32  `mapstructure:"data_disk_size" required:"true"`
}

type VolcengineEcsConfig struct {
	//vpc
	VpcId        string `mapstructure:"vpc_id" required:"false"`
	VpcName      string `mapstructure:"vpc_name" required:"false"`
	VpcCidrBlock string `mapstructure:"vpc_cidr_block" required:"false"`
	DNS1         string `mapstructure:"dns1" required:"false"`
	DNS2         string `mapstructure:"dns2" required:"false"`
	//subnet
	SubnetId         string `mapstructure:"subnet_id" required:"false"`
	SubnetName       string `mapstructure:"subnet_name" required:"false"`
	SubnetCidrBlock  string `mapstructure:"subnet_cidr_block" required:"false"`
	AvailabilityZone string `mapstructure:"availability_zone" required:"false"`
	//sg
	SecurityGroupId   string `mapstructure:"security_group_id" required:"false"`
	SecurityGroupName string `mapstructure:"security_group_name" required:"false"`
	//eip
	PublicIpId               string `mapstructure:"public_ip_id" required:"false"`
	AssociatePublicIpAddress bool   `mapstructure:"associate_public_ip_address" required:"false"`
	PublicIpBandWidth        int64  `mapstructure:"public_ip_band_width" required:"false"`
	//ssh
	Comm communicator.Config `mapstructure:",squash"`
	//ecs
	InstanceType    string                     `mapstructure:"instance_type" required:"true"`
	SourceImageId   string                     `mapstructure:"source_image_id" required:"true"`
	TargetImageName string                     `mapstructure:"target_image_name" required:"true"`
	SystemDiskType  string                     `mapstructure:"system_disk_type" required:"true"`
	SystemDiskSize  int32                      `mapstructure:"system_disk_size" required:"true"`
	DataDisks       []VolcengineDataDiskConfig `mapstructure:"data_disks" required:"false"`
	InstanceName    string                     `mapstructure:"instance_name" required:"false"`
	UserData        string                     `mapstructure:"user_data" required:"false"`
	//hpc
	HpcClusterId string `mapstructure:"hpc_cluster_id" required:"false"`
}
