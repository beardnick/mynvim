package main

import (
	"github.com/beardnick/mynvim/component"
	"github.com/beardnick/mynvim/container"
	"github.com/beardnick/mynvim/global"
	"github.com/beardnick/mynvim/remote"
	"github.com/beardnick/mynvim/text"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		global.Nvm = p.Nvim
		p.HandleFunction(&plugin.FunctionOptions{Name: "PushBuf"}, container.PushBuf)
		p.HandleFunction(&plugin.FunctionOptions{Name: "ToggleContainer"}, container.ToggleContainer)
		p.HandleCommand(&plugin.CommandOptions{Name: "Expand", Range: "."}, text.AwkExpand)
		p.HandleCommand(&plugin.CommandOptions{Name: "PluginDir", NArgs: "+"}, remote.PluginDir)
		p.HandleCommand(&plugin.CommandOptions{Name: "Plugin", NArgs: "+"}, remote.Plugin)
		p.HandleCommand(&plugin.CommandOptions{Name: "PluginInstall", NArgs: "0"}, remote.PluginInstall)
		p.HandleCommand(&plugin.CommandOptions{Name: "Push", NArgs: "+"}, remote.Push)
		p.HandleCommand(&plugin.CommandOptions{Name: "Output"}, component.OutPut)
		p.HandleCommand(&plugin.CommandOptions{Name: "Tree"}, component.Tree)
		return nil
	})
}
