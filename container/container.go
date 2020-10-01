package container

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type Container struct {
	Nvm      *nvim.Nvim
	bufs     []nvim.Buffer
	Position string
	Height   int
}

func (c *Container) PushBufs(buffer ...nvim.Buffer) {
	c.bufs = append(c.bufs, buffer...)
}

func (c *Container) PushBuf(buffer nvim.Buffer) {
	c.bufs = append(c.bufs, buffer)
}

func (c *Container) PopBuf() {
	length := len(c.bufs)
	if length == 0 {
		return
	}
	c.bufs = c.bufs[:length-1]
}

func (c *Container) RemoveBuf(buffer nvim.Buffer) {
	length := len(c.bufs)
	if length == 0 {
		return
	}
	target := -1
	for i := 0; i < length; i++ {
		if c.bufs[i] == buffer {
			target = i
			break
		}
	}
	if target == -1 {
		return
	}
	c.bufs = append(c.bufs[:target], c.bufs[target+1:]...)
}

func (c *Container) Show() error {
	_, err := c.Nvm.Exec(
		fmt.Sprintf("%s %d split", c.Position, c.Height),
		false,
	)
	if err != nil {
		return err
	}
	length := len(c.bufs)
	c.Nvm.SetCurrentBuffer(c.bufs[0])
	for i := 1; i < length; i++ {
		_, err := c.Nvm.Exec(
			fmt.Sprintf("wincmd v"),
			false,
		)
		if err != nil {
			return err
		}
		_, err = c.Nvm.Exec(
			fmt.Sprintf("wincmd l"),
			false,
		)
		if err != nil {
			return err
		}
		err = c.Nvm.SetCurrentBuffer(c.bufs[i])
		if err != nil {
			return err
		}
	}
	return nil
}
