package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestCambiarRolAsignacionEvaluador(t *testing.T) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8090/v1/cambio_rol_evaluador/1"), nil)
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
	t.Log("Test CambiarRolAsignacionEvaluador finalizado correctamente.")
}
