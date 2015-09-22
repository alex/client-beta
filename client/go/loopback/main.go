package keybase

import (
	"encoding/base64"
	"fmt"
	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/service"
	"net"
	"sync"
)

var val string
var con net.Conn
var startOnce sync.Once

type dummyCmd struct{}

type debuggingConfig struct {
	libkb.NullConfiguration
	homeDir *string
	runMode *string
}

func (n debuggingConfig) GetDebug() (bool, bool) {
	// if you want helpful debug info in xcode
	return true, true
	// return false, false
}

func (n debuggingConfig) GetLocalRPCDebug() string {
	// if you want helpful debug info in xcode
	return "Acsvip"
	// return ""
}

func (n debuggingConfig) GetRunMode() (libkb.RunMode, error) {
	if n.runMode == nil {
		return libkb.DevelRunMode, nil
	}

	return libkb.StringToRunMode(*n.runMode)
}

func (n debuggingConfig) GetHome() string {
	if n.homeDir == nil {
		return ""
	}

	return *n.homeDir
}

func (d dummyCmd) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:    true,
		API:       true,
		KbKeyring: true,
	}
}

func start(cmdline libkb.CommandLine) {
	startOnce.Do(func() {
		libkb.G.Init()
		libkb.G.ConfigureAll(cmdline, dummyCmd{})
		(service.NewService(false)).StartLoopbackServer(libkb.G)
		Reset()
	})
}

func Init(homeDir string, runMode string) {
	start(debuggingConfig{libkb.NullConfiguration{}, &homeDir, &runMode})
}

// Takes base64 encoded msgpack rpc payload
func WriteB64(str string) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err == nil {
		con.Write(data)
	} else {
		fmt.Println("write error:", err, str)
	}
}

// Blocking read, returns base64 encoded msgpack rpc payload
func ReadB64() string {
	data := make([]byte, 50*1024)

	n, err := con.Read(data)
	if n > 0 && err == nil {
		str := base64.StdEncoding.EncodeToString(data[0:n])
		return str
	}

	if err != nil {
		fmt.Println("read error:", err)
		// attempt to fix the connection
		Reset()
	}

	return ""
}

func Reset() {
	if con != nil {
		con.Close()
	}

	var err error
	libkb.G.SocketWrapper = nil
	con, _, err = libkb.G.GetSocket()

	if err != nil {
		fmt.Println("loopback socker error:", err)
	}
}