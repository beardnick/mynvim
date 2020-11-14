package nvimutil

import (
	"fmt"
	"github.com/neovim/go-client/nvim"
	"github.com/zchee/nvim-go/pkg/nvimutil"
	"strconv"
)

func Bufwinnr(nvm *nvim.Nvim, buf int) (win int, err error) {
	result, err := nvm.Exec(fmt.Sprintf("echo bufwinnr(%d)", buf), true)
	if err != nil {
		nvimutil.Echomsg(nvm, err)
	}
	win, err = strconv.Atoi(result)
	return
}
