package xlsx

import (
	"fmt"
	"inventory-juanfe/models"

	"github.com/xuri/excelize/v2"
)

func ElectionFieldsAsset(asset models.AssetExport, cell, sheetName, field string, file *excelize.File) {
	switch field {
	case "code":
		file.SetCellValue(sheetName, cell, asset.Code)

	case "description":
		file.SetCellValue(sheetName, cell, asset.Description)

	case "historical_cost":
		if asset.HistoricalCost != nil {
			file.SetCellValue(sheetName, cell, *asset.HistoricalCost)
		}

	case "activation_date":
		file.SetCellValue(sheetName, cell, asset.ActivationDate.Format("2006-01-02"))

	case "logical_status":
		file.SetCellValue(sheetName, cell, translateLogicalStatus(asset.LogicalStatus))

	case "physical_status":
		file.SetCellValue(sheetName, cell, translatePhysicalStatus(asset.PhysicalStatus))

	case "category":
		file.SetCellValue(sheetName, cell, asset.CategoryName)

	case "area":
		file.SetCellValue(sheetName, cell, asset.AreaName)

	case "city":
		file.SetCellValue(sheetName, cell, asset.CityName)

	case "responsible_name":
		file.SetCellValue(sheetName, cell, asset.ResponsibleName)

	case "responsible_position":
		file.SetCellValue(sheetName, cell, asset.ResponsiblePosition)

	case "accounting_group":
		file.SetCellValue(sheetName, cell, fmt.Sprintf("%d", asset.AccountCodeGroup))

	case "sub_code":
		file.SetCellValue(sheetName, cell, fmt.Sprintf("%d", asset.SubCode))

	// ── Campos de auditoría — solo presentes en ExportTypeAudit ──
	case "period_year":
		if asset.PeriodYear != nil {
			file.SetCellValue(sheetName, cell, *asset.PeriodYear)
		}

	case "period_month":
		if asset.PeriodMonth != nil {
			file.SetCellValue(sheetName, cell, *asset.PeriodMonth)
		}

	case "confirmed":
		if asset.Confirmed != nil {
			if *asset.Confirmed {
				file.SetCellValue(sheetName, cell, "Sí")
			} else {
				file.SetCellValue(sheetName, cell, "No")
			}
		}

	case "deactivated":
		if asset.Deactivated != nil {
			if *asset.Deactivated {
				file.SetCellValue(sheetName, cell, "Sí")
			} else {
				file.SetCellValue(sheetName, cell, "No")
			}
		}

	case "has_label":
		if asset.HasLabel != nil {
			if *asset.HasLabel {
				file.SetCellValue(sheetName, cell, "Sí")
			} else {
				file.SetCellValue(sheetName, cell, "No")
			}
		}
	}
}
