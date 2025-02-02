package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func TestFirmarEvaluacion(t *testing.T) {
	body := models.PeticionFirmaElectronica{
		PersonaId:    "79777053",
		AsignacionId: 41,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Error al convertir el cuerpo a JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8090/v1/firma_electronica/firmar_evaluacion", bytes.NewReader(bodyJSON))
	if err != nil {
		t.Fatalf("Error al crear la solicitud: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error al ejecutar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Se esperaba el código de estado 200, pero se obtuvo %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error al decodificar la respuesta: %v", err)
	}

	success, ok := response["Success"].(bool)
	if !ok || !success {
		t.Fatalf("La respuesta no es exitosa o no contiene el campo 'Success'.")
	}

	message, ok := response["Message"].(string)
	if !ok || message == "" {
		t.Fatalf("La respuesta no contiene el campo 'Message' o está vacío.")
	}

	t.Log("Firma de evaluación exitosa:", response)
}
