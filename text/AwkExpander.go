package text

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/beardnick/mynvim/buffer"
	"github.com/beardnick/mynvim/neovim"
	"github.com/benhoyt/goawk/interp"
	"github.com/benhoyt/goawk/parser"
	"github.com/neovim/go-client/nvim"
	"regexp"
	"strings"
)

var delimiter = regexp.MustCompile(`>>>+`)

func AwkExpand(nvm *nvim.Nvim, ranges [2]int) {
	lines, err := buffer.CurrentBufferLines(nvm, ranges)
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	result, err := Expend(lines)
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	b, err := nvm.CurrentBuffer()
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
	out := []byte(result)
	err = nvm.SetBufferLines(b, ranges[0]-1, ranges[1], true, bytes.Split(out, []byte("\n")))
	if err != nil {
		neovim.EchoErrStack(err)
		return
	}
}

func Expend(lines string) (expended string, err error) {
	data, opt, template, err := parseTemplateArg(lines, delimiter)
	if err != nil {
		return
	}
	template = parseTemplate(template)
	src := fmt.Sprintf(`{print "%s"}`, template)
	prog, err := parser.ParseProgram([]byte(src), nil)
	if err != nil {
		return
	}

	vars, err := parseAwkOpt(opt)
	if err != nil {
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
		return
	}
	expended = string(out.Bytes())
	return
}

func parseTemplateArg(lines string, delimiter *regexp.Regexp) (data, opt, template string, err error) {
	parts := delimiter.Split(lines, 3)
	if len(parts) < 2 {
		err = errors.New(
			fmt.Sprintf("one or two %v match is needed to split data and template", delimiter))
		return
	}
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if len(parts) == 2 {
		data, template = parts[0], parts[1]
	}
	if len(parts) == 3 {
		data, opt, template = parts[0], parts[1], parts[2]
	}
	return
}

var arch = regexp.MustCompile("{{.*?}}")

func parseTemplate(t string) string {
	archs := arch.FindAllString(t, -1)
	r := strings.NewReplacer(
		`"`, `\"`,
		`'`, `'\\\''`,
		"\n", "\"\nprint \"",
	)
	t = r.Replace(t)
	r = strings.NewReplacer(
		`{{`, `"`,
		`}}`, `"`,
	)
	archIndex := arch.FindAllStringIndex(t, -1)
	for i := 0; i < len(archs); i++ {
		s := r.Replace(archs[i])
		b, e := archIndex[i][0], archIndex[i][1]
		t = t[:b] + s + t[e:]
	}
	return t
}

func parseAwkOpt(opt string) (vars []string, err error) {
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
			err = errors.New("-v flag must be in format name=value")
			return
		}
		vars = append(vars, parts[0], parts[1])
	}
	return
}
