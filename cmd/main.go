package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/hunterMotko/csvq/utils"
)

const version = "0.0.1"

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

// if -v flag is set, print version and exit
func getVersion() {
	if *v {
		fmt.Fprintf(os.Stdout, "go_csv-%s\n", version)
		os.Exit(1)
	}
}

type Csv struct {
  csvReader *csv.Reader
  Headers map[string]int
  Records [][]string
  Lines int
}

func (c *Csv) GetHeaders() []string {
  headers, err := c.csvReader.Read()
  if err != nil {
    log.Fatalf("error reading headers: %v", err)
  }
  for i, header := range headers {
    c.Headers[header] = i
  }
  return headers
}

func (c *Csv) GetRecords() {
  for {
    record, err := c.csvReader.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      log.Fatalf("error reading record: %v", err)
    }
    c.Records = append(c.Records, record)
  }
}

func (c *Csv) GetColumnsByIndex(args []string) {
  for _, arg := range args {
    if _, ok := c.Headers[arg]; !ok {
      log.Fatalf("header not found: %s", arg)
    }
  }
  for _, record := range c.Records {
    for _, arg := range args {
      fmt.Fprintf(os.Stdout, "%s,", record[c.Headers[arg]])
    }
    fmt.Fprintf(os.Stdout, "\n")
  }
}

func main() {
	flag.Parse()
	utils.PipeCheck()
	getVersion()
	csvReader := csv.NewReader(os.Stdin)
  var csvFile Csv
  csvFile.csvReader = csvReader
  csvFile.Headers = make(map[string]int)
  csvFile.Records = make([][]string, 0)
  headers := csvFile.GetHeaders()
  // if -hd flag is set, print headers and exit
  if *hd {
    fmt.Fprintf(os.Stdout, "%s%v%s\n", red, strings.Join(headers, ","), reset)
    os.Exit(1)
  }
  csvFile.GetRecords()
	// if -c flag is set, get columns by header name
	if *c {
		args := flag.Args()
		if len(args) == 0 {
			log.Fatalf("no columns specified")
		}
    csvFile.GetColumnsByIndex(args)
	}
	os.Exit(1)
}
