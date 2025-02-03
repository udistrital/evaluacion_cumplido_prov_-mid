package models

type NotificacionAsignacionEvaluacionEmail struct {
	Source              string                            `json:"Source"`
	Template            string                            `json:"Template"`
	Destinations        []DestinationAsignacionEvaluacion `json:"Destinations"`
	DefaultTemplateData TemplateAsignacionData            `json:"DefaultTemplateData"`
}

type DestinationAsignacionEvaluacion struct {
	Destination             Destination            `json:"Destination"`
	ReplacementTemplateData TemplateAsignacionData `json:"ReplacementTemplateData"`
	Attachments             []string               `json:"Attachments"`
}

type Destination struct {
	BccAddresses []string `json:"BccAddresses"`
	CcAddresses  []string `json:"CcAddresses"`
	ToAddresses  []string `json:"ToAddresses"`
}

type TemplateAsignacionData struct {
	RolEvaluador    string `json:"rol_evaluador"`
	NombreProveedor string `json:"nombre_proveedor"`
	ItemsEvaluar    string `json:"items_evaluar"`
}
