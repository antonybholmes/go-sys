package sys

import (
	"bytes"
	"errors"
	"strings"

	"github.com/antonybholmes/go-sys/log"
	"github.com/xuri/excelize/v2"
)

type Table struct {
	IndexNames []string   `json:"indexNames"`
	Index      [][]string `json:"index"`
	Columns    [][]string `json:"columns"`
	Data       [][]string `json:"data"`
}

func XlsxSheetNames(reader *bytes.Reader) ([]string, error) {

	f, err := excelize.OpenReader(reader) // .OpenFile("Book1.xlsx")

	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the spreadsheet.
		err := f.Close()

		if err != nil {
			log.Debug().Msgf("err closing xlsx: %s", err)
		}
	}()

	// Always pick the first sheet
	sheets := f.GetSheetList()

	return sheets, nil
}

// Convert the first sheet of an excel file into a json representation
func XlsxToJson(reader *bytes.Reader,
	sheet string,
	indexes int,
	headers int,
	skipRows int,
	trimWhitespace bool) (*Table, error) {

	f, err := excelize.OpenReader(reader) // .OpenFile("Book1.xlsx")

	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the spreadsheet.
		err := f.Close()

		if err != nil {
			log.Debug().Msgf("err closing xlsx: %s", err)
		}
	}()

	// Always pick the first sheet

	if sheet == "" {
		sheet = f.GetSheetName(0)
	}

	if sheet == "" {
		return nil, errors.New("no sheets")
	}

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheet)

	if err != nil {
		return nil, err
	}

	headers = max(0, headers)
	indexes = max(0, indexes)
	skipRows = max(0, skipRows)

	// rows we don't care about
	rows = rows[skipRows:]

	colStart := indexes
	colCount := len(rows[0]) - colStart

	//dataStart := max(header+1, 0)

	columns := make([][]string, 0, colCount)

	if headers > 0 {
		for c := range colCount {
			columns = append(columns, make([]string, headers))

			for r := range headers {
				columns[c][r] = rows[r][colStart+c]

				if trimWhitespace {
					columns[c][r] = strings.TrimSpace(columns[c][r])
				}
			}
		}
	}

	indexNames := make([]string, 0, indexes)

	if headers > 0 && indexes > 0 {
		indexNames = rows[headers-1][0:indexes]
	}

	// remove headers
	rows = rows[headers:]
	rowCount := len(rows)

	//log.Debug().Msgf("xlsx: %d %d %d", len(rows), headers, indexes)

	indexData := make([][]string, 0, rowCount)

	// for i := range rowCount {
	// 	indexData[i] = make([]string, indexes)
	// }

	data := make([][]string, rowCount)

	for ri, row := range rows {
		if indexes > 0 {
			indexData = append(indexData, row[0:indexes])
		}

		data[ri] = row[indexes:]
	}

	ret := Table{
		IndexNames: indexNames,
		Index:      indexData,
		Columns:    columns,
		Data:       data}

	return &ret, nil
}
