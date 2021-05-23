package app

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/beardnick/mynvim/component"
	"github.com/beardnick/mynvim/config"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
)

func ToggleSSH(nvm *nvim.Nvim) {
	servers := config.Conf.Servers
	accounts := make([]string, 0, len(servers))
	for _, server := range servers {
		accounts = append(accounts, server.Account)
	}
	data := component.CommonTreeData{
		Data:       "ssh",
		Nodes:      accounts,
		Children:   nil,
		NodeAction: ":<C-U>SshConnect<CR>",
	}
	tree := component.NewCommonTree(data)
	err := tree.Show(nvm)
	neovim.EchoErrStack(err)
}

func SshConnect(nvm *nvim.Nvim) {
	servers := config.Conf.Servers
	var account []byte
	b := nvm.NewBatch()
	account, err := nvm.CurrentLine()
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	jobid := 0

	// todo 使用wincontainer来统一布局
	_,err = component.BottomBar(nvm,20)
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	pass := ""
	for _, s := range servers {
		if s.Account != string(account) {
			continue
		}
		pass = s.Password
		b.Eval(fmt.Sprintf("termopen('ssh %s -p %d')", s.Account, s.Port), &jobid)
	}
	err = b.Execute()
	// try to send password
	err = trySendPassword(nvm, jobid, pass)
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	err = nvm.Command("normal I")
	neovim.EchoErrStack(err)
}

func trySendPassword(nvm *nvim.Nvim, jobid int, passwd string) (err error) {
	var line [][]byte
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		err = nvm.Eval("getbufline(bufnr('%'),1,'$')", &line)
		if err != nil {
			return
		}
		content := bytes.Join(line, []byte{})
		if strings.Contains(string(content), "#") || strings.Contains(string(content), "$") {
			return
		}
		if strings.Contains(string(content), "yes") {
			err = nvm.Eval(fmt.Sprintf("chansend(%d,\"%s\")", jobid, "yes\n"), nil)
			return
		}
		if strings.Contains(string(content), "password:") {
			err = nvm.Eval(fmt.Sprintf("chansend(%d,\"%s\")", jobid, passwd+"\n"), nil)
			return
		}
	}
	return
}
