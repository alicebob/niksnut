package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alicebob/niksnut/httpd"
	"github.com/alicebob/niksnut/niks"
)

var (
	defaultListen    = "localhost:3141"
	defaultBuildsDir = "./builds/"
)

func main() {
	cli := parseFlags()
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	// slog.Info("hello.")

	if cli.command == "help" {
		fmt.Printf(`usage: niksnut [--help] [--version] [--buildsdir=%s]
	       [--configfile=./config.nix]
	       <command> [--help] [<args>]
`, defaultBuildsDir)
		fmt.Printf("   niksnut help -- same as `niksnut --help`\n")
		fmt.Printf("   niksnut version -- same as `niksnut --version`\n")
		fmt.Printf("   niksnut check\n")
		fmt.Printf("   niksnut httpd [--listen=%s]\n", defaultListen)
		fmt.Printf("   niksnut run <projectname> [<git branch>]\n")
		return
	}
	if cli.version {
		fmt.Printf("niksnut version v0.0.1\n")
		return
	}

	config, err := niks.ReadConfig(cli.configFile)
	if err != nil {
		fmt.Printf("error reading %s: %s\n", cli.configFile, err)
		os.Exit(1)
	}

	switch cli.command {
	case "check":
		if cli.help {
			fmt.Printf("usage: niksnut check\n")
			return
		}
		fmt.Printf("configfile: %s\n", cli.configFile)
		fmt.Printf("configfile parsed successfully\n")
		fmt.Printf("found %d projects:\n", len(config.Projects))
		for _, p := range config.Projects {
			fmt.Printf("  - %s\n", p.ID)
		}
	case "run":
		if cli.help {
			fmt.Printf("usage: niksnut run <projectid> [<git branch>]\n")
			return
		}
		cliRun(cli.buildsDir, config, cli.runProject)
	case "httpd":
		if cli.help {
			fmt.Printf("usage: niksnut httpd [--listen=%s]\n", defaultListen)
			return
		}

		s := &httpd.Server{
			BuildsDir: cli.buildsDir,
			Root:      os.DirFS("./httpd/"), // FIXME
			Config:    *config,
			Addr:      cli.httpdListen,
		}

		fmt.Printf("starting httpd on %s\n", s.Addr)
		go func() {
			err := s.Run()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

	}
}

type cliFlags struct {
	configFile  string
	buildsDir   string
	command     string
	httpdListen string
	runProject  string
	help        bool
	version     bool
}

func parseFlags() *cliFlags {
	var (
		fl   = &cliFlags{}
		args = os.Args[1:]
	)

	f := &flag.FlagSet{}
	f.StringVar(&fl.configFile, "config", "./config.nix", "config file (.nix)")
	f.StringVar(&fl.buildsDir, "buildsdir", defaultBuildsDir, "builds directory")
	f.BoolVar(&fl.help, "help", false, "run help")
	f.BoolVar(&fl.version, "version", false, "show version")
	if err := f.Parse(args); err != nil {
		fl.help = true
		return fl
	}
	args = f.Args()

	if len(args) == 0 {
		fl.command = "help"
		return fl
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "run":
		fl.command = "run"
		if len(args) != 1 {
			fl.help = true
			return fl
		}
		fl.runProject = args[0]
		// optional branch name
	case "httpd":
		fl.command = "httpd"
		f := &flag.FlagSet{}
		f.StringVar(&fl.httpdListen, "listen", defaultListen, "listen")
		f.BoolVar(&fl.help, "help", false, "run help")
		if err := f.Parse(args); err != nil {
			fl.help = true
			return fl
		}
		if len(f.Args()) != 0 {
			fl.help = true
			return fl
		}
	case "version":
		fl.command = "version"
		if len(f.Args()) != 0 {
			fl.help = true
			return fl
		}
	case "check":
		fl.command = "check"
		if len(f.Args()) != 0 {
			fl.help = true
			return fl
		}
	default:
		// "help" and invalid commands
		fl.command = "help"
	}

	if !strings.HasSuffix(fl.buildsDir, "/") {
		fl.buildsDir += "/"
	}

	return fl
}

func cliRun(buildsDir string, config *niks.Config, projectID string) {
	fmt.Printf("run of %s:\n", projectID)
	// fixme all
	p := config.Projects[0]
	branch := "main"
	build, err := niks.SetupBuild(buildsDir, p)
	if err != nil {
		fmt.Printf("error setting up build %s: %s\n", p.ID, err)
		os.Exit(1)
	}
	if err := build.Run(p, branch); err != nil {
		fmt.Printf("error running %s: %s\n", p.ID, err)
		os.Exit(1)
	}
	fmt.Printf("run of %s:\n", p.ID)

	fmt.Print(build.Stdout())
}
