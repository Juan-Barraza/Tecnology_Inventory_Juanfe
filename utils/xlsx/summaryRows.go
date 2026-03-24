package xlsx

import (
	"fmt"
	dtos "inventory-juanfe/dtos/response"

	"github.com/xuri/excelize/v2"
)

var summaryRows = []struct {
	Key   string
	Label string
}{
	{"total_confirmated", "Total Confirmados"},
	{"total_desactivated", "Total Dados de Baja"},
	{"total_with_label", "Total Con Etiqueta"},
	{"total_without_label", "Total Sin Etiqueta"},
}

func buildSummarySheet(file *excelize.File, counters *dtos.CounterAssetsToExport) error {
	const sheetName = "Resumen"

	file.NewSheet(sheetName)

	// Header de la hoja
	file.SetCellValue(sheetName, "A1", "Concepto")
	file.SetCellValue(sheetName, "B1", "Total")

	// Valores de los contadores en orden definido
	values := map[string]int64{
		"total_confirmated":   counters.TotalConfirmated,
		"total_desactivated":  counters.TotalDesactivated,
		"total_with_label":    counters.TotalWithLabel,
		"total_without_label": counters.TotalWithoutLabel,
	}

	for i, row := range summaryRows {
		rowNum := i + 2 // empieza en fila 2 porque la 1 es el header
		cellA := fmt.Sprintf("A%d", rowNum)
		cellB := fmt.Sprintf("B%d", rowNum)
		file.SetCellValue(sheetName, cellA, row.Label)
		file.SetCellValue(sheetName, cellB, values[row.Key])
	}

	return nil
}
