package neovim

import (
	"fmt"
	"github.com/neovim/go-client/nvim"
	"strconv"
)

func Bufwinnr(nvm *nvim.Nvim, buf int) (win int, err error) {
	result, err := nvm.Exec(fmt.Sprintf("echo bufwinnr(%d)", buf), true)
	if err != nil {
		Echomsg(err)
	}
	win, err = strconv.Atoi(result)
	return
}
