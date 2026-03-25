package handlers

import (
	"bytes"
	"inventory-juanfe/services"
	"inventory-juanfe/utils/xlsx"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ExportHandlerXlsx struct {
	service *services.ExportAssetsToXlsx
}

func NewExportHandlerXlsx(service *services.ExportAssetsToXlsx) *ExportHandlerXlsx {
	return &ExportHandlerXlsx{service: service}
}

func (h *ExportHandlerXlsx) ExportXlsx(c fiber.Ctx) error {
	exportType := c.Query("export_type")
	yearstr := c.Query("year")
	monthStr := c.Query("month")
	dayStr := c.Query("day")
	var year int
	var month int
	var day int
	var err error

	if yearstr != "" {
		year, err = strconv.Atoi(yearstr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error al parsear query",
			})
		}
	}
	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error al parsear query",
			})
		}
	}
	if dayStr != "" {
		day, err = strconv.Atoi(dayStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error al parsear query",
			})
		}
	}

	fileXlsx, err := h.service.ExportToXlsx(year, month, day, xlsx.ExportType(exportType))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var buffer bytes.Buffer
	if err := fileXlsx.Write(&buffer); err != nil {
		log.Println("Error writing file to buffer:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	downloadName := "inventrio_" + time.Now().UTC().Format("2006-01-02") + ".xlsx"
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+downloadName)
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return c.SendStream(bytes.NewReader(buffer.Bytes()))
}
