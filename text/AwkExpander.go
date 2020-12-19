package text

import (
	"bytes"
	"fmt"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
	"os/exec"
	"regexp"
	"strings"
)

var delimiter = regexp.MustCompile(`>+`)

func AwkExpand(nvm *nvim.Nvim, ranges [2]int) {
	batch := nvm.NewBatch()
	var (
		b       nvim.Buffer
		content [][]byte
	)
	batch.CurrentBuffer(&b)
	batch.BufferLines(b, ranges[0]-1, ranges[1], false, &content)
	err := batch.Execute()
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	sb := strings.Builder{}
	for _, l := range content {
		sb.Write(l)
		sb.Write([]byte("\n"))
	}
	lines := sb.String()
	parts := delimiter.Split(lines, -1)
	if len(parts) < 2 {
		neovim.Echomsg("one or two >>> is needed to split data and cmd")
		return
	}
	var (
		data string
		opt  string
		cmd  string
	)
	if len(parts) == 2 {
		data, cmd = parts[0], parts[1]
	} else if len(parts) == 3 {
		data, opt, cmd = parts[0], parts[1], parts[2]
	} else {
		neovim.Echomsg("one or two >>> is needed to split data and cmd")
		return
	}
	r := strings.NewReplacer(
		`"`, `\"`,
		`'`, `'\\\''`,
		`{{`, `"`,
		`}}`, `"`,
	)
	data = strings.Trim(data, "\n")
	opt = strings.Trim(opt, "\n")
	cmd = strings.Trim(cmd, "\n")
	cmd = r.Replace(cmd)
	cmd = fmt.Sprintf(
		`echo '%s' | awk %s '{ print "%s" }'`, data, opt, cmd)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		neovim.Echomsg(err, string(out))
		return
	}
	err = nvm.SetBufferLines(b, ranges[0]-1, ranges[1], true, bytes.Split(out, []byte("\n")))
	if err != nil {
		neovim.Echomsg(err, string(out))
		return
	}
}
