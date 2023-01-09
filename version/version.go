package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	Version           = "0.0.1"
	VersionPrerelease = "dev"
	PluginVersion     = version.InitializePluginVersion(Version, VersionPrerelease)
)
