package component

import (
	"fmt"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
)

func OutPut(nvm *nvim.Nvim, args []string) {
	for i := 0; i < 1000; i++ {
		err := nvm.Eval(fmt.Sprintf(`append('$',"test%d")`, i), nil)
		if err != nil {
			neovim.EchoErrStack(err)
		}
		_, err = nvm.Exec(fmt.Sprintf(`normal j`), false)
		if err != nil {
			neovim.EchoErrStack(err)
		}
	}
}
