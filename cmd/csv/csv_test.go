package csv_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	Csv "github.com/hunterMotko/csvq/cmd/csv"
)

var data []string = []string{
	"LatD,LatM,LatS,NS,LonD,LonM,LonS,EW,City,State",
	"41,5,59,N,80,39,0,W,Youngstown,OH",
	"42,52,48,N,97,23,23,W,Yankton,SD",
	"46,35,59,N,120,30,36,W,Yakima,WA",
	"42,16,12,N,71,48,0,W,Worcester,MA",
	"43,37,48,N,89,46,11,W,Wisconsin Dells,WI",
	"36,5,59,N,80,15,0,W,Winston-Salem,NC",
	"49,52,48,N,97,9,0,W,Winnipeg,MB",
	"39,11,23,N,78,9,36,W,Winchester,VA",
	"34,14,24,N,77,55,11,W,Wilmington,NC",
	"39,45,0,N,75,33,0,W,Wilmington,DE",
	"48,9,0,N,103,37,12,W,Williston,ND",
	"41,15,0,N,77,0,0,W,Williamsport,PA",
	"37,40,48,N,82,16,47,W,Williamson,WV",
	"33,54,0,N,98,29,23,W,Wichita Falls,TX",
	"37,41,23,N,97,20,23,W,Wichita,KS",
	"40,4,11,N,80,43,12,W,Wheeling,WV",
	"26,43,11,N,80,3,0,W,West Palm Beach,FL",
	"47,25,11,N,120,19,11,W,Wenatchee,WA",
	"41,25,11,N,122,23,23,W,Weed,CA",
}

var expHeaders map[string]int = map[string]int{
	"LatD": 0, "LatM": 1, "LatS": 2, "NS": 3, "LonD": 4, "LonM": 5, "LonS": 6, "EW": 7, "City": 8, "State": 9,
}

func TestCsvInit(t *testing.T) {
	cv := initCsv()
	t.Run("Check inital headers on NewCsv", func(t *testing.T) {
		for k, v := range cv.Headers {
			val, ok := expHeaders[k]
			if !ok {
				t.Errorf("Not in map! %s, %v", k, ok)
			}
			if val != v {
				t.Errorf("Val %d, not equal to %d", val, v)
			}
		}
	})
	t.Run("Check inital columns on NewCsv", func(t *testing.T) {
		exLen := 19
		exCols := data[1:]
		if cv.Lines != exLen {
			t.Errorf("lines %v, should equal %v\n", cv.Lines, exLen)
		}
		for i, rec := range cv.Records {
			col := strings.Split(exCols[i], ",")
			if !reflect.DeepEqual(rec, col) {
				t.Errorf("column %v,\n should equal %v\n", rec, col)
			}
		}
	})
}

func TestSlices(t *testing.T) {
	cv := initCsv()
	testCases := []struct {
		desc     string
		sliceStr string
		expect   string
	}{
		{
			desc:     "common slice",
			sliceStr: "1-4",
			expect:   MockCSV(SliceColumns(1, 4)),
		},
		{
			desc:     "start col to end of columns slice",
			sliceStr: "5-",
			expect:   MockCSV(SliceColumns(5, cv.HeadLen)),
		},
		{
			desc:     "begining to N column",
			sliceStr: "-5",
			expect:   MockCSV(SliceColumns(0, 5)),
		},
		// {
		// 	desc:     "out of bounds error",
		// 	sliceStr: "[5-20]",
		// 	expect:   `slice bounds out of range [:19] with capacity 10`,
		// },
	}
	for _, cur := range testCases {
		t.Run(cur.desc, func(t *testing.T) {
			got := captureStdout(func() {
				cv.GetColumnsBySlice(cur.sliceStr)
			})
			if got != cur.expect {
				t.Errorf("got %s, but expected %s\n", got, cur.expect)
			}
		})
	}
}

func initCsv() *Csv.Csv {
	buf := bytes.NewBuffer([]byte(strings.Join(data, "\n")))
	reader := csv.NewReader(buf)
	head, _ := reader.Read()
	cv, err := Csv.NewCsv(reader, head)
	if err != nil {
		log.Fatal("something didnt work on initalize")
	}
	return cv
}

func MockCSV(m [][]string) string {
	var buf bytes.Buffer
	wrt := csv.NewWriter(&buf)
	wrt.WriteAll(m)
	return buf.String()
}

func SliceColumns(x, y int) [][]string {
	var res [][]string
	for _, row := range data {
		cols := strings.Split(row, ",")
		res = append(res, cols[x:y])
	}
	return res
}

func captureStdout(f func()) string {
	defer func(org *os.File) {
		os.Stdout = org
	}(os.Stdout)
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	f()
	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(out)
}
