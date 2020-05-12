package main

import (
	"clustercloner/clusters/launcher"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/urfave/cli/v2"
)

var (
	mainCtx context.Context
	//Version ...
	Version = "dev"
	//BuildDate build date
	BuildDate = "unknown"
	// GitCommit git commit SHA
	GitCommit = "dirty"
	//GitBranch git branch
	GitBranch = "master"
)

func mainCmd(cliCtx *cli.Context) error {
	printFlags(cliCtx)
	launcher.Launch(cliCtx)
	//crossCloud.PocLaunch()

	return nil
}

func printFlags(cliCtx *cli.Context) {
	var s = "\n"
	for _, flagName := range cliCtx.FlagNames() {
		value := cliCtx.String(flagName)
		s += fmt.Sprintf("\t\t%s: %s\n", flagName, value)
	}
	log.Println(s)
}

func init() {
	// handle termination signal
	mainCtx = handleSignals()
	_ = mainCtx
}

func handleSignals() context.Context {
	// Graceful shut-down on SIGINT/SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// create cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		sid := <-sig
		log.Printf("Received signal: %d\n", sid)
		log.Println("canceling main command ...")
	}()

	return ctx
}

func main() {
	log.Println("Starting")

	flags := launcher.CLIFlags()
	app := &cli.App{
		Flags: flags,
		Name:  "Cluster Cloner",
		Usage: "CLI",

		Action:  mainCmd,
		Version: Version,
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("goapp %s\n", Version)
		fmt.Printf("  Build date: %s\n", BuildDate)
		fmt.Printf("  Git commit: %s\n", GitCommit)
		fmt.Printf("  Git branch: %s\n", GitBranch)
		fmt.Printf("  Built with: %s\n", runtime.Version())
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}