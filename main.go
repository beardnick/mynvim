package main

import (
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func Buffers(nvm *nvim.Nvim, args []string) (string, error) {
	bufs, err := nvm.Buffers()
	if err != nil {
		return "", err
	}
	results := []string{}
	for _, b := range bufs {
		name, err := nvm.BufferName(b)
		if err != nil {
			return "", err
		}
		results = append(results, name)
	}
	return strings.Join(results, ";"), nil
}

func hello(args []string) (string, error) {
	return "Hello " + strings.Join(args, " "), nil
}

func NvimConnect(addr string) (conn *nvim.Nvim, err error) {
	return nvim.Dial(addr, []nvim.DialOption{}...)
}

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleFunction(&plugin.FunctionOptions{Name: "Hello"}, hello)
		p.HandleFunction(&plugin.FunctionOptions{Name: "Buffers"}, Buffers)
		return nil
	})
}
