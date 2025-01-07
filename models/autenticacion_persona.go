package models

type AutenticacionPersona struct {
	Rol                []string `json:"role"`
	Documento          string   `json:"documento"`
	DocumentoCompuesto string   `json:"documento_compuesto"`
	Email              string   `json:"email"`
	FamilyName         string   `json:"FamilyName"`
	Codigo             string   `json:"codigo"`
	Estado             string   `json:"estado"`
}
