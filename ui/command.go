package ui

const (
	commandActionString = "Running"
	skippedString       = "Skipping"

	outputPrefix = "=> "
)

// PrintCommand prints the command to be executed.
func PrintCommand(command string) {
	if Quiet {
		return
	}

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
		"[%s] %s\n%s%s\n",
		yellow(skippedString),
		bold(command),
		cyan(prefixOutput()),
		reason,
	)
}

// PrintCommandError prints an error from a running command.
func PrintCommandError(err error) {
	if Quiet {
		return
	}

	printf(
		Stderr,
		"%s%s\n",
		red(prefixOutput()),
		err.Error(),
	)
}

func prefixOutput() string {
	if Quiet {
		return ""
	}

	return outputPrefix
}
