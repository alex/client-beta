// +build darwin

package client

import (
	"fmt"
	"os"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/launchd"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
)

func NewCmdLaunchd(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "launchd",
		Usage:       "keybase launchd [...]",
		Description: "Manage keybase launchd services.",
		Subcommands: []cli.Command{
			NewCmdLaunchdInstall(cl),
			NewCmdLaunchdUninstall(cl),
			NewCmdLaunchdList(cl),
			NewCmdLaunchdStatus(cl),
			NewCmdLaunchdStart(cl),
			NewCmdLaunchdStop(cl),
			NewCmdLaunchdRestart(cl),
		},
	}
}

func NewCmdLaunchdInstall(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "install",
		Usage:       "keybase launchd install <label> <path/to/keybase>",
		Description: "Install a keybase launchd service.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "run-mode",
				Usage: fmt.Sprintf("Run mode (%s)", libkb.RunModes),
			},
		},
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) < 1 {
				G.Log.Fatalf("No label specified.")
			}
			if len(args) < 2 {
				G.Log.Fatalf("No path to keybase executable specified.")
			}
			plistArgs := []string{"--log-format=file"}

			runMode := c.String("run-mode")
			if runMode != "" {
				plistArgs = append(plistArgs, fmt.Sprintf("--run-mode=%s", runMode))
			}

			plistArgs = append(plistArgs, "service")

			plist := launchd.NewPlist(args[0], args[1], plistArgs)
			err := launchd.Install(plist)
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func NewCmdLaunchdUninstall(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "uninstall",
		Usage:       "keybase launchd uninstall <label>",
		Description: "Uninstall a keybase launchd service.",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) < 1 {
				G.Log.Fatalf("No label specified.")
			}
			err := launchd.Uninstall(args[0])
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func NewCmdLaunchdList(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "list",
		Usage:       "keybase launchd list",
		Description: "List keybase launchd services.",
		Action: func(c *cli.Context) {
			err := launchd.ShowServices("keybase")
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func NewCmdLaunchdStatus(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "status",
		Usage:       "keybase launchd status <label>",
		Description: "Status for keybase launchd service.",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) < 1 {
				G.Log.Fatalf("No label specified.")
			}
			err := launchd.ShowStatus(args[0])
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func NewCmdLaunchdRestart(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "restart",
		Usage:       "keybase launchd restart <label>",
		Description: "Restart a keybase launchd service.",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) < 1 {
				G.Log.Fatalf("No label specified.")
			}
			err := launchd.Restart(args[0])
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func NewCmdLaunchdStart(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "start",
		Usage:       "keybase launchd start <label>",
		Description: "Start a keybase launchd service",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) < 1 {
				G.Log.Fatalf("No label specified")
			}
			err := launchd.Start(args[0])
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func NewCmdLaunchdStop(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "stop",
		Usage:       "keybase launchd stop <label>",
		Description: "Stop a keybase launchd service.",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) < 1 {
				G.Log.Fatalf("No label specified.")
			}
			err := launchd.Stop(args[0])
			if err != nil {
				G.Log.Fatalf("%v", err)
			}
			os.Exit(0)
		},
	}
}

func DiagnoseSocketError(err error) {
	services, err := launchd.ListServices("keybase.")
	if err != nil {
		GlobUI.Printf("Error checking launchd services: %v\n\n", err)
		return
	}

	if len(services) == 0 {
		GlobUI.Println("\nThere are no Keybase services installed. You may need to re-install.")
	} else if len(services) > 1 {
		GlobUI.Println("\nWe found multiple services:")
		for _, service := range services {
			GlobUI.Println("  " + service.StatusDescription())
		}
		GlobUI.Println("")
	} else if len(services) == 1 {
		service := services[0]
		status, err := service.Status()
		if err != nil {
			G.Log.Errorf("Error checking service status(%s): %v\n\n", service.Label(), err)
		} else {
			if status == nil || !status.IsRunning() {
				GlobUI.Printf("\nWe found a Keybase service (%s) but it's not running.\n", service.Label())
				cmd := fmt.Sprintf("keybase launchd start %s", service.Label())
				GlobUI.Println("You might try starting it: " + cmd + "\n")
			} else {
				GlobUI.Printf("\nWe couldn't connect but there is a Keybase service (%s) running (%s).\n", status.Label(), status.Pid())
				cmd := fmt.Sprintf("keybase launchd restart %s", service.Label())
				GlobUI.Println("You might try restarting it: " + cmd + "\n")
			}
		}
	}
}
