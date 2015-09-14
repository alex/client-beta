package client

import (
	"github.com/keybase/cli"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/protocol/go"
	"github.com/maxtaco/go-framed-msgpack-rpc/rpc2"
)

type CmdPGPPull struct {
	userAsserts []string
}

func (v *CmdPGPPull) ParseArgv(ctx *cli.Context) error {
	v.userAsserts = ctx.Args()
	return nil
}

func (v *CmdPGPPull) Run() (err error) {
	cli, err := GetPGPClient()
	if err != nil {
		return err
	}

	protocols := []rpc2.Protocol{
		NewLogUIProtocol(),
	}
	if err = RegisterProtocols(protocols); err != nil {
		return err
	}

	return cli.PGPPull(keybase1.PGPPullArg{
		UserAsserts: v.userAsserts,
	})
}

func NewCmdPGPPull(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:         "pull",
		ArgumentHelp: "<usernames...>",
		Usage:        "Download the latest PGP keys for people you track.",
		Flags:        []cli.Flag{},
		Action: func(c *cli.Context) {
			cl.ChooseCommand(&CmdPGPPull{}, "pull", c)
		},
	}
}

func (v *CmdPGPPull) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:     true,
		GpgKeyring: true,
		KbKeyring:  true,
		API:        true,
	}
}
