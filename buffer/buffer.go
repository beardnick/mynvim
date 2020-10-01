package buffer

import "github.com/neovim/go-client/nvim"

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
