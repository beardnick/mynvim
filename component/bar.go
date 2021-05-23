package component

import (
	"fmt"

	"github.com/beardnick/mynvim/util"
	"github.com/neovim/go-client/nvim"
)

type Position int

const (
	LeftMost Position = iota
	RightMost
	Bottom
	Top
)

func LeftBar(nvm *nvim.Nvim, size int) (buffer util.Buffer, err error) {
	return newBar(nvm, size, LeftMost)
}

func RightBar(nvm *nvim.Nvim, size int) (buffer util.Buffer, err error) {
	return newBar(nvm, size, RightMost)
}


func BottomBar(nvm *nvim.Nvim, size int) (buffer util.Buffer, err error) {
	return newBar(nvm, size, Bottom)
}

func TopBar(nvm *nvim.Nvim, size int) (buffer util.Buffer, err error) {
	return newBar(nvm, size, Top)
}


func newBar(nvm *nvim.Nvim, size int, pos Position) (buffer util.Buffer, err error) {
	bufnr, err := nvm.CreateBuffer(false, false)
	if err != nil {
		return
	}
	buffer.Buf = bufnr
	winpos := ""
	b := nvm.NewBatch()
	buffer.Batch = b
	switch pos {
	case LeftMost:
		winpos = "wincmd H"
	case RightMost:
		winpos = "wincmd L"
	case Bottom:
		winpos = "wincmd J"
	case Top:
		winpos = "wincmd K"
	}
	b.Command("vsplit")
	b.Command(winpos)
	b.SetCurrentBuffer(bufnr)
	b.Command(fmt.Sprintf("autocmd QuitPre * %dbwipeout!", bufnr))
	b.CurrentWindow(&buffer.Win)
	err = b.Execute()
	if err != nil {
		return
	}
	switch pos {
	case LeftMost, RightMost:
		err = nvm.SetWindowWidth(buffer.Win, size)
	case Bottom, Top:
		err = nvm.SetWindowHeight(buffer.Win, size)
	}
	return
}
