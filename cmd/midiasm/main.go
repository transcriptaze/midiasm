package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"plugin"

	"github.com/transcriptaze/midiasm/commands"
	"github.com/transcriptaze/midiasm/log"
	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

var cli = []struct {
	cmd     string
	command commands.Command
}{
	{"disassemble", &DISASSEMBLE},
	{"assemble", &ASSEMBLE},
	{"notes", &NOTES},
	{"click", &CLICK},
	{"export", &EXPORT},
	{"transpose", &TRANSPOSE},
}

const version = "v0.1.0"

type Plugin interface {
	GetCommand() (string, commands.Command)
}

func main() {
	// ... load plugins
	bindir := filepath.Dir(os.Args[0])
	plugins := filepath.Join(bindir, "plugins")

	fs.WalkDir(os.DirFS(plugins), ".", func(path string, d fs.DirEntry, err error) error {
		file := filepath.Join(plugins, path)

		if err != nil {
			return err
		} else if !d.Type().IsRegular() {
			return nil
		} else if p, err := plugin.Open(file); err != nil {
			fmt.Printf("Error loading plugin %q (%v)", path, err)
		} else if tsv, err := p.Lookup("TSV"); err != nil {
			fmt.Printf("Error loading plugin %q (%v)", path, err)
		} else if plugin, ok := tsv.(Plugin); ok {
			cmd, command := plugin.GetCommand()

			cli = append(cli, struct {
				cmd     string
				command commands.Command
			}{
				cmd:     cmd,
				command: command,
			})
		}

		return nil
	})

	// ... add 'help' and 'version' commands to CLI

	cli = append(cli, struct {
		cmd     string
		command commands.Command
	}{
		cmd:     "help",
		command: &HELP,
	})

	cli = append(cli, struct {
		cmd     string
		command commands.Command
	}{
		cmd:     "version",
		command: &VERSION,
	})

	// ... parse command line
	cmd, err := parse()
	if err != nil {
		fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
		return
	}

	// ... set manufacturer specific options
	if conf := cmd.Config(); conf != "" {
		f, err := os.Open(conf)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		manufacturers, err := lib.LoadManufacturers(f)

		f.Close()

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		lib.AddManufacturers(manufacturers)
	}

	// ... set global stuff

	if cmd.Debug() {
		log.SetLogLevel(log.Debug)
	} else if cmd.Verbose() {
		log.SetLogLevel(log.Info)
	}

	context.SetMiddleC(cmd.MiddleC())

	// ... process
	if err := cmd.Execute(); err != nil {
		fmt.Println()
		fmt.Printf("   *** ERROR: %v\n", err)
		fmt.Println()

		os.Exit(1)
	}
}

func parse() (commands.Command, error) {
	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.cmd == os.Args[1] {
				cmd := c.command
				flagset := cmd.Flagset()
				if err := flagset.Parse(os.Args[2:]); err != nil {
					return cmd, err
				} else {
					return cmd, nil
				}
			}
		}
	}

	cmd := &DISASSEMBLE
	flagset := cmd.Flagset()
	if err := flagset.Parse(os.Args[1:]); err != nil {
		return cmd, err
	}

	return cmd, nil
}
