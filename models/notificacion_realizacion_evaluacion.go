package models

type NotificacionRealizacionEvaluacionEmail struct {
	Source              string                             `json:"Source"`
	Template            string                             `json:"Template"`
	Destinations        []DestinationRealizacionEvaluacion `json:"Destinations"`
	DefaultTemplateData TemplateRealizacionData            `json:"DefaultTemplateData"`
}

type DestinationRealizacionEvaluacion struct {
	Destination             Destination             `json:"Destination"`
	ReplacementTemplateData TemplateRealizacionData `json:"ReplacementTemplateData"`
	Attachments             []string                `json:"Attachments"`
}

type TemplateRealizacionData struct {
	NombreProveedor string `json:"nombre_proveedor"`
	Vigencia        string `json:"vigencia"`
	NumeroContrato  string `json:"numero_contrato"`
}
