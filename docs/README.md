The Volcengine plugin is intended as a starting point for creating Packer plugins, containing:

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    gridscale = {
      version = ">= 1.0.0"
      source  = "github.com/volcengine/ecs"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/volcengine/ecs
```

### Components

#### Builders

- [volcengine-ecs builder](/packer/integrations/volcengine/latest/components/builder/ecs) - provides the capability to build customized images based
  on an existing base image.
