package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestObtenerInformacionEvaluacion(t *testing.T) {

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8090/v1/informacion_evaluacion/41"), nil)
	if err != nil {
		t.Fatalf("Error al crear la solicitud: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error al ejecutar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Se esperaba el código de estado 200, pero se obtuvo: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error al decodificar la respuesta: %v", err)
	}

	success, ok := response["Success"].(bool)
	if !ok || !success {
		t.Fatalf("La respuesta no es exitosa o no contiene el campo 'Success'.")
	}

	data, ok := response["Data"].(map[string]interface{})
	if !ok || data == nil {
		t.Fatalf("La respuesta no contiene el campo 'Data' o es nula.")
	}

	nombreEvaluador, ok := data["NombreEvaluador"].(string)
	if !ok || nombreEvaluador == "" {
		t.Fatalf("No se recibió el campo 'NombreEvaluador' o está vacío.")
	}

	message, ok := response["Message"].(string)
	if !ok || message == "" {
		t.Fatalf("La respuesta no contiene el campo 'Message' o está vacío.")
	}

	t.Log("Respuesta de la solicitud GET:", response)
	t.Log("Test GET ObtenerInformacionEvaluacion finalizado correctamente.")
}
