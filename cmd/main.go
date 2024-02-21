package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"

	Csv "github.com/hunterMotko/csvq/cmd/csv"
	"github.com/hunterMotko/csvq/utils"
)

var (
	v  *bool
	hd *bool
	c  *bool
	f  string
	s  string
)

func init() {
	v = flag.Bool("v", false, "print version")
	c = flag.Bool("c", false, "get columns by header name")
	hd = flag.Bool("hd", false, "get csv headers")
	flag.StringVar(&f, "f", "", "file path")
	flag.StringVar(&s, "s", "", "slice of columns")
}

func main() {
	flag.Parse()
	if *v {
		utils.GetVersion()
	}

	var reader *csv.Reader
	if f != "" {
		file, err := os.Open(f)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	}

	if reader == nil {
		utils.PipeCheck()
		reader = csv.NewReader(os.Stdin)
	}

	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Error reading csv: %v\n", err)
	}

	if *hd {
		utils.PrintHeaders(headers)
	}

	csvFile, err := Csv.NewCsv(reader, headers)
	if err != nil {
		os.Exit(1)
	}

	if s != "" {
		csvFile.GetColumnsBySlice(s)
	} else if *c {
		args := flag.Args()
		csvFile.GetColumnsByIndex(args)
	}
}
