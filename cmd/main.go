package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
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

type Csv struct {
	reader  *csv.Reader
	Headers map[string]int
	Records [][]string
	Lines   int
}

func (csv *Csv) Init() {
	csv.Headers = make(map[string]int)
	csv.Records = make([][]string, 0)
}

// SetReader sets the reader for the csv file
func (csv *Csv) SetReader(r *csv.Reader) {
	csv.reader = r
}

func (c *Csv) GetHeaders() []string {
	headers, err := c.reader.Read()
	if err != nil {
		log.Fatalf("error reading headers: %v", err)
	}
	for i, header := range headers {
		c.Headers[header] = i
	}
	return headers
}

// GetRecords gets all records/body from csv file
func (c *Csv) GetRecords() {
	for {
		record, err := c.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading record: %v", err)
		}
		c.Records = append(c.Records, record)
	}
}

// GetColumnsByIndex gets columns by header index
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

// GetColumnsBySlice gets columns by slice 
// slice should be in the form of [start:end]
func (c *Csv) GetColumnsBySlice(slice string, headers []string) {
	slice = strings.Trim(slice, "[]")
	sliceStr := strings.Split(slice, ":")
	if len(sliceStr) != 2 {
		log.Fatalf("invalid slice: %s", slice)
	}
	start, err := strconv.Atoi(sliceStr[0])
	if err != nil {
		log.Fatalf("Invaild Atoi Conversion: %s", slice)
	}
	end, err := strconv.Atoi(sliceStr[1])
	if err != nil {
		log.Fatalf("Invaild Atoi Conversion: %s", slice)
	}
	var res []string
	h := strings.Join(headers[start:end], ",")
	res = append(res, h)
	for _, record := range c.Records {
		res = append(res, strings.Join(record[start:end], ","))
	}
	for i, row := range res {
		if i == 0 {
      fmt.Fprintf(os.Stdout, "%s%v%s\n", red, row, reset)
		} else {
      fmt.Fprintf(os.Stdout, "%s%v%s\n", green, row, reset)
		}
	}
}

// if -v flag is set, print version and exit
func getVersion() {
	fmt.Fprintf(os.Stdout, "csvq version %s\n", version)
	os.Exit(0)
}

// if -hd flag is set, print headers and exit
func PrintHeaders(headers []string) {
	fmt.Fprintf(os.Stdout, "%s%v%s\n", red, strings.Join(headers, ","), reset)
	os.Exit(0)
}

// if -c flag is set, get columns by header name
func HandleColumns(csvFile Csv) {
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("no columns specified")
	}
	csvFile.GetColumnsByIndex(args)
}

func main() {
	var in *os.File = os.Stdin
	flag.Parse()
	if *v {
		getVersion()
	}

	var csvFile Csv
	csvFile.Init()

	if f != "" {
		file, err := os.Open(f)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
		csvFile.SetReader(csv.NewReader(file))
	}

	if csvFile.reader == nil {
		utils.PipeCheck()
		csvFile.SetReader(csv.NewReader(in))
	}

	headers := csvFile.GetHeaders()
	if *hd {
		PrintHeaders(headers)
	} else if s != "" {
		csvFile.GetRecords()
		csvFile.GetColumnsBySlice(s, headers)
	} else if *c {
		csvFile.GetRecords()
		HandleColumns(csvFile)
	}
}
