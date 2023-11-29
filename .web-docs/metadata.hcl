# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Volcengine"
  description = "A multi-component plugin to create custom images."
  identifier = "packer/volcengine/volcengine"
  component {
    type = "builder"
    name = "Volcengine Image Builder"
    slug = "ecs"
  }
}
