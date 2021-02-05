package buffer

import (
	"github.com/neovim/go-client/nvim"
	"strings"
)

func BufferOfType(nvm *nvim.Nvim, filetype string) ([]nvim.Buffer, error) {
	allBufs, err := nvm.Buffers()
	if err != nil {
		return []nvim.Buffer{}, err
	}
	bufs := []nvim.Buffer{}
	for _, b := range allBufs {
		realType, err := BufferFileType(nvm, b)
		if err != nil {
			return []nvim.Buffer{}, err
		}
		if realType != filetype {
			continue
		}
		bufs = append(bufs, b)
	}
	return bufs, nil
}

func BufferFileType(nvm *nvim.Nvim, buffer nvim.Buffer) (filetype string, err error) {
	err = nvm.BufferOption(buffer, "filetype", &filetype)
	return
}

func CurrentBufferLines(nvm *nvim.Nvim, ranges [2]int) (lines string, err error) {
	batch := nvm.NewBatch()
	var (
		b       nvim.Buffer
		content [][]byte
	)
	batch.CurrentBuffer(&b)
	batch.BufferLines(b, ranges[0]-1, ranges[1], false, &content)
	err = batch.Execute()
	if err != nil {
		return
	}
	sb := strings.Builder{}
	for _, l := range content {
		sb.Write(l)
		sb.Write([]byte("\n"))
	}
	lines = sb.String()
	return
}
