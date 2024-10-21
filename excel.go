package sys

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type Table struct {
	IndexName string     `json:"indexName"`
	Index     [][]string `json:"index"`
	Columns   [][]string `json:"columns"`
	Data      [][]string `json:"data"`
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

func XlsxToJson(reader *bytes.Reader, sheet string, indexes int, headers int, skipRows int) (*Table, error) {

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
		return nil, fmt.Errorf("no sheets")
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
	//dataStart := max(header+1, 0)

	columns := make([][]string, headers)

	for i := 0; i < headers; i++ {
		columns[i] = rows[i][colStart:]
	}

	indexName := ""

	if headers > 0 && indexes > 0 {
		indexName = rows[headers-1][0]
	}

	// remove headers
	rows = rows[headers:]
	rowCount := len(rows)

	//log.Debug().Msgf("xlsx: %d %d %d", len(rows), headers, indexes)

	indexNames := make([][]string, indexes)
	for i := 0; i < indexes; i++ {
		indexNames[i] = make([]string, rowCount)
	}

	data := make([][]string, rowCount)

	for ri, row := range rows {

		for i := 0; i < indexes; i++ {
			indexNames[i][ri] = row[i]
		}

		data[ri] = row[indexes:]

	}

	ret := Table{IndexName: indexName, Index: indexNames, Columns: columns, Data: data}

	return &ret, nil
}
