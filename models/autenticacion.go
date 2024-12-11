package models

type Autentiacion struct {
	Role               []string `json:"role"`
	Documento          string   `json:"documento"`
	DocumentoCompuesto string   `json:"documento_compuesto"`
	Email              string   `json:"email"`
	FamilyName         string   `json:"FamilyName"`
	Codigo             string   `json:"Codigo"`
	Estado             string   `json:"Estado"`
}

type PeticionAutenticacion struct {
	Numero string `json:"numero"`
	Rol    string `json:"rol"`
	User   string `json:"user"`
}
