package utils

import (
	"fmt"
	"os"
)

// check if stdin is being used and if so, check if it's being piped and not <, <<, <<<, etc.
func PipeCheck() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	if (stat.Mode() & os.ModeNamedPipe) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", "If using stdin, please pipe in a file.")
		os.Exit(-1)
	}
}
