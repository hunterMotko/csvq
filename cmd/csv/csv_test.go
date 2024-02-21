package csv_test

import (
	"encoding/csv"
	"reflect"
	"testing"

	Csv "github.com/hunterMotko/csvq/cmd/csv"
)

const data string = `
"LatD", "LatM", "LatS", "NS", "LonD", "LonM", "LonS", "EW", "City", "State"
   41,    5,   59, "N",     80,   39,    0, "W", "Youngstown", OH
   42,   52,   48, "N",     97,   23,   23, "W", "Yankton", SD
   46,   35,   59, "N",    120,   30,   36, "W", "Yakima", WA
   42,   16,   12, "N",     71,   48,    0, "W", "Worcester", MA
   43,   37,   48, "N",     89,   46,   11, "W", "Wisconsin Dells", WI
   36,    5,   59, "N",     80,   15,    0, "W", "Winston-Salem", NC
   49,   52,   48, "N",     97,    9,    0, "W", "Winnipeg", MB
   39,   11,   23, "N",     78,    9,   36, "W", "Winchester", VA
   34,   14,   24, "N",     77,   55,   11, "W", "Wilmington", NC
   39,   45,    0, "N",     75,   33,    0, "W", "Wilmington", DE
   48,    9,    0, "N",    103,   37,   12, "W", "Williston", ND
   41,   15,    0, "N",     77,    0,    0, "W", "Williamsport", PA
   37,   40,   48, "N",     82,   16,   47, "W", "Williamson", WV
   33,   54,    0, "N",     98,   29,   23, "W", "Wichita Falls", TX
   37,   41,   23, "N",     97,   20,   23, "W", "Wichita", KS
   40,    4,   11, "N",     80,   43,   12, "W", "Wheeling", WV
   26,   43,   11, "N",     80,    3,    0, "W", "West Palm Beach", FL
   47,   25,   11, "N",    120,   19,   11, "W", "Wenatchee", WA
   41,   25,   11, "N",    122,   23,   23, "W", "Weed", CA
`

func initalizeFile() *Csv.Csv {
	var reader *csv.Reader
	reader = csv.NewReader()
	return Csv.NewCsv(reader)
}

func TestCsv(t *testing.T) {
	cv := initalizeFile()
	t.Run("test files on creating new csv struct", func(t *testing.T) {
		headers := cv.GetHeaders()
		expected := []string{"LatD", "LatM", "LatS", "NS", "LonD", "LonM", "LonS", "EW", "City", "State"}
		if !reflect.DeepEqual(headers, expected) {
			t.Errorf("wanted headers %v, to be a matching array as expected %v\n", headers, expected)
		}
	})
}
