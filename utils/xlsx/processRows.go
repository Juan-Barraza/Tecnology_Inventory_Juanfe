package xlsx

import (
	"fmt"
	"inventory-juanfe/models"

	"github.com/xuri/excelize/v2"
)

func processRowsToXlsx(assests []models.AssetExport, file *excelize.File, selectedFields []string, sheetName string) {
	for rowIdx, assests := range assests {
		row := rowIdx + 2 // iniciar en la segunda fila
		for colIdx, field := range selectedFields {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			cell := colName + fmt.Sprintf("%d", row)
			ElectionFieldsAsset(assests, cell, sheetName, field, file)
		}
	}
}
