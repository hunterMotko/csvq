package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

func HandleSliceString(slice string, hLen int) (int, int, error) {
	slice = strings.Trim(slice, "[]")
	var err error
	var start, end int
	if len(slice) == 2 {
		sl := []rune(slice)
		if sl[0] == '-' {
			start = 0
			end = int(sl[1] - '0')
		} else if sl[1] == '-' {
			start = int(sl[0] - '0')
			end = hLen
		}
	} else {
		sliceStr := strings.Split(slice, "-")
		if len(sliceStr) != 2 {
			return 0, 0, fmt.Errorf("invalid slice: %s", slice)
		}
		start, err = strconv.Atoi(sliceStr[0])
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid Atoi Conversion: %s", slice)
		}
		end, err = strconv.Atoi(sliceStr[1])
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid Atoi Conversion: %s", slice)
		}
		if end > hLen {
			return start, end, fmt.Errorf("Index out of bounds: %d\n Your column length is:  %d\n", end, hLen)
		}
	}
	return start, end, nil
}

func (c *Csv) GetColumnsBySlice(s string) error {
	start, end, err := HandleSliceString(s, c.HeadLen)
	if err != nil {
		return err
	}

	var headers []string
	for k := range c.Headers {
		headers = append(headers, k)
	}

	sort.SliceStable(headers, func(i, j int) bool {
		return c.Headers[headers[i]] < c.Headers[headers[j]]
	})
	
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
