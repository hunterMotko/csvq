package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/hunterMotko/csvq/utils"
	Csv "github.com/hunterMotko/csvq/cmd/csv"
)


var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"

	v  *bool
	hd *bool
	c  *bool
	f  string
	s  string
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
	flag.StringVar(&s, "s", "", "slice of columns")
}

func main() {
	var in *os.File = os.Stdin
	flag.Parse()
	if *v {
		utils.GetVersion()
	}

	var csvFile Csv.Csv
	csvFile.Init()

	if f != "" {
		file, err := os.Open(f)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
		csvFile.SetReader(csv.NewReader(file))
	}

	if csvFile.Reader == nil {
		utils.PipeCheck()
		csvFile.SetReader(csv.NewReader(in))
	}

	headers := csvFile.GetHeaders()
	if *hd {
		utils.PrintHeaders(headers)
	} else if s != "" {
		csvFile.GetRecords()
		csvFile.GetColumnsBySlice(s, headers)
	} else if *c {
		csvFile.GetRecords()
		csvFile.HandleColumns(flag.Args())
	}
}
