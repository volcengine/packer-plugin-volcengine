# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Volcengine"
  description = "TODO"
  identifier = "packer/volcengine/volcengine"
  component {
    type = "builder"
    name = "Volcengine Image Builder"
    slug = "volcengine-ecs"
  }
}
