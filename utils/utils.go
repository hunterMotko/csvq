package utils

import (
	"fmt"
	"os"
	"strconv"
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

func StrToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Invalid Atoi Conversion: %s", s)
	}
	return i, nil
}

func HandleSliceString(slice string, hLen int) (int, int, error) {
	if !strings.Contains(slice, "-") {
		return 0, 0, fmt.Errorf("invalid slice: %s", slice)
	}
	front := slice[0]
	back := slice[len(slice)-1]
	sliceStr := strings.Split(slice, "-")
	sStrLen := len(sliceStr)
	if sStrLen != 2 {
		return 0, 0, fmt.Errorf("invalid slice: %s", slice)
	}
	var err error
	var start, end int
	if front == '-' {
		start = 0
		if end, err = StrToInt(sliceStr[1]); err != nil {
			return 0, 0, err
		}
	} else if back == '-' {
		if start, err = StrToInt(sliceStr[0]); err != nil {
			return 0, 0, err
		}
		end = hLen
	} else {
		if start, err = StrToInt(sliceStr[0]); err != nil {
			return 0, 0, err
		}
		if end, err = StrToInt(sliceStr[1]); err != nil {
			return 0, 0, err
		}
	}
	if end > hLen {
		return 0, 0, fmt.Errorf("Index out of bounds: %d\n Your column length is:  %d\n", end, hLen)
	}
	return start, end, nil
}

