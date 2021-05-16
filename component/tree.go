package component

import (
	"bytes"
	"github.com/neovim/go-client/nvim"
	"strings"
)

//func Tree(nvm *nvim.Nvim, args []string) {
//	lines := "\uE5FF bin\n\uE5FF global"
//	out := bytes.Split([]byte(lines), []byte("\n"))
//	b := nvm.NewBatch()
//	var buffer nvim.Buffer
//	b.Exec("setlocal nonumber",false,nil)
//	b.CurrentBuffer(&buffer)
//	b.SetBufferLines(buffer,0,len(out),false,out)
//	b.Exec("setlocal nomodifiable",false,nil)
//	err := b.Execute()
//	if err != nil {
//		neovim.EchoErrStack(err)
//	}
//}

type TreeNode interface {
	Data() string
}

type Tree interface {
	Show(nvm *nvim.Nvim) (err error)
}

type CommonTree struct {
	//ChildTree []CommonTree
	//ChildNode []CommonTreeNode
	//Parent *CommonTree
	TreeData CommonTreeData
}

func (c CommonTree) Show(nvm *nvim.Nvim) (err error) {
	data := strings.Join(c.TreeData.Nodes, "\n")
	out := bytes.Split([]byte(data), []byte("\n"))
	b := nvm.NewBatch()
	var buffer nvim.Buffer
	b.Exec("30vsplit tree",false,nil)
	b.Exec("setlocal nonumber",false,nil)
	b.CurrentBuffer(&buffer)
	b.SetBufferLines(buffer,0,len(out),false,out)
	b.SetBufferKeyMap(buffer,"n","q","<C-U>:q!<CR>", map[string]bool{"noremap":true,"silent":true})
	b.SetBufferKeyMap(buffer,"n","<CR>","<C-U>:SshConnect<CR>", map[string]bool{"noremap":true,"silent":true})
	b.Exec("setlocal nomodifiable",false,nil)
	err = b.Execute()
	return
}

func NewCommonTree(data CommonTreeData) Tree {
	return CommonTree{data}
}

type CommonTreeData struct {
	Data     string
	Nodes    []string
	Children []CommonTreeData
}

type CommonTreeNode struct {
	Content string
}

func (n CommonTreeNode) Data() string {
	return n.Content
}
