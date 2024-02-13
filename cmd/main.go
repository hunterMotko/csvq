package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

const version = "0.0.1"

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	v      *bool
	c      *bool
	hd     *bool
	f      string
)

func init() {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		green = ""
		yellow = ""
		blue = ""
	}
	v = flag.Bool("v", false, "print version")
	c = flag.Bool("c", false, "get columns by header name")
	hd = flag.Bool("hd", false, "get csv headers")
	flag.StringVar(&f, "f", "", "file path")
}

func getCsvMatrix(r *csv.Reader) <-chan []string {
	c := make(chan []string)
	go func(c chan []string) {
		defer close(c)
		for {
			c <- readByLine(r)
		}
	}(c)
	return c
}

func readCsvFile(file *os.File) {
	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
    if err == io.EOF {
      return
    }
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(-1)
		}
		fmt.Println(record)
	}
}

func getIndexes(args []string, record []string) []int {
	var indexes []int
  fmt.Println(args, record)
	for _, arg := range args {
		for i, v := range record {
			if v == arg {
				indexes = append(indexes, i)
			}
		}
	}
	return indexes
}

func getColumn(idxs []int, record []string) {
  if len(idxs) == 0 || len(record) == 0 {
    os.Exit(0)
  }
	for _, idx := range idxs {
    if idx > len(record) {
      os.Exit(0)
    }
		fmt.Fprintf(os.Stdout, "%v ", record[idx])
	}
  fmt.Fprintf(os.Stdout, "\n")
}

// check if stdin is being used and if so, check if it's being piped and not <, <<, <<<, etc.
func pipeCheck() {
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

// if -v flag is set, print version and exit
func getVersion() {
	if *v {
		fmt.Fprintf(os.Stdout, "go_csv-%s\n", version)
		os.Exit(1)
	}
}

// if -hd flag is set, print headers and exit
func readByLine(cr *csv.Reader) []string {
  headers, err := cr.Read()
   if err == io.EOF {
     return nil
   }
   if err != nil {
     fmt.Fprintf(os.Stderr, "%v\n", err)
     os.Exit(-1)
   }
   return headers
}

func GetCols(args []string, headers []string, csvReader *csv.Reader) {
		var indexes []int
    indexes = getIndexes(args, headers)
		for record := range getCsvMatrix(csvReader) {
			getColumn(indexes, record)
		}
		os.Exit(1)
}

func main() {
	flag.Parse()
	pipeCheck()
	getVersion()
	var stdin *os.File = os.Stdin
	csvReader := csv.NewReader(stdin)
  var headers []string
  headers = readByLine(csvReader)
  // if -hd flag is set, print headers and exit
	if *hd {
    fmt.Fprintf(os.Stdout, "%s%v%s\n", red, headers, reset)
    os.Exit(1)
	}
	// if -c flag is set, get columns by header name
	if *c {
    args := flag.Args()
    fmt.Println(args)
    GetCols(args, headers, csvReader)
	}
	os.Exit(1)
}
