package plugins

import "github.com/pernydev/plugstep/pkg/plugstep/config"

type CustomPluginSource struct{}

func (source *CustomPluginSource) GetPluginDownload(c config.PluginConfig) (*PluginDownload, error) {
	return &PluginDownload{
		URL:          *c.DownloadURL,
		Checksum:     "nocheck",
		ChecksumType: ChecksumTypeSha256,
	}, nil
}
