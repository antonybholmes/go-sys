package sys

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type Table struct {
	Index   [][]string `json:"index"`
	Columns [][]string `json:"columns"`
	Data    [][]string `json:"data"`
}

func XlsxToText(xlsx []byte, index int, header int) (*Table, error) {
	r := bytes.NewReader(xlsx)

	f, err := excelize.OpenReader(r) // .OpenFile("Book1.xlsx")

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

	sheets := f.GetSheetList()

	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets")
	}

	// Always pick the first sheet
	firstSheet := sheets[0]

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(firstSheet)
	if err != nil {

		return nil, err
	}

	n := len(rows[0])
	colStart := max(index+1, 0)
	//dataStart := max(header+1, 0)
	cols := len(rows) - colStart

	indexNames := make([][]string, 0, 100)
	columns := make([][]string, cols)
	data := make([][]string, 0, 100)

	for i := 0; i < cols; i++ {
		columns[i] = make([]string, header+1)
	}

	//log.Debug().Msgf("hmm:%d %d %d %d", cols, len(rows), colStart, header)

	for ri, row := range rows {
		if ri <= header {
			for i := colStart; i < n; i++ {
				columns[i-colStart][ri] = row[i]
			}
		} else {
			indexNames = append(indexNames, make([]string, colStart))

			for i := 0; i <= index; i++ {
				indexNames[len(indexNames)-1][i] = row[i]
			}

			data = append(data, row[colStart:n])

		}
	}

	ret := Table{Index: indexNames, Columns: columns, Data: data}

	return &ret, nil
}
