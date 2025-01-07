package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/utils_oas/requestresponse"
)

func TestConsultarAsignaciones(t *testing.T) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8090/v1/asignaciones/consultar/79777053", nil)
	if err != nil {
		t.Fatalf("Error al crear la solicitud GET: %v", err)
	}

	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error al ejecutar la solicitud GET: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Errorf("Se esperaba el código de estado 200, pero se obtuvo: %d", r.StatusCode)
	}

	var response requestresponse.APIResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		t.Fatalf("Error al decodificar la respuesta: %v", err)
	}

	if !response.Success {
		t.Fatal("La respuesta no indica éxito.")
	}

	if response.Data == nil {
		t.Fatal("No se recibió el campo 'Data' en la respuesta")
	}

	t.Log("Respuesta de la solicitud GET:", response)
	t.Log("Test ConsultarAsignaciones finalizado correctamente.")
}

func TestCambiarEstadoAsignacionEvaluacion(t *testing.T) {
	item := models.PeticionCambioEstadoAsignacion{
		AsignacionId:      &models.AsignacionEvaluador{Id: 41},
		AbreviacionEstado: "ER",
	}

	body, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Error al convertir a JSON: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8090/v1/asignaciones/cambiar-estado", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Error al crear la solicitud POST: %v", err)
	}

	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error al ejecutar la solicitud POST: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Errorf("Se esperaba el código de estado 200, pero se obtuvo: %d", r.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		t.Fatalf("Error al decodificar la respuesta: %v", err)
	}

	if success, ok := response["Success"].(bool); !ok || !success {
		t.Fatal("La respuesta no indica éxito.")
	}

	message, ok := response["Message"].(string)
	if !ok || message == "" {
		t.Fatal("No se recibió un mensaje válido en la respuesta.")
	}

	if response["Data"] != nil {
		t.Fatalf("Se esperaba 'Data' como null, pero se obtuvo: %v", response["Data"])
	}

	t.Log("Respuesta de la solicitud POST:", response)
	t.Log("Test CambiarEstadoAsignacionEvaluacion finalizado correctamente.")
}
