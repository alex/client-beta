package client

import (
	keybase1 "github.com/keybase/client/go/protocol"
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	"golang.org/x/net/context"
)

func newChangeArg(newPassphrase string, force bool) keybase1.PassphraseChangeArg {
	return keybase1.PassphraseChangeArg{
		Passphrase: newPassphrase,
		Force:      force,
	}
}

func passphraseChange(arg keybase1.PassphraseChangeArg) error {
	cli, err := GetAccountClient()
	if err != nil {
		return err
	}
	protocols := []rpc.Protocol{
		NewSecretUIProtocol(G),
	}
	if err := RegisterProtocols(protocols); err != nil {
		return err
	}
	return cli.PassphraseChange(context.TODO(), arg)
}
