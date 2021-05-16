package app

import (
	"fmt"
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

func SshConnect(nvm *nvim.Nvim){
	servers := config.Conf.Servers
	var account []byte
	b := nvm.NewBatch()
	account, err := nvm.CurrentLine()
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	b.Exec("40split term",false,nil)
	for _, s := range servers {
		if  s.Account != string(account) {
			continue
		}
		b.Eval(fmt.Sprintf("termopen('ssh %s -p %d')",s.Account, s.Port),nil)
	}
	err = b.Execute()
	neovim.EchoErrStack(err)
}
