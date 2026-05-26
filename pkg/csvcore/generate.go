package csvcore

import (
	"encoding/csv"
	"fmt"
	"io"
)

func GenerateSynthetic(writer io.Writer, rows int, cols int) error {
	if rows < 0 {
		rows = 0
	}
	if cols < 2 {
		cols = 2
	}

	csvWriter := csv.NewWriter(writer)
	header := make([]string, cols)
	header[0] = "id"
	for col := 1; col < cols; col++ {
		header[col] = fmt.Sprintf("field_%02d", col)
	}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	record := make([]string, cols)
	for row := 0; row < rows; row++ {
		record[0] = fmt.Sprintf("%d", row+1)
		for col := 1; col < cols; col++ {
			switch col % 4 {
			case 0:
				record[col] = fmt.Sprintf("%.2f", float64((row+1)*(col+3))/7.0)
			case 1:
				record[col] = fmt.Sprintf("segment_%d", row%12)
			case 2:
				record[col] = fmt.Sprintf("2026-05-%02d", row%28+1)
			default:
				record[col] = fmt.Sprintf("note row %d col %d", row+1, col+1)
			}
		}
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
