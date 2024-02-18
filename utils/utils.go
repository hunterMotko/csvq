package utils

import (
	"fmt"
	"os"
	"strings"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

const version = "0.0.1"
// if -v flag is set, print version and exit
func GetVersion() {
	fmt.Fprintf(os.Stdout, "csvq version %s\n", version)
	os.Exit(0)
}

// if -hd flag is set, print headers and exit
func PrintHeaders(headers []string) {
	fmt.Fprintf(os.Stdout, "%s%v%s\n", red, strings.Join(headers, ","), reset)
	os.Exit(0)
}
