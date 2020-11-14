package container

import (
	"fmt"
	nvimutil "github.com/beardnick/mynvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

var (
	container_exists           = false
	global_container Container = Container{
		bufs:     nil,
		Position: "botright",
		Height:   20,
	}
)

type Container struct {
	Nvm      *nvim.Nvim
	bufs     []nvim.Buffer
	Position string
	Height   int
}

func (c *Container) PushBufs(buffer ...nvim.Buffer) {
	for _, buf := range buffer {
		c.PushBuf(buf)
	}
}

func (c *Container) CreateEval(eval string) (err error) {
	_, err = c.Nvm.Exec(
		fmt.Sprintf("%s %d split", c.Position, c.Height),
		false,
	)
	if err != nil {
		return err
	}
	_, err = c.Nvm.Exec(eval, false)
	if err != nil {
		return
	}
	buffer, err := c.Nvm.CurrentBuffer()
	nvimutil.Echomsg(c.Nvm, buffer)
	if err != nil {
		return
	}
	c.bufs = append(c.bufs, buffer)
	return
}

func (c *Container) PushBufEval(eval string) (err error) {
	if len(c.bufs) == 0 {
		return c.CreateEval(eval)
	}
	if !container_exists && len(c.bufs) > 0 {
		err = c.Show()
		if err != nil {
			return
		}
	}
	if !container_exists && len(c.bufs) == 0 {
		_, err := c.Nvm.Exec(
			fmt.Sprintf("%s %d split", c.Position, c.Height),
			false,
		)
		if err != nil {
			return err
		}
	}
	nvimutil.Echomsg(c.Nvm, eval)
	if len(c.bufs) == 0 {
	}
	err = c.Nvm.SetCurrentBuffer(c.bufs[len(c.bufs)-1])
	if err != nil {
		return
	}
	_, err = c.Nvm.Exec(
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
	_, err = c.Nvm.Exec(eval, false)
	if err != nil {
		return
	}
	buffer, err := c.Nvm.CurrentBuffer()
	nvimutil.Echomsg(c.Nvm, buffer)
	if err != nil {
		return
	}
	c.bufs = append(c.bufs, buffer)
	return
}

func (c *Container) PushBuf(buffer nvim.Buffer) (err error) {
	if container_exists {
		err = c.Show()
		if err != nil {
			return
		}
	}
	err = c.Nvm.SetCurrentBuffer(c.bufs[len(c.bufs)-1])
	if err != nil {
		return
	}
	_, err = c.Nvm.Exec(
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
	err = c.Nvm.SetCurrentBuffer(buffer)
	if err != nil {
		return err
	}
	c.bufs = append(c.bufs, buffer)
	return
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

func (c *Container) Hide() (err error) {
	for _, buf := range c.bufs {
		win, err := nvimutil.Bufwinnr(c.Nvm, int(buf))
		if err != nil {
			return err
		}
		_, err = c.Nvm.Exec(fmt.Sprintf("%vhide", win), false)
		if err != nil {
			return err
		}
	}
	container_exists = false
	return
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
	container_exists = true
	return nil
}
func PushBuf(nvm *nvim.Nvim, args []string) (err error) {
	if !container_exists {
		global_container.Nvm = nvm
	}
	for _, arg := range args {
		//nvimutil.Echomsg(nvm,arg)
		err = global_container.PushBufEval(arg)
		if err != nil {
			return
		}
	}
	return
}

func ToggleContainer(nvm *nvim.Nvim, args []string) (err error) {
	if container_exists {
		err = global_container.Hide()
		return
	}
	if len(global_container.bufs) == 0 {
		nvimutil.Echomsg(nvm, "win container empty")
		return
	}
	if err != nil {
		return err
	}
	err = global_container.Show()
	return
}
