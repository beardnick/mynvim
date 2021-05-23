package component

import (
	"bytes"
	"strings"

	"github.com/neovim/go-client/nvim"
)


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
	buffer, err:= LeftBar(nvm,30)
	if err != nil {
		return
	}
	b := buffer.Batch
	bufnr:= buffer.Buf
	b.Command("setlocal nonumber")
	b.Command("setlocal nowrap")
	b.SetBufferLines(bufnr, 0, len(out), false, out)
	b.SetBufferKeyMap(bufnr, "n", "q", ":<C-U>q!<CR>", map[string]bool{"noremap": true, "silent": true})
	b.SetBufferKeyMap(bufnr, "n", "<CR>", c.TreeData.NodeAction, map[string]bool{"noremap": true, "silent": true})
	b.SetBufferKeyMap(bufnr, "n", "o", c.TreeData.NodeAction, map[string]bool{"noremap": true, "silent": true})
	b.Command("setlocal nomodifiable")
	err = b.Execute()
	return
}

func NewCommonTree(data CommonTreeData) Tree {
	return CommonTree{data}
}

type CommonTreeData struct {
	Data       string
	Nodes      []string
	Children   []CommonTreeData
	NodeAction string
	TreeAction string
}

type CommonTreeNode struct {
	Content string
}

func (n CommonTreeNode) Data() string {
	return n.Content
}
