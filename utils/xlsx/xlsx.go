package xlsx

import (
	"fmt"
	dtos "inventory-juanfe/dtos/response"
	"inventory-juanfe/models"

	"github.com/xuri/excelize/v2"
)

type ExportType string

const (
	ExportTypeGeneral ExportType = "general"
	ExportTypeAudit   ExportType = "audit"
)

func ExportDataAuditoryToXlsx(assets []models.AssetExport, exportType ExportType, countersInformation *dtos.CounterAssetsToExport) (*excelize.File, error) {
	file := excelize.NewFile()
	var sheetName string

	if exportType == ExportTypeAudit {
		sheetName = "Auditoria"
	} else {
		sheetName = "Activos"
	}
	sheetIndex, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error to create sheet")
	}
	file.DeleteSheet("Sheet1")

	fields := getHeaders(exportType)
	for i, field := range fields {
		collName, _ := excelize.ColumnNumberToName(i + 1)
		cell := collName + "1"
		headersTraduced := translateField(field)
		file.SetCellValue(sheetName, cell, headersTraduced)
	}

	if exportType == ExportTypeAudit {
		if err := buildSummarySheet(file, countersInformation); err != nil {
			return nil, fmt.Errorf("error to build summary")
		}
	}

	processRowsToXlsx(assets, file, fields, sheetName)

	file.SetActiveSheet(sheetIndex)

	return file, nil
}
