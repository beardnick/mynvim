package main

import (
	"mynvim/buffer"
	"mynvim/container"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func NewContainer(nvm *nvim.Nvim, args []string) error {
	c := container.Container{
		Nvm:      nvm,
		Position: "botright",
		Height:   20,
	}
	bufs, err := buffer.BufferOfType(nvm, "vim")
	if err != nil {
		return err
	}
	c.PushBufs(bufs...)
	return c.Show()
}

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		//p.HandleFunction(&plugin.FunctionOptions{Name: "Hello"}, hello)
		//p.HandleFunction(&plugin.FunctionOptions{Name: "Buffers"}, Buffers)
		p.HandleFunction(&plugin.FunctionOptions{Name: "NewContainer"}, NewContainer)
		return nil
	})
}

// ContainerList
// TabPageContainer
// ContainerLayout
// PushContainer(bufnr,eval)
// PopContainer(bufnr)
// ContainerToggle(container,position)
