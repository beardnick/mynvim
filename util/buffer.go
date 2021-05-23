package util

import "github.com/neovim/go-client/nvim"

type Buffer struct {
	Buf nvim.Buffer
	Win nvim.Window
	Batch *nvim.Batch
	nvm *nvim.Nvim
}
