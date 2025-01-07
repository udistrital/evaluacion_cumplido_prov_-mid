package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGuardarResultadoEvaluacion(t *testing.T) {

	body := `{
		"AsignacionEvaluadorId": 4,
		"ClasificacionId": 1,
		"ResultadoEvaluacion": {
			"ResultadosIndividuales": [
				{
					"Categoria": "CUMPLIMIENTO",
					"Titulo": "TIEMPOS DE ENTREGA",
					"Respuesta": {
						"Pregunta": "¿Se cumplieron los tiempos de entrega de bienes o la prestación del servicios ofertados por el proveedor?",
						"Cumplimiento": "Si",
						"ValorAsignado": 12
					}
				},
				{
					"Categoria": "CUMPLIMIENTO",
					"Titulo": "CANTIDADES",
					"Respuesta": {
						"Pregunta": "¿Se entregan las cantidades solicitadas?",
						"Cumplimiento": "No",
						"ValorAsignado": 0
					}
				},
				{
					"Categoria": "CALIDAD",
					"Titulo": "CONFORMIDAD",
					"Respuesta": {
						"Pregunta": "¿El bien o servicio cumplió con las especificaciones y requisitos pactados en el momento de entrega?",
						"Cumplimiento": "No",
						"ValorAsignado": 0
					}
				},
				{
					"Categoria": "CALIDAD",
					"Titulo": "FUNCIONALIDAD ADICIONAL",
					"Respuesta": {
						"Pregunta": "¿El producto comprado o el servicio prestado proporcionó más herramientas o funciones de las solicitadas originalmente?",
						"Cumplimiento": "No",
						"ValorAsignado": 0
					}
				},
				{
					"Categoria": "POS CONTRACTUAL",
					"Titulo": "RECLAMACIONES",
					"Respuesta": {
						"Pregunta": "¿Se han presentado reclamaciones al proveedor en calidad o gestión?",
						"Cumplimiento": "No",
						"ValorAsignado": 12
					}
				},
				{
					"Categoria": "POS CONTRACTUAL",
					"Titulo": "RECLAMACIONES",
					"Respuesta": {
						"Pregunta": "¿El proveedor soluciona oportunamente las no conformidades de calidad y gestión de los bienes o servicios recibidos?",
						"Cumplimiento": "No",
						"ValorAsignado": 0
					}
				},
				{
					"Categoria": "POS CONTRACTUAL",
					"Titulo": "SERVICIO POS VENTA",
					"Respuesta": {
						"Pregunta": "¿El proveedor cumple con los compromisos pactados dentro del contrato u orden de servicio o compra? (aplicación de garantías, mantenimiento, cambios, reparaciones, capacitaciones, entre otras)",
						"Cumplimiento": "Si",
						"ValorAsignado": 10
					}
				},
				{
					"Categoria": "GESTIÓN",
					"Titulo": "PROCEDIMIENTOS",
					"Respuesta": {
						"Pregunta": "¿El contrato es suscrito en el tiempo pactado, entrega las pólizas a tiempo y las facturas son radicadas en el tiempo indicado con las condiciones y soportes requeridos para su trámite contractual?",
						"Cumplimiento": "Excelente",
						"ValorAsignado": 9
					}
				},
				{
					"Categoria": "GESTIÓN",
					"Titulo": "GARANTÍA",
					"Respuesta": {
						"Pregunta": "¿Se requirió hacer uso de la garantía del producto o servicio?",
						"Cumplimiento": "No",
						"ValorAsignado": 15
					}
				},
				{
					"Categoria": "GESTIÓN",
					"Titulo": "GARANTÍA",
					"Respuesta": {
						"Pregunta": "¿El proveedor cumplió a satisfacción con la garantía pactada?",
						"Cumplimiento": "No",
						"ValorAsignado": 0
					}
				}
			]
		}
	}`

	req, err := http.NewRequest("POST", "http://localhost:8090/v1/resultado/resultado-evaluacion", bytes.NewBuffer([]byte(body)))
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

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Se esperaba el código de estado 201, pero se obtuvo: %d", resp.StatusCode)
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
	if !ok || message != "Recurso creado con éxito" {
		t.Fatalf("Se esperaba el mensaje 'Recurso creado con éxito', pero se obtuvo: %s", message)
	}

	t.Log("Respuesta de la solicitud POST:", response)
	t.Log("Test POST GuardarResultadoEvaluacion finalizado correctamente.")
}
