package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Csv struct {
	Headers map[string]int
	Records [][]string
	Lines   int
}

func NewCsv(reader *csv.Reader, headers []string) (*Csv, error) {
	head := make(map[string]int)
	for i, v := range headers {
		head[v] = i
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading csv: %v\n", err)
	}

	return &Csv{
		Headers: head,
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
	temp := make([]string, len(args))
	for _, rec := range c.Records {
		for i, arg := range args {
			if idx, ok := c.Headers[arg]; ok {
				temp[i] = rec[idx]
			} else {
				return fmt.Errorf("header not found: %s", arg)
			}
		}
		if err := writer.Write(temp); err != nil {
			return err
		}
	}

	return nil
}

func (c *Csv) GetColumnsBySlice(slice string) error {
	slice = strings.Trim(slice, "[]")
	sliceStr := strings.Split(slice, ":")
	if len(sliceStr) != 2 {
		return fmt.Errorf("invalid slice: %s", slice)
	}
	start, err := strconv.Atoi(sliceStr[0])
	if err != nil {
		return fmt.Errorf("Invalid Atoi Conversion: %s", slice)
	}
	end, err := strconv.Atoi(sliceStr[1])
	if err != nil {
		return fmt.Errorf("Invalid Atoi Conversion: %s", slice)
	}
	var headers []string
	for k := range c.Headers {
		headers = append(headers, k)
	}

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
	// 	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(headers[start:end], ","))
	// 	for _, record := range c.Records {
	// 		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(record[start:end], ","))
	// 	}
	return nil
}
