package app

import (
	"fmt"
	"github.com/beardnick/mynvim/component"
	"github.com/beardnick/mynvim/config"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
	"strings"
	"time"
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
		NodeAction: ":<C-U>SshConnect<CR>",
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
	jobid := 0

	buffer, err := nvm.CreateBuffer(false, false)
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	b.Command("split")
	b.Command("wincmd J")
	b.SetCurrentBuffer(buffer)

	var win nvim.Window
	b.CurrentWindow(&win)
	err = b.Execute()
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	b.SetWindowHeight(win, 30)
	pass := ""
	for _, s := range servers {
		if  s.Account != string(account) {
			continue
		}
		pass = s.Password
		b.Eval(fmt.Sprintf("termopen('ssh %s -p %d')",s.Account, s.Port),&jobid)
	}
	err = b.Execute()
	// try to send password
	var line []byte
	for i := 0; i < 100 ; i++ {
		time.Sleep(time.Millisecond * 100)
		line, err = nvm.CurrentLine()
		if err != nil {
			break
		}
		if strings.Contains(string(line),"#") || strings.Contains(string(line),"$"){
			break
		}
		if strings.Contains(string(line),"password:") {
			err = nvm.Eval(fmt.Sprintf("chansend(%d,\"%s\")",jobid,pass + "\n"),nil)
			break
		}
	}
	if err !=nil {
		neovim.EchoErrStack(err)
		return
	}
	err = nvm.Command("normal I")
	neovim.EchoErrStack(err)
}
