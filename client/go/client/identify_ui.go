package client

import (
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	"golang.org/x/net/context"
)

type IdentifyUIServer struct {
	ui libkb.IdentifyUI
}

func NewIdentifyUIProtocol() rpc.Protocol {
	return keybase1.IdentifyUiProtocol(&IdentifyUIServer{GlobUI.GetIdentifyUI()})
}

func NewIdentifyTrackUIProtocol() rpc.Protocol {
	ui := GlobUI.GetIdentifyTrackUI(true)
	return keybase1.IdentifyUiProtocol(&IdentifyUIServer{ui})
}

func (i *IdentifyUIServer) Confirm(_ context.Context, arg keybase1.ConfirmArg) (bool, error) {
	return i.ui.Confirm(&arg.Outcome)
}

func (i *IdentifyUIServer) FinishWebProofCheck(_ context.Context, arg keybase1.FinishWebProofCheckArg) error {
	i.ui.FinishWebProofCheck(arg.Rp, arg.Lcr)
	return nil
}

func (i *IdentifyUIServer) FinishSocialProofCheck(_ context.Context, arg keybase1.FinishSocialProofCheckArg) error {
	i.ui.FinishSocialProofCheck(arg.Rp, arg.Lcr)
	return nil
}

func (i *IdentifyUIServer) DisplayCryptocurrency(_ context.Context, arg keybase1.DisplayCryptocurrencyArg) error {
	i.ui.DisplayCryptocurrency(arg.C)
	return nil
}

func (i *IdentifyUIServer) DisplayKey(_ context.Context, arg keybase1.DisplayKeyArg) error {
	i.ui.DisplayKey(arg.Key)
	return nil
}

func (i *IdentifyUIServer) ReportLastTrack(_ context.Context, arg keybase1.ReportLastTrackArg) error {
	i.ui.ReportLastTrack(arg.Track)
	return nil
}

func (i *IdentifyUIServer) LaunchNetworkChecks(_ context.Context, arg keybase1.LaunchNetworkChecksArg) error {
	return nil
}

func (i *IdentifyUIServer) DisplayTrackStatement(_ context.Context, arg keybase1.DisplayTrackStatementArg) error {
	i.ui.DisplayTrackStatement(arg.Stmt)
	return nil
}

func (i *IdentifyUIServer) Start(_ context.Context, arg keybase1.StartArg) error {
	i.ui.Start(arg.Username)
	return nil
}

func (i *IdentifyUIServer) Finish(_ context.Context, sessionID int) error {
	i.ui.Finish()
	return nil
}
