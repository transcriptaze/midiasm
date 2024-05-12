package main

import (
	"flag"
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
	{"disassemble", &commands.DISASSEMBLE},
	{"assemble", &commands.ASSEMBLE},
	{"notes", &commands.NOTES},
	{"click", &commands.CLICK},
	{"export", &commands.EXPORT},
	{"transpose", &commands.TRANSPOSE},
}

var options = struct {
	conf    string
	c4      bool
	verbose bool
	debug   bool
}{}

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
	cmd, flagset, err := parse()
	if err != nil {
		fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
		return
	}

	// ... set manufacturer specific options
	if options.conf != "" {
		f, err := os.Open(options.conf)
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
	if options.debug {
		log.SetLogLevel(log.Debug)
	} else if options.verbose {
		log.SetLogLevel(log.Info)
	}

	if options.c4 {
		context.SetMiddleC(lib.C4)
	} else {
		context.SetMiddleC(lib.C3)
	}

	// ... process
	if err := cmd.Execute(flagset); err != nil {
		fmt.Println()
		fmt.Printf("   *** ERROR: %v\n", err)
		fmt.Println()

		os.Exit(1)
	}
}

func parse() (commands.Command, *flag.FlagSet, error) {
	flagset := flag.NewFlagSet("midiasm", flag.ExitOnError)

	flagset.BoolVar(&options.c4, "C4", options.c4, "Sets middle C to C4 (Yamaho convention). Defaults to C3")
	flagset.BoolVar(&options.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&options.debug, "debug", false, "Enable debugging information")

	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.cmd == os.Args[1] {
				cmd := c.command
				flagset = cmd.Flagset(flagset)
				if err := flagset.Parse(os.Args[2:]); err != nil {
					return cmd, flagset, err
				} else {
					return cmd, flagset, nil
				}
			}
		}
	}

	cmd := &commands.DISASSEMBLE
	flagset = cmd.Flagset(flagset)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		return cmd, flagset, err
	}

	return cmd, flagset, nil
}
