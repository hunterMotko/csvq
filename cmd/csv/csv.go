package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
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

type Csv struct {
	Reader  *csv.Reader
	Headers map[string]int
	Records [][]string
	Lines   int
}

func (c *Csv) Init() {
	c.Headers = make(map[string]int)
	c.Records = make([][]string, 0)
}

// SetReader sets the reader for the csv file
func (c *Csv) SetReader(r *csv.Reader) {
	c.Reader = r
}

func (c *Csv) GetHeaders() []string {
	headers, err := c.Reader.Read()
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
		record, err := c.Reader.Read()
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

// if -c flag is set, get columns by header name
func (c *Csv) HandleColumns(args []string) {
	if len(args) == 0 {
		log.Fatalf("no columns specified")
	}
	c.GetColumnsByIndex(args)
}
