package commands

import (
	"github.com/charmbracelet/log"
	"github.com/pernydev/plugstep/pkg/plugstep"
	"github.com/pernydev/plugstep/pkg/plugstep/install/plugins"
	"github.com/pernydev/plugstep/pkg/plugstep/install/server"
)

func InstallCommand(ps *plugstep.Plugstep) {
	log.Debug("Installing server JAR and all plugins...", "serverjar", ps.Config.Server.Project, "minecraft-version", ps.Config.Server.MinecraftVersion, "plugins", len(ps.Config.Plugins))
	server.InstallServer(ps)
	plugins.InstallPlugins(ps)
}
