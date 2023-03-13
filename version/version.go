package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	Version           = "1.0.0"
	VersionPrerelease = "dev"
	PluginVersion     = version.InitializePluginVersion(Version, VersionPrerelease)
)
