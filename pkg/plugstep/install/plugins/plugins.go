package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/pernydev/plugstep/pkg/plugstep"
	"github.com/pernydev/plugstep/pkg/plugstep/config"
	"github.com/pernydev/plugstep/pkg/plugstep/utils"
)

func InstallPlugins(ps *plugstep.Plugstep) error {
	log.Info("Starting plugin download", "plugins", len(ps.Config.Plugins))

	var err error
	var status PluginInstallStatus
	var errorPlugin string

	fmt.Println("")

	installed := 0
	checked := 0

	for _, plugin := range ps.Config.Plugins {
		status, err = installPlugin(ps, &plugin)
		if err != nil {
			errorPlugin = *plugin.Resource
			break
		}
		fmt.Printf("\r")
		if status == PluginInstallStatusInstalled {
			log.Infof("Installed %s", *plugin.Resource)
			installed++
		} else {
			log.Infof("Checked %s", *plugin.Resource)
			checked++
		}
	}
	fmt.Printf("\r")
	fmt.Println("")

	if err != nil {
		log.Error("Plugin installation failed.", "plugin", errorPlugin, "err", err)
		return err
	}

	removed := removeOld(ps)

	log.Info("Plugins ready.", "installed", installed, "checked", checked, "removed", removed)

	return nil
}

// TODO: Add error handling
func removeOld(ps *plugstep.Plugstep) int {
	entries, err := os.ReadDir(filepath.Join(ps.ServerDirectory, "plugins"))
	if err != nil {
		log.Error("Error reading directory", "err", err)
		return 0
	}

	removed := 0
	for _, f := range entries {
		if f.IsDir() {
			continue
		}
		found := false
		for _, p := range ps.Config.Plugins {
			if f.Name() == *p.Resource+".jar" {
				found = true
				continue
			}
		}
		if found == true {
			continue
		}

		os.Remove(filepath.Join(ps.ServerDirectory, "plugins", f.Name()))
		log.Infof("Removed %s", strings.Split(f.Name(), ".")[0])
		removed++
	}

	return removed
}

type PluginInstallStatus string

const (
	PluginInstallStatusInstalled PluginInstallStatus = "installed"
	PluginInstallStatusChecked   PluginInstallStatus = "checked"
)

func installPlugin(ps *plugstep.Plugstep, p *config.PluginConfig) (PluginInstallStatus, error) {
	badge := lipgloss.NewStyle().
		Background(lipgloss.Color("#ca9ee6")).
		Foreground(lipgloss.Color("#232634")).
		PaddingRight(2).
		PaddingLeft(2).
		Transform(strings.ToUpper).
		Render("INSTALLING")

	sourceBadge := lipgloss.NewStyle().
		Background(lipgloss.Color("#e78284")).
		Foreground(lipgloss.Color("#232634")).
		PaddingRight(2).
		PaddingLeft(2).
		Transform(strings.ToUpper).
		Render(string(p.Source))

	fmt.Printf("\r%s%s %s", badge, sourceBadge, *p.Resource)
	source := GetSource(p.Source)
	if source == nil {
		return "", fmt.Errorf("invalid source")
	}

	download, err := source.GetPluginDownload(*p)
	if err != nil {
		return "", err
	}

	file := filepath.Join(ps.ServerDirectory, "plugins", *p.Resource+".jar")

	hash := ""
	switch download.ChecksumType {
	case ChecksumTypeSha256:
		hash, _ = utils.CalculateFileSHA256(file)
	case ChecksumTypeSha512:
		hash, _ = utils.CalculateFileSHA512(file)
	}

	if hash == download.Checksum {
		return PluginInstallStatusChecked, nil
	}

	utils.DownloadFile(download.URL, file)

	return PluginInstallStatusInstalled, nil
}
