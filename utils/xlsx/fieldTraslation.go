package xlsx

var inventoryFieldTranslations = map[string]string{
	"code":                 "Código",
	"description":          "Descripción",
	"historical_cost":      "Costo Histórico",
	"activation_date":      "Fecha de Activación",
	"logical_status":       "Estado Lógico",
	"physical_status":      "Estado Físico",
	"category":             "Categoría",
	"area":                 "Área",
	"city":                 "Ciudad",
	"responsible_name":     "Responsable",
	"responsible_position": "Cargo del Responsable",
	"period_day":           "Dia del periodo",
	"period_year":          "Año del Período",
	"period_month":         "Mes del Período",
	"accounting_group":     "Grupo Contable",
	"sub_code":             "Subcuenta",
	"confirmed":            "Confirmado",
	"deactivated":          "Dado de Baja",
	"has_label":            "Tiene Etiqueta",
}

func translateField(key string) string {
	if val, ok := inventoryFieldTranslations[key]; ok {
		return val
	}
	return key
}
