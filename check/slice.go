package check

import "github.com/neovim/go-client/nvim"

func ContainsBuffer(src []nvim.Buffer, buf nvim.Buffer) bool {
	for _, b := range src {
		if b == buf {
			return true
		}
	}
	return false
}
