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
	}
	tree := component.NewCommonTree(data)
	err := tree.Show(nvm)
	neovim.EchoErrStack(err)
}

// todo: send text to ssh job with chansend()
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
	b.Exec("40split term",false,nil)
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
	for i := 0; i < 100 ; i++ {
		time.Sleep(time.Millisecond * 100)
		line, e := nvm.CurrentLine()
		if e != nil {
			neovim.EchoErrStack(e)
			return
		}
		if strings.Contains(string(line),"#") || strings.Contains(string(line),"$"){
			e = nvm.Eval(fmt.Sprintf("chansend(%d,\"%s\")",jobid,"ls -al\n"),nil)
			neovim.EchoErrStack(e)
			return
		}
		if strings.Contains(string(line),"password:") {
			e = nvm.Eval(fmt.Sprintf("chansend(%d,\"%s\")",jobid,pass + "\n"),nil)
			neovim.EchoErrStack(e)
			return
		}
	}
}
