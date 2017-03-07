package report

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"io"
	"os"
	"sync"
)

func WriteCsv(f io.Writer, report chan ReportRow, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()
	w := csv.NewWriter(f)
	defer w.Flush()

	for r := range report {
		switch row := r.(type) {
		case ReportHeader:
			w.Write(row.Headers)

		case ReportData:
			rowStr := make([]string, 0)
			for _, a := range row.Data {
				rowStr = append(rowStr, fmt.Sprintf("%v", a))
			}
			w.Write(rowStr)

		case ReportEOF:
			return nil

		default:
			seelog.Warnf("Unexpected row")
			return errors.New("Unexpected row detected")
		}
	}

	return nil
}

func WriteCsvFile(file string, report chan ReportRow, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	f, err := os.Create(file)
	if err != nil {
		seelog.Errorf("Unable to write file[%s] erorr[%s]", file, err)
		return err
	}
	defer f.Close()
	return WriteCsv(f, report, wg)
}