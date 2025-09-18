package plugins

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/pernydev/plugstep/pkg/plugstep/config"
)

type PluginSource interface {
	GetPluginDownload(c config.PluginConfig) (*PluginDownload, error)
}

type ChecksumType string

const (
	ChecksumTypeSha256 ChecksumType = "sha256"
	ChecksumTypeSha512 ChecksumType = "sha512"
)

type PluginDownload struct {
	URL          string
	Checksum     string
	ChecksumType ChecksumType
}

type ModrinthPluginSource struct {
	apiURL string
}

type ModrinthVersion struct {
	VersionNumber string         `json:"version_number"`
	Files         []ModrinthFile `json:"files"`
}

type ModrinthFile struct {
	Hashes struct {
		Sha512 string `json:"sha512"`
	} `json:"hashes"`
	URL     string `json:"url"`
	Primary bool   `json:"primary"`
}

func (m *ModrinthPluginSource) GetPluginDownload(c config.PluginConfig) (*PluginDownload, error) {
	url := fmt.Sprintf("%s/project/%s/version", m.apiURL, *c.Resource)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("got %d", r.StatusCode)
	}

	var response []ModrinthVersion
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	version := findModrinthVersion(response, *c.Version)
	if version == nil {
		log.Error("Plugin version not found!", "source", c.Source, "plugin", c.Resource, "version", c.Version)
	}

	file := findModrinthPrimaryFile(version.Files)
	if version == nil {
		log.Error("Plugin version has no primary file!", "source", c.Source, "plugin", c.Resource, "version", c.Version)
	}

	return &PluginDownload{
		URL:          file.URL,
		Checksum:     file.Hashes.Sha512,
		ChecksumType: ChecksumTypeSha512,
	}, nil
}

func findModrinthVersion(response []ModrinthVersion, version string) *ModrinthVersion {
	for _, resp := range response {
		if resp.VersionNumber == version {
			return &resp
		}
	}
	return nil
}

func findModrinthPrimaryFile(files []ModrinthFile) *ModrinthFile {
	for _, f := range files {
		if f.Primary == true {
			return &f
		}
	}
	return nil
}

type CustomPluginSource struct{}

func (source *CustomPluginSource) GetPluginDownload(c config.PluginConfig) (*PluginDownload, error) {
	return &PluginDownload{
		URL:          *c.DownloadURL,
		Checksum:     "nocheck",
		ChecksumType: ChecksumTypeSha256,
	}, nil
}

func GetSource(source config.PluginSource) PluginSource {
	switch source {
	case config.PluginSourceModrinth:
		return &ModrinthPluginSource{
			apiURL: "https://api.modrinth.com/v2",
		}
	case config.PluginSourceCustom:
		return &CustomPluginSource{}
	}
	return nil
}
