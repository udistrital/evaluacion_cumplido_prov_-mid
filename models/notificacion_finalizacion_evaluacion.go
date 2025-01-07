package models

type NotificacionFinalizacionEvaluacionEmail struct {
	Source              string                              `json:"Source"`
	Template            string                              `json:"Template"`
	Destinations        []DestinationFinalizacionEvaluacion `json:"Destinations"`
	DefaultTemplateData TemplateFinalizacionData            `json:"DefaultTemplateData"`
}

type DestinationFinalizacionEvaluacion struct {
	Destination             Destination              `json:"Destination"`
	ReplacementTemplateData TemplateFinalizacionData `json:"ReplacementTemplateData"`
	Attachments             []string                 `json:"Attachments"`
}

type TemplateFinalizacionData struct {
	NombreProveedor   string `json:"nombre_proveedor"`
	Vigencia          string `json:"vigencia"`
	NumeroContrato    string `json:"numero_contrato"`
	NombreEvaluadores string `json:"nombre_evaluadores"`
}
