package config

import (
	"os"
	"path/filepath"

	"github.com/Microsoft/kunlun/common/configuration"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/ui"
	flags "github.com/jessevdk/go-flags"
)

type fs interface {
	fileio.Stater
	fileio.TempFiler
	fileio.FileReader
	fileio.FileWriter
}

func NewConfig(ui *ui.UI, fs fs) Config {
	return Config{
		// stateBootstrap: bootstrap,
		ui: ui,
		fs: fs,
	}
}

type Config struct {
	ui *ui.UI
	fs fs
}

func ParseArgs(args []string) (GlobalFlags, []string, error) {
	var globals GlobalFlags
	parser := flags.NewParser(&globals, flags.IgnoreUnknown)
	remainingArgs, err := parser.ParseArgs(args[1:])
	if err != nil {
		return GlobalFlags{}, remainingArgs, err
	}

	if !filepath.IsAbs(globals.StateDir) {
		workingDir, err := os.Getwd()
		if err != nil {
			return GlobalFlags{}, remainingArgs, err
		}
		globals.StateDir = filepath.Join(workingDir, globals.StateDir)
	}

	return globals, remainingArgs, nil
}

func (c Config) Bootstrap(globalFlags GlobalFlags, remainingArgs []string, argsLen int) (configuration.Configuration, error) {
	if argsLen == 1 { // if run kid.
		return configuration.Configuration{
			Command: "help",
		}, nil
	}

	var command string
	if len(remainingArgs) > 0 {
		command = remainingArgs[0]
	}

	if globalFlags.Version || command == "version" {
		command = "version"
		return configuration.Configuration{
			ShowCommandHelp: globalFlags.Help,
			Command:         command,
		}, nil
	}

	if len(remainingArgs) == 0 {
		return configuration.Configuration{
			Command: "help",
		}, nil
	}

	if len(remainingArgs) == 1 && command == "help" {
		return configuration.Configuration{
			Command: command,
		}, nil
	}

	if command == "help" {
		return configuration.Configuration{
			ShowCommandHelp: true,
			Command:         remainingArgs[1],
		}, nil
	}

	if globalFlags.Help {
		return configuration.Configuration{
			ShowCommandHelp: true,
			Command:         command,
		}, nil
	}

	return configuration.Configuration{
		Global: configuration.GlobalConfiguration{
			Debug:    globalFlags.Debug,
			StateDir: globalFlags.StateDir,
			Name:     globalFlags.EnvID,
		},
		Command:         command,
		SubcommandFlags: remainingArgs[1:],
		ShowCommandHelp: false,
	}, nil
}
