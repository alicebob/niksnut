package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alicebob/niksnut/httpd"
	"github.com/alicebob/niksnut/niks"
)

var (
	//go:embed static
	staticRoot embed.FS
)

var (
	defaultListen    = "localhost:3141"
	defaultBuildsDir = "./builds/"

	MaxBuildAge = 4 * 24 * time.Hour
)

func main() {
	cli := parseFlags()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil))) // TODO: different in httpd mode?

	if cli.command == "help" {
		fmt.Printf(`usage: niksnut [--help] [--version] [--buildsdir=%s]
	       [--configfile=./config.nix] [--root=""] [--offline]
	       <command> [--help] [<args>]
`, defaultBuildsDir)
		fmt.Printf("   niksnut help -- same as `niksnut --help`\n")
		fmt.Printf("   niksnut version -- same as `niksnut --version`\n")
		fmt.Printf("   niksnut check\n")
		fmt.Printf("   niksnut httpd [--listen=%s]\n", defaultListen)
		fmt.Printf("   niksnut run <projectname> [<git branch>]\n")
		fmt.Printf("   niksnut gc\n")
		return
	}
	if cli.version {
		fmt.Printf("niksnut version v0.0.1\n")
		return
	}

	var config *niks.Config
	if !cli.help {
		c, err := niks.ReadConfig(cli.configFile)
		if err != nil {
			fmt.Printf("error reading %s: %s\n", cli.configFile, err)
			os.Exit(1)
		}
		config = c
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
			fmt.Printf("usage: niksnut run <projectid> <git branch>\n")
			return
		}
		cliRun(cli.buildsDir, config, cli.runProject, cli.branch)
	case "httpd":
		if cli.help {
			fmt.Printf("usage: niksnut httpd [--listen=%s]\n", defaultListen)
			return
		}

		// use embedded templates and /static/ files, unless --root set
		var (
			root   = fs.FS(httpd.TemplateRoot)
			static = fs.FS(staticRoot)
		)
		if cli.devRoot != "" {
			slog.Info("dev mode because -root is set.")
			root = os.DirFS(cli.devRoot + "/httpd/")
			static = os.DirFS(cli.devRoot + "/")
		}

		s := &httpd.Server{
			BuildsDir: cli.buildsDir,
			Root:      root,
			Static:    static,
			Config:    *config,
			Addr:      cli.httpdListen,
			Offline:   cli.offline,
		}

		wg := sync.WaitGroup{}
		defer wg.Wait()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		slog.Info("starting httpd", "addr", s.Addr)
		go func() {
			err := s.Run()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}()

		wg.Add(1)
		go func() {
			bgGC(ctx, s.BuildsDir)
			wg.Done()
		}()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		cancel()

	case "gc":
		if cli.help {
			fmt.Printf("usage: niksnut gc\nremoved old builds and runs nix-collect-garbage")
			return
		}
		cliGC(cli.buildsDir)
	}
}

type cliFlags struct {
	configFile  string
	buildsDir   string
	devRoot     string
	command     string
	httpdListen string
	runProject  string
	branch      string
	offline     bool
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
	f.StringVar(&fl.buildsDir, "buildsdir", defaultBuildsDir, "builds and checkouts directory")
	f.StringVar(&fl.devRoot, "root", "", "if empty, use embedded templates and /static/ files")
	f.BoolVar(&fl.offline, "offline", false, "offline git")
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
		if len(args) != 2 {
			fl.help = true
			return fl
		}
		fl.runProject = args[0]
		fl.branch = args[1]
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
		fl.version = true
		if len(f.Args()) != 1 {
			fl.help = true
			return fl
		}
	case "check", "gc":
		fl.command = cmd
		if len(f.Args()) != 1 {
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

func cliRun(buildsDir string, config *niks.Config, projectID, branch string) {
	ctx := context.Background()
	fmt.Printf("run of %s - %s:\n", projectID, branch)
	var proj *niks.Project
	for _, p := range config.Projects {
		if p.ID == projectID {
			proj = &p
			break
		}
	}
	if proj == nil {
		fmt.Printf("project not found\n")
		os.Exit(1)
	}

	build, err := niks.SetupBuild(buildsDir, *proj)
	if err != nil {
		fmt.Printf("error setting up build %s: %s\n", proj.ID, err)
		os.Exit(1)
	}
	if err := build.Run(ctx, buildsDir, *proj, branch); err != nil {
		fmt.Printf("error running %s: %s\n", proj.ID, err)
		os.Exit(1)
	}

	fmt.Print(build.Stdout())
}

func cliGC(root string) {
	ctx := context.Background()
	err := niks.GarbageCollect(ctx, root, time.Now().UTC().Add(-MaxBuildAge))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func bgGC(ctx context.Context, buildsDir string) {
	// slog.Info("bg GC loop")
	// defer slog.Info("bg GC loop done")
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
			keep := time.Now().UTC().Add(-MaxBuildAge)
			if err := niks.GarbageCollect(ctx, buildsDir, keep); err != nil {
				slog.Error("background GC run", "error", err)
			}
		}
	}
}
