package main

import (
	aks_ "clusterCloner/clusters/aks"
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
	// main context
	mainCtx context.Context
	//Version version of app
	Version = "dev"
	//BuildDate build date
	BuildDate = "unknown"
	// GitCommit git commit SHA
	GitCommit = "dirty"
	//GitBranch git branch
	GitBranch = "master"
)

func mainCmd(cliCtx *cli.Context) error {
	var s = ""
	for _, flagName := range cliCtx.FlagNames() {
		value := cliCtx.String(flagName)
		s += fmt.Sprintf("\t\t%s: %s\n", flagName, value)
	}
	aks_.CreateCluster()

	return nil
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
		log.Printf("received signal: %d\n", sid)
		log.Println("canceling main command ...")
	}()

	return ctx
}

func main() {
	log.Print("Starting")

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Usage:    "GCP project",
				Required: true, //todo use current GCP default
			},
			&cli.StringFlag{
				Name:  "location",
				Usage: "GCP zone",
			},
		},
		Name:    "goapp",
		Usage:   "goapp CLI",
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
