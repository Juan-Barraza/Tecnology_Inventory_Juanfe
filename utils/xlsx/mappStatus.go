package xlsx

import "inventory-juanfe/models"

var logicalStatusTranslations = map[string]string{
	"active":      "Activo",
	"inactive":    "Inactivo",
	"written_off": "Dado de baja",
}

var physicalStatusTranslations = map[string]string{
	"optimal":        "Óptimo",
	"good":           "Bueno",
	"fair":           "Regular",
	"deteriorated":   "Deteriorado",
	"out_of_service": "Fuera de servicio",
}

func translateLogicalStatus(status models.LogicalStatus) string {
	if val, ok := logicalStatusTranslations[string(status)]; ok {
		return val
	}
	return string(status)
}

func translatePhysicalStatus(status models.PhysicalStatus) string {
	if val, ok := physicalStatusTranslations[string(status)]; ok {
		return val
	}
	return string(status)
}
