package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	"github.com/hunterMotko/csvq/utils"
)

type Csv struct {
	Headers map[string]int
	HeadLen int
	Records [][]string
	Lines   int
}

func NewCsv(reader *csv.Reader, headers []string) (*Csv, error) {
	head := make(map[string]int)
	hLen := 0
	for i, v := range headers {
		head[v] = i
		hLen++
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading csv: %v\n", err)
	}
	return &Csv{
		Headers: head,
		HeadLen: hLen,
		Records: records,
		Lines:   len(records),
	}, nil
}

func (c *Csv) GetColumnsByIndex(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no columns specified")
	}
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()
	for _, h := range args {
		if _, ok := c.Headers[h]; !ok {
			return fmt.Errorf("header not found: %s", h)
		}
	}
	if err := writer.Write(args); err != nil {
		return err
	}
	temp := make([]string, len(args))
	for _, rec := range c.Records {
		for i, arg := range args {
			temp[i] = rec[c.Headers[arg]]
		}
		if err := writer.Write(temp); err != nil {
			return err
		}
	}
	return nil
}

func (c *Csv) HeaderToArray() []string {
	var headers []string
	for k := range c.Headers {
		headers = append(headers, k)
	}
	sort.SliceStable(headers, func(i, j int) bool {
		return c.Headers[headers[i]] < c.Headers[headers[j]]
	})
	return headers
}

func (c *Csv) GetColumnsBySlice(s string) error {
	start, end, err := utils.HandleSliceString(s, c.HeadLen)
	if err != nil {
		return err
	}

	headers := c.HeaderToArray()
	writer := csv.NewWriter(os.Stdout)
	if err := writer.Write(headers[start:end]); err != nil {
		return err
	}
	for _, rec := range c.Records {
		if err := writer.Write(rec[start:end]); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("Write ERROR: %v", err)
	}
	return nil
}
