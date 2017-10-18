package appcli

import (
	"fmt"
	"strings"

	"github.com/rliebz/tusk/config"
	"github.com/urfave/cli"
)

// CompletionFlag is the flag passed when performing shell completions.
var CompletionFlag = "--" + cli.BashCompletionFlag.GetName()

// createBashComplete prints as the completion type then a list of options.
// Currently, the only available completion types are "file" and "tasks".
func createBashComplete(app *cli.App, meta *config.Metadata) func(c *cli.Context) {
	return func(c *cli.Context) {
		if c.NArg() > 0 {
			return
		}

		if !meta.Completion.IsFlagValue {
			fmt.Println("tasks")
			for _, command := range app.Commands {
				if command.Hidden {
					continue
				}
				fmt.Printf(
					"%s:%s\n",
					command.Name,
					strings.Replace(command.Usage, "\n", "", -1),
				)
			}
			for _, flag := range app.Flags {
				values := strings.Split(flag.GetName(), ", ")
				for _, value := range values {
					if len(value) == 1 || c.IsSet(value) {
						continue
					}
					fmt.Printf("--%s\n", value)
				}
			}
			return
		}

		// Default to file completion
		fmt.Println("file")
	}
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
