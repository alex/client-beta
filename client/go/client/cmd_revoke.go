package client

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	rpc "github.com/keybase/go-framed-msgpack-rpc"
)

type CmdRevoke struct {
	id string
}

func (c *CmdRevoke) ParseArgv(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("Revoke takes exactly one key.")
	}
	c.id = ctx.Args()[0]
	return nil
}

func (c *CmdRevoke) Run() (err error) {
	cli, err := GetRevokeClient()
	if err != nil {
		return err
	}

	protocols := []rpc.Protocol{
		NewSecretUIProtocol(G),
	}
	if err = RegisterProtocols(protocols); err != nil {
		return err
	}

	return cli.RevokeKey(context.TODO(), keybase1.RevokeKeyArg{
		KeyID: c.id,
	})
}

func NewCmdRevoke(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:         "revoke",
		ArgumentHelp: "<key-id>",
		Usage:        "Revoke a key",
		Flags:        []cli.Flag{},
		Action: func(c *cli.Context) {
			cl.ChooseCommand(&CmdRevoke{}, "revoke", c)
		},
	}
}

func (c *CmdRevoke) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:     true,
		GpgKeyring: true,
		KbKeyring:  true,
		API:        true,
	}
}
