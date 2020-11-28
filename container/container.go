package container

import (
	"fmt"
	"github.com/beardnick/mynvim/check"
	"github.com/beardnick/mynvim/neovim"
	"github.com/neovim/go-client/nvim"
	"strconv"
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
	b := c.Nvm.NewBatch()
	b.Command(fmt.Sprintf("%s %d split", c.Position, c.Height))
	b.Command(eval)
	err = b.Execute()
	if err != nil {
		return
	}
	buffer, err := c.Nvm.CurrentBuffer()
	neovim.Echomsg(buffer)
	if err != nil {
		return
	}
	c.bufs = append(c.bufs, buffer)
	container_exists = true
	return
}

func (c *Container) PushBufEval(eval string) (err error) {
	neovim.Echomsg(eval)
	neovim.Echomsg(c.bufs)
	if len(c.bufs) == 0 {
		return c.CreateEval(eval)
	}
	if !container_exists {
		err = c.Show()
		if err != nil {
			return
		}
	}
	win, err := neovim.Bufwinnr(c.Nvm, int(c.bufs[len(c.bufs)-1]))
	if err != nil {
		return
	}
	b := c.Nvm.NewBatch()
	b.Command(fmt.Sprintf("%dwincmd w", win))
	b.Command("wincmd v")
	b.Command("wincmd l")
	b.Command(eval)
	err = b.Execute()
	if err != nil {
		return err
	}
	buffer, err := c.Nvm.CurrentBuffer()
	if err != nil {
		return
	}
	c.bufs = append(c.bufs, buffer)
	return
}

func (c *Container) PushBuf(buffer nvim.Buffer) (err error) {
	if check.ContainsBuffer(c.bufs, buffer) {
		return
	}
	win, err := neovim.Bufwinnr(c.Nvm, int(buffer))
	if err != nil {
		return err
	}
	err = c.Nvm.Command(fmt.Sprintf("%vhide", win))
	if err != nil {
		return err
	}
	return c.PushBufEval(fmt.Sprintf("%dbuffer", int(buffer)))
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
		win, err := neovim.Bufwinnr(c.Nvm, int(buf))
		if err != nil {
			return err
		}
		err = c.Nvm.Command(fmt.Sprintf("%vhide", win))
		if err != nil {
			return err
		}
	}
	container_exists = false
	return
}

func (c *Container) Show() (err error) {
	_, err = c.Nvm.Exec(
		fmt.Sprintf("%s %d split", c.Position, c.Height),
		false,
	)
	if err != nil {
		return
	}
	length := len(c.bufs)
	err = c.Nvm.SetCurrentBuffer(c.bufs[0])
	if err != nil {
		return
	}
	b := c.Nvm.NewBatch()
	for i := 1; i < length; i++ {
		b.Command("wincmd v")
		b.Command("wincmd l")
		b.SetCurrentBuffer(c.bufs[i])
		err = b.Execute()
		if err != nil {
			return
		}
	}
	container_exists = true
	return
}
func PushBuf(nvm *nvim.Nvim, args []string) (err error) {
	if !container_exists {
		global_container.Nvm = nvm
	}
	for _, arg := range args {
		if buf, err := strconv.Atoi(arg); err == nil {
			global_container.PushBuf(nvim.Buffer(buf))
			continue
		}
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
		neovim.Echomsg("win container empty")
		return
	}
	if err != nil {
		return err
	}
	err = global_container.Show()
	return
}
