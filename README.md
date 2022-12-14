# Packer Builder for Volcengine ECS

This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Volcengine ECS](https://www.volcengine.com/product/ecs) image.

## Requirements
* [Go 1.15+](https://golang.org/doc/install)
* [Packer](https://www.packer.io/intro/getting-started/install.html)

## Build & Installation

### Install from source:

Clone repository to `$GOPATH/src/github.com/volcengine/packer-plugin-volcengine`

```sh
$ mkdir -p $GOPATH/src/github.com/volcengine; 
$ cd $GOPATH/src/github.com/volcengine
$ git clone git@github.com:volcengine/packer-plugin-volcengine.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/volcengine/packer-plugin-volcengine
$ make install
```

### Install from HCL:
```hcl
packer {
  required_plugins {
    volcengine = {
      version = ">= 0.0.1"
      source  = "github.com/volcengine/volcengine"
    }
  }
}
```


### Install from release:

* Download binaries from the [releases page](https://github.com/volcengine/packer-plugin-volcengine/releases).
* [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin, or simply put it into the same directory with JSON templates.
* Move the downloaded binary to `~/.packer.d/plugins/`

## Usage for ECS
Here is a sample template, which you can also find in the `example/` directory
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
Enter the API user credentials in your terminal with the following commands. Replace the <AK> and <SK> with your user details.
```sh
export VOLCENGINE_ACCESS_KEY=<AK>
export VOLCENGINE_SECRET_KEY=<SK>
```
Then run Packer using the example template with the command underneath.
```
# use for ECS
packer build example/volcengine.json
```


