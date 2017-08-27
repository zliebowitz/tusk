package ui

const (
	commandActionString = "Running"
	skippedString       = "Skipping"

	outputPrefix = "  =>"
)

// PrintCommand prints the command to be executed.
func PrintCommand(command string) {
	printf(
		Stderr,
		"[%s] %s\n",
		blue(commandActionString),
		bold(command),
	)
}

// PrintSkipped prints the command skipped and the reason.
func PrintSkipped(command string, reason string) {
	if !Verbose {
		return
	}

	printf(
		Stderr,
		"[%s] %s\n%s %s\n",
		yellow(skippedString),
		bold(command),
		cyan(outputPrefix),
		reason,
	)
}

// PrintCommandOutput prints output from a running command.
func PrintCommandOutput(text string) {
	printf(
		Stderr,
		"%s %s\n",
		cyan(outputPrefix),
		text,
	)
}

// PrintCommandError prints an error from a running command.
func PrintCommandError(err error) {
	printf(
		Stderr,
		"%s [%s] %s\n",
		red(outputPrefix),
		red(errorString),
		err.Error(),
	)
}
