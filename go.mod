module github.com/volcengine/packer-plugin-volcengine

go 1.15

require (
	github.com/hashicorp/hcl/v2 v2.16.2
	github.com/hashicorp/packer-plugin-sdk v0.5.1
	github.com/volcengine/volcengine-go-sdk v1.0.44
	github.com/zclconf/go-cty v1.12.1
)

replace github.com/zclconf/go-cty => github.com/nywilken/go-cty v1.12.1 // added by packer-sdc fix as noted in github.com/hashicorp/packer-plugin-sdk/issues/187
