package app

import (
	"github.com/beardnick/mynvim/component"
	"github.com/beardnick/mynvim/config"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
)


func ToggleSsh(nvm *nvim.Nvim){
	servers := config.Conf.Servers
	var accounts []string
	for _, server := range servers {
		accounts = append(accounts, server.Account)
	}
	data := component.CommonTreeData{
		Data:     "ssh",
		Nodes:    accounts,
		Children: nil,
	}
	tree := component.NewCommonTree(data)
	err := tree.Show(nvm)
	neovim.EchoErrStack(err)
}
