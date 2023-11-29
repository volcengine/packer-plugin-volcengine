Type: `volcengine-ecs`
Artifact BuilderId: `volcengine.ecs`

The `volcengine-ecs` Packer builder plugin provide the capability to build
customized images based on an existing base images.

## Configuration Reference

The following configuration options are available for building Volcengine images.
In addition to the options listed here,
a [communicator](/packer/docs/templates/legacy_json_templates/communicator) can be configured for this builder.

### Required:

<!-- Code generated from the comments of the VolcengineAuthenticationConfig struct in builder/ecs/volcengine_authentication_config.go; DO NOT EDIT MANUALLY -->

- `access_key` (string) - Volcengine Access Key

- `secret_key` (string) - Volcengine Secret Key

- `region` (string) - Volcengine Region

<!-- End of code generated from the comments of the VolcengineAuthenticationConfig struct in builder/ecs/volcengine_authentication_config.go; -->


<!-- Code generated from the comments of the VolcengineEcsConfig struct in builder/ecs/volcengine_ecs_config.go; DO NOT EDIT MANUALLY -->

- `instance_type` (string) - ecs

- `source_image_id` (string) - Source Image Id

- `target_image_name` (string) - Target Image Name

- `system_disk_type` (string) - System Disk Type

- `system_disk_size` (int32) - System Disk Size

<!-- End of code generated from the comments of the VolcengineEcsConfig struct in builder/ecs/volcengine_ecs_config.go; -->


### Optional:

<!-- Code generated from the comments of the VolcengineAuthenticationConfig struct in builder/ecs/volcengine_authentication_config.go; DO NOT EDIT MANUALLY -->

- `session_key` (string) - Volcengine Session Key

- `endpoint` (string) - Volcengine Endpoint

- `disable_ssl` (\*bool) - Volcengine Disable SSL

<!-- End of code generated from the comments of the VolcengineAuthenticationConfig struct in builder/ecs/volcengine_authentication_config.go; -->


<!-- Code generated from the comments of the VolcengineEcsConfig struct in builder/ecs/volcengine_ecs_config.go; DO NOT EDIT MANUALLY -->

- `vpc_id` (string) - vpc

- `vpc_name` (string) - Vpc Name

- `vpc_cidr_block` (string) - Vpc Cidr Block

- `dns1` (string) - DNS 1

- `dns2` (string) - DNS 2

- `subnet_id` (string) - subnet

- `subnet_name` (string) - Subnet Name

- `subnet_cidr_block` (string) - Subnet Cidr Block

- `availability_zone` (string) - Availability Zone

- `security_group_id` (string) - sg

- `security_group_name` (string) - Security Group Name

- `public_ip_id` (string) - eip

- `associate_public_ip_address` (bool) - Associate Public Ip Address

- `public_ip_band_width` (int64) - Public Ip Band Width

- `data_disks` ([]VolcengineDataDiskConfig) - Data Disks

- `instance_name` (string) - Instance Name

- `user_data` (string) - User Data

- `hpc_cluster_id` (string) - hpc

<!-- End of code generated from the comments of the VolcengineEcsConfig struct in builder/ecs/volcengine_ecs_config.go; -->


<!-- Code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `temporary_key_pair_type` (string) - `dsa` | `ecdsa` | `ed25519` | `rsa` ( the default )
  
  Specifies the type of key to create. The possible values are 'dsa',
  'ecdsa', 'ed25519', or 'rsa'.
  
  NOTE: DSA is deprecated and no longer recognized as secure, please
  consider other alternatives like RSA or ED25519.

- `temporary_key_pair_bits` (int) - Specifies the number of bits in the key to create. For RSA keys, the
  minimum size is 1024 bits and the default is 4096 bits. Generally, 3072
  bits is considered sufficient. DSA keys must be exactly 1024 bits as
  specified by FIPS 186-2. For ECDSA keys, bits determines the key length
  by selecting from one of three elliptic curve sizes: 256, 384 or 521
  bits. Attempting to use bit lengths other than these three values for
  ECDSA keys will fail. Ed25519 keys have a fixed length and bits will be
  ignored.
  
  NOTE: DSA is deprecated and no longer recognized as secure as specified
  by FIPS 186-5, please consider other alternatives like RSA or ED25519.

<!-- End of code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; -->


- `ssh_keypair_name` (string) - If specified, this is the key that will be used for SSH with the
  machine. The key must match a key pair name loaded up into the remote.
  By default, this is blank, and Packer will generate a temporary keypair
  unless [`ssh_password`](#ssh_password) is used.
  [`ssh_private_key_file`](#ssh_private_key_file) or
  [`ssh_agent_auth`](#ssh_agent_auth) must be specified when
  [`ssh_keypair_name`](#ssh_keypair_name) is utilized.


- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


- `ssh_agent_auth` (bool) - If true, the local SSH agent will be used to authenticate connections to
  the source instance. No temporary keypair will be created, and the
  values of [`ssh_password`](#ssh_password) and
  [`ssh_private_key_file`](#ssh_private_key_file) will be ignored. The
  environment variable `SSH_AUTH_SOCK` must be set for this option to work
  properly.


# Disk Devices Configuration:

<!-- Code generated from the comments of the VolcengineDataDiskConfig struct in builder/ecs/volcengine_ecs_config.go; DO NOT EDIT MANUALLY -->

- `data_disk_type` (string) - Data Disk Type

- `data_disk_size` (int32) - Data Disk Size

<!-- End of code generated from the comments of the VolcengineDataDiskConfig struct in builder/ecs/volcengine_ecs_config.go; -->


## Basic Example

Here is a basic example for Volcengine.

### Example for Use Hcl

```hcl
source "volcengine-ecs" "foo" {
    associate_public_ip_address = true
    availability_zone           = "cn-beijing-a"
    instance_type               = "ecs.g1e.large"
    source_image_id             = "image-yby53q82s721qhil65"
    ssh_clear_authorized_keys   = true
    ssh_username                = "root"
    system_disk_size            = "50"
    system_disk_type            = "ESSD_PL0"
    target_image_name           = "packer_test"
    temporary_key_pair_name     = "packer-key"
}

build {
    sources = ["source.volcengine-ecs.foo"]
    provisioner "shell" {
        inline = ["sleep 30", "yum install mysql -y"]
    }
}
```
### Example for Use Json
```json
{
  "variables": {
    "access_key": "{{ env `VOLCENGINE_ACCESS_KEY` }}",
    "secret_key": "{{ env `VOLCENGINE_SECRET_KEY` }}"
  },
  "builders": [{
    "type":"volcengine-ecs",
    "access_key":"{{user `access_key`}}",
    "secret_key":"{{user `secret_key`}}",
    "region":"cn-beijing",
    "target_image_name":"packer_xym_test",
    "source_image_id":"image-38deyjkaisf6kiyswzn9",
    "availability_zone": "cn-beijing-b",
    "instance_type":"ecs.g2i.large",
    "ssh_username":"root",
    "temporary_key_pair_name": "packer-key",
    "ssh_clear_authorized_keys": true,
    "associate_public_ip_address": true,
    "public_ip_id": "eip-13frvze8oo8hs3n6nu4m13mgo",
    "system_disk_size": "50",
    "system_disk_type": "ESSD_PL0"
  }],
  "provisioners": [{
    "type": "shell",
    "inline": [
      "sleep 30",
      "yum install mysql -y"
    ]
  }]
}
```
