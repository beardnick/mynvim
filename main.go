package main

import (
	"github.com/beardnick/mynvim/container"
	"github.com/beardnick/mynvim/global"
	"github.com/beardnick/mynvim/text"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		global.Nvm = p.Nvim
		p.HandleFunction(&plugin.FunctionOptions{Name: "PushBuf"}, container.PushBuf)
		p.HandleFunction(&plugin.FunctionOptions{Name: "ToggleContainer"}, container.ToggleContainer)
		p.HandleCommand(&plugin.CommandOptions{Name: "Expand", Range: "."}, text.AwkExpand)
		return nil
	})
}

// ContainerList
// TabPageContainer
// ContainerLayout
// PushContainer(bufnr,eval)
// PopContainer(bufnr)
// ContainerToggle(container,position)
