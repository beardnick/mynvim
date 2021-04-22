// Copyright 2016 The neovim-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package neovim

import (
	"fmt"
	"github.com/beardnick/mynvim/global"
	"github.com/pkg/errors"
)

var (
	// ErrorColor highlight error message use Identifier syntax color.
	ErrorColor = "Identifier"
	// ProgressColor highlight progress message use Identifier syntax color.
	ProgressColor = "Identifier"
	// SuccessColor highlight success message use Identifier syntax color.
	SuccessColor = "Function"
)

// Echo provide the vim 'echo' command.
func Echo(format string, a ...interface{}) error {
	global.Nvm.Command("redraw")
	return global.Nvm.Command("echo '" + fmt.Sprintf(format, a...) + "'")
}

// EchoRaw provide the raw output vim 'echo' command.
func EchoRaw(a string) error {
	global.Nvm.Command("redraw")
	return global.Nvm.Command("echo \"" + a + "\"")
}

// Echomsg provide the vim 'echomsg' command.
func Echomsg(a ...interface{}) error {
	return global.Nvm.WriteOut(fmt.Sprintln(a...))
}

// Echoerr provide the vim 'echoerr' command.
func Echoerr(format string, a ...interface{}) error {
	return global.Nvm.WritelnErr(fmt.Sprintf(format, a...))
}

func EchoErrStack(err error) error {
	return global.Nvm.WritelnErr(fmt.Sprintf("%+v", err))
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// EchohlBefore provide the vim 'echo' command with the 'echohl' highlighting prefix text.
func EchohlBefore(prefix string, highlight string, format string, a ...interface{}) error {
	global.Nvm.Command("redraw")
	suffix := "\" | echohl None | echon \""
	if prefix != "" {
		suffix += ": "
	}
	return global.Nvm.Command("echohl " + highlight + " | echo \"" + prefix + suffix + fmt.Sprintf(format, a...) + "\" | echohl None")
}

// EchohlAfter provide the vim 'echo' command with the 'echohl' highlighting message text.
func EchohlAfter(prefix string, highlight string, format string, a ...interface{}) error {
	global.Nvm.Command("redraw")
	if prefix != "" {
		prefix += ": "
	}
	return global.Nvm.Command("echo \"" + prefix + "\" | echohl " + highlight + " | echon \"" + fmt.Sprintf(format, a...) + "\" | echohl None")
}

// EchoProgress displays a command progress message to echo area.
func EchoProgress(prefix, format string, a ...interface{}) error {
	global.Nvm.Command("redraw")
	msg := fmt.Sprintf(format, a...)
	return global.Nvm.Command(fmt.Sprintf("echo \"%s: \" | echohl %s | echon \"%s ...\" | echohl None", prefix, ProgressColor, msg))
}

// EchoSuccess displays the success of the command to echo area.
func EchoSuccess(prefix string, msg string) error {
	global.Nvm.Command("redraw")
	if msg != "" {
		msg = " | " + msg
	}
	return global.Nvm.Command(fmt.Sprintf("echo \"%s: \" | echohl %s | echon 'SUCCESS' | echohl None | echon '%s' | echohl None", prefix, SuccessColor, msg))
}

// ReportError output of the accumulated errors report.
// TODO(zchee): research global.Nvm.m.ReportError behavior
// Why it does not immediately display error?
// func ReportError(v *vim.Nvim, format string, a ...interface{}) error {
// 	return global.Nvm.ReportError(fmt.Sprintf(format, a...))
// }

// ClearMsg cleanups the echo area.
func ClearMsg() error {
	return global.Nvm.Command("echon")
}
