package services

import (
	"fmt"
	dtos "inventory-juanfe/dtos/response"
	"inventory-juanfe/models"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/utils/xlsx"

	"github.com/xuri/excelize/v2"
)

type ExportAssetsToXlsx struct {
	repoExpo *repository.ExporterRepository
}

func NewExportAssetsToXlsx(repoExpo *repository.ExporterRepository) *ExportAssetsToXlsx {
	return &ExportAssetsToXlsx{repoExpo: repoExpo}
}

func (s *ExportAssetsToXlsx) ExportToXlsx(year, month, day int, exportType xlsx.ExportType) (*excelize.File, error) {
	var err error
	var countersResults *dtos.CounterAssetsToExport
	var assets = []models.AssetExport{}
	if exportType == xlsx.ExportTypeAudit {
		assets, err = s.repoExpo.GetAssetsWithDate(year, month, day)
		if err != nil {
			return nil, fmt.Errorf("error to get assets")
		}
		countersResults, err = s.repoExpo.CountAssetsConfirmatedAndDesactivated(year, month, day)
		if err != nil {
			return nil, fmt.Errorf("error to get totals counters")
		}
	}
	if exportType == xlsx.ExportTypeGeneral {
		assets, err = s.repoExpo.GetAssetsToExport()
		if err != nil {
			return nil, fmt.Errorf("error to get assets in type general")
		}
	}

	file, err := xlsx.ExportDataAuditoryToXlsx(assets, exportType, countersResults)
	if err != nil {
		return nil, err
	}

	return file, nil

}
