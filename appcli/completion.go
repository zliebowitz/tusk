package appcli

import (
	"fmt"
	"strings"

	"github.com/rliebz/tusk/config"
	"github.com/urfave/cli"
)

// CompletionFlag is the flag passed when performing shell completions.
var CompletionFlag = "--" + cli.BashCompletionFlag.GetName()

// createDefaultComplete prints the completion metadata for the top-level app.
// The metadata includes the completion type followed by a list of options.
// The available completion types are "normal" and "file". Normal will return
// tasks and flags, while file allows completion engines to use system files.
func createDefaultComplete(app *cli.App, meta *config.Metadata) func(c *cli.Context) {
	return func(c *cli.Context) {
		if c.NArg() > 0 {
			return
		}

		if !meta.Completion.IsFlagValue {
			fmt.Println("normal")
			for _, command := range app.Commands {
				printCommand(command)
			}
			for _, flag := range app.Flags {
				printFlag(c, flag)
			}
			return
		}

		// Default to file completion
		fmt.Println("file")
	}
}

// createCommandComplete prints the completion metadata for a cli command.
// The metadata includes the completion type followed by a list of options.
// The available completion types are "normal" and "file". Normal will return
// task-specific flags, while file allows completion engines to use system files.
func createCommandComplete(command *cli.Command, meta *config.Metadata) func(c *cli.Context) {
	return func(c *cli.Context) {
		if !meta.Completion.IsFlagValue {
			fmt.Println("normal")
			for _, flag := range command.Flags {
				printFlag(c, flag)
			}
			return
		}

		// Default to file completion
		fmt.Println("file")
	}
}

func printCommand(command cli.Command) {
	if command.Hidden {
		return
	}
	fmt.Printf(
		"%s:%s\n",
		command.Name,
		strings.Replace(command.Usage, "\n", "", -1),
	)
}

func printFlag(c *cli.Context, flag cli.Flag) {
	values := strings.Split(flag.GetName(), ", ")
	for _, value := range values {
		if len(value) == 1 || c.IsSet(value) {
			continue
		}
		fmt.Printf(
			"--%s:%s\n",
			value,
			strings.Replace(getDescription(flag), "\n", "", -1),
		)
	}
}

func getDescription(flag cli.Flag) string {
	return strings.SplitN(flag.String(), "\t", 2)[1]
}

func removeCompletionArg(args []string) ([]string, bool) {
	var output []string
	isCompleting := false
	for _, arg := range args {
		if arg != CompletionFlag {
			output = append(output, arg)
		} else {
			isCompleting = true
		}
	}

	return output, isCompleting
}
