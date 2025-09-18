package config

type PlugstepConfig struct {
	Server  ServerConfig   `toml:"server"`
	Plugins []PluginConfig `toml:"plugins"`
}

type ServerJarVendor string

const (
	ServerJarVendorPaperMC ServerJarVendor = "papermc"
	// TODO: Add more
)

type ServerConfig struct {
	Vendor           ServerJarVendor `toml:"vendor"`
	Project          string          `toml:"project"`
	MinecraftVersion string          `toml:"minecraft_version"`
	Version          string          `toml:"version"`
}

type PluginSource string

const (
	PluginSourcePaperHangar PluginSource = "paper-hangar"
	PluginSourcePolymart    PluginSource = "polymart"
	PluginSourceModrinth    PluginSource = "modrinth"
	PluginSourceCustom      PluginSource = "custom"
)

type PluginConfig struct {
	Source      PluginSource `toml:"source"`
	Resource    *string      `toml:"resource"`
	Version     *string      `toml:"version"`
	DownloadURL *string      `toml:"download_url"`
}
