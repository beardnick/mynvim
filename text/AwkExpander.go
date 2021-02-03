package text

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/beardnick/mynvim/neovim"
	"github.com/benhoyt/goawk/interp"
	"github.com/benhoyt/goawk/parser"
	"github.com/neovim/go-client/nvim"
	"regexp"
	"strings"
)

var delimiter = regexp.MustCompile(`>>>+`)

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
	src := fmt.Sprintf(`{print "%s"}`, cmd)
	prog, err := parser.ParseProgram([]byte(src), nil)
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	vars, err := ParseAwkOpt(opt)
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	out := bytes.NewBuffer(nil)
	config := &interp.Config{
		Stdin:  bytes.NewReader([]byte(data)),
		Vars:   vars,
		Output: out,
	}
	_, err = interp.ExecProgram(prog, config)
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	err = nvm.SetBufferLines(b, ranges[0]-1, ranges[1], true, bytes.Split(out.Bytes(), []byte("\n")))
	if err != nil {
		neovim.Echomsg(err)
		return
	}
}

func ParseAwkOpt(opt string) (vars []string, err error) {
	fieldSep := " "
	vs := make([]string, 0)
	args := strings.Fields(opt)
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-F":
			if i+1 >= len(args) {
				err = errors.New("flag needs an argument: -F")
				return
			}
			i++
			fieldSep = args[i]
		case "-v":
			if i+1 >= len(args) {
				err = errors.New("flag needs an argument: -v")
				return
			}
			i++
			vs = append(vs, args[i])
		}
	}
	vars = []string{"FS", fieldSep}
	for _, v := range vs {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			err = fmt.Errorf("-v flag must be in format name=value")
			return
		}
		vars = append(vars, parts[0], parts[1])
	}
	return
}
