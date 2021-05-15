package component

import (
	"bytes"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
)

func Tree(nvm *nvim.Nvim, args []string) {
	lines := "\uE5FF bin\n\uE5FF global"
	out := bytes.Split([]byte(lines), []byte("\n"))
	b := nvm.NewBatch()
	var buffer nvim.Buffer
	b.Exec("setlocal nonumber",false,nil)
	b.CurrentBuffer(&buffer)
	b.SetBufferLines(buffer,0,len(out),false,out)
	b.Exec("setlocal nomodifiable",false,nil)
	err := b.Execute()
	if err != nil {
		neovim.EchoErrStack(err)
	}
}
