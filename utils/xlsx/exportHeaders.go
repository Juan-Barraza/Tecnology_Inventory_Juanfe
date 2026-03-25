package xlsx

var baseHeaders = []string{
	"code",
	"description",
	"historical_cost",
	"activation_date",
	"logical_status",
	"physical_status",
	"category",
	"area",
	"city",
	"responsible_name",
	"responsible_position",
	"accounting_group",
	"sub_code",
}

var auditHeaders = []string{
	"period_day",
	"period_year",
	"period_month",
	"confirmed",
	"deactivated",
	"has_label",
}

func getHeaders(exportType ExportType) []string {
	if exportType == ExportTypeAudit {
		all := make([]string, len(baseHeaders)+len(auditHeaders))
		copy(all, baseHeaders)
		copy(all[len(baseHeaders):], auditHeaders)
		return all
	}
	return baseHeaders
}
