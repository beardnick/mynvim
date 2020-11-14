package main

import (
	"github.com/beardnick/mynvim/container"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		//p.HandleFunction(&plugin.FunctionOptions{Name: "Hello"}, hello)
		//p.HandleFunction(&plugin.FunctionOptions{Name: "Buffers"}, Buffers)
		p.HandleFunction(&plugin.FunctionOptions{Name: "PushBuf"}, container.PushBuf)
		p.HandleFunction(&plugin.FunctionOptions{Name: "ToggleContainer"}, container.ToggleContainer)
		return nil
	})
}

// ContainerList
// TabPageContainer
// ContainerLayout
// PushContainer(bufnr,eval)
// PopContainer(bufnr)
// ContainerToggle(container,position)
