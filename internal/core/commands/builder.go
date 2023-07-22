package commands

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
)

type CMD struct {
	Binary   string
	Commands []string
	Error    error
}

type CMDNewArgs struct {
	Binary      string
	CommandArgs string
}

type CMDBuilder struct {
	binary      string
	commandArgs string
	commands    []string
	error       error
}

func (b *CMDBuilder) WithBinary(binary string) *CMDBuilder {
	b.binary = binary
	return b
}

func (b *CMDBuilder) WithCommands(commandArgs string) *CMDBuilder {
	commands, err := utils.GetCommandArgs(commandArgs)
	if err != nil {
		b.error = fmt.Errorf("could not parse the jobcmd: %s", err)
		return b
	}

	if b.binary == "" {
		b.commands = commands
		return b
	}

	b.commands = append([]string{b.binary}, commands...)
	return b
}

func (b *CMDBuilder) Build() (*CMD, error) {
	if b.error != nil {
		return nil, errors.NewConfigurationError("Could not create a valid 'jobcmd' instance", b.error)
	}

	return &CMD{
		Binary:   b.binary,
		Commands: b.commands,
	}, nil
}

func NewCMD() *CMDBuilder {
	return &CMDBuilder{}
}
