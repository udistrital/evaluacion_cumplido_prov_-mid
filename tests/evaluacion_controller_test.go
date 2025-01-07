package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCambiarEstadoEvaluacio(t *testing.T) {
	body := map[string]interface{}{
		"EvaluacionId": map[string]interface{}{
			"Id": 1,
		},
		"AbreviacionEstado": "AEV",
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Error al convertir el cuerpo a JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8090/v1/evaluacion/cambiar-estado", bytes.NewReader(bodyJSON))
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

	t.Log("Cambio de estado exitoso:", response)
}
