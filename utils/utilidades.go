package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
)

func GetJsonWSO2Test(urlp string, target interface{}) (status int, err error) {
	b := new(bytes.Buffer)
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(nil, err)
		}
	}()
	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func GetJsonTest(url string, target interface{}) (status int, err error) {
	r, err := http.Get(url)
	if err != nil {
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	fmt.Println("JSON recibido:", r.StatusCode)
	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func SendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func SendJsonTls(url string, trequest string, target interface{}, datajson interface{}) error {
	// Crear un buffer para el cuerpo de la solicitud
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}

	// Configurar el transporte del cliente para que ignore las validaciones de certificado
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Crear la solicitud HTTP
	req, err := http.NewRequest(trequest, url, b)
	if err != nil {
		beego.Error("error creando la solicitud", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	// Decodificar la respuesta en la variable target
	return json.NewDecoder(r.Body).Decode(target)
}

func ValorLetras(n int) string {
	var unidades = []string{"", "uno", "dos", "tres", "cuatro", "cinco", "seis", "siete", "ocho", "nueve", "diez", "once", "doce", "trece", "catorce", "quince", "dieciséis", "diecisiete", "dieciocho", "diecinueve", "veinte"}
	var decenas = []string{"", "", "veinti", "treinta", "cuarenta", "cincuenta", "sesenta", "setenta", "ochenta", "noventa"}
	var centenas = []string{"", "cien", "doscientos", "trescientos", "cuatrocientos", "quinientos", "seiscientos", "setecientos", "ochocientos", "novecientos"}

	if n == 0 {
		return "cero"
	}

	// Millones
	if n >= 1000000 {
		if n == 1000000 {
			return "un millón"
		}
		millones := fmt.Sprintf("%s millones", ValorLetras(n/1000000))
		resto := n % 1000000
		if resto == 0 {
			return millones
		}
		return fmt.Sprintf("%s %s", millones, ValorLetras(resto))
	}

	// Miles
	if n >= 1000 {
		if n == 1000 {
			return "mil"
		}
		miles := fmt.Sprintf("%s mil", ValorLetras(n/1000))
		resto := n % 1000
		// Si el resto es 0, solo retornamos la parte de los miles.
		if resto == 0 {
			return miles
		}
		return fmt.Sprintf("%s %s", miles, ValorLetras(resto))
	}

	// Centenas
	if n >= 100 {
		if n == 100 {
			return "cien"
		}
		centenasStr := fmt.Sprintf("%s %s", centenas[n/100], ValorLetras(n%100))
		resto := n % 100
		if resto == 0 {
			return strings.TrimSpace(centenas[n/100])
		}
		return strings.TrimSpace(centenasStr)
	}

	// Decenas
	if n >= 20 {
		if n%10 == 0 {
			return decenas[n/10]
		}
		if n < 30 {
			return decenas[n/10] + unidades[n%10]
		}
		return fmt.Sprintf("%s y %s", decenas[n/10], unidades[n%10])
	}

	// Unidades
	return unidades[n]
}

func FormatNumber(value interface{}, precision int, thousand string, decimal string) string {
	v := reflect.ValueOf(value)
	var x string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = fmt.Sprintf("%d", v.Int())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x = fmt.Sprintf("%d", v.Uint())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Float32, reflect.Float64:
		x = fmt.Sprintf(fmt.Sprintf("%%.%df", precision), v.Float())
	case reflect.Ptr:
		switch v.Type().String() {
		case "*big.Rat":
			x = value.(*big.Rat).FloatString(precision)

		default:
			panic("Unsupported type - " + v.Type().String())
		}
	default:
		panic("Unsupported type - " + v.Kind().String())
	}

	return FormatNumberString(x, precision, thousand, decimal)
}

func FormatNumberString(x string, precision int, thousand string, decimal string) string {
	lastIndex := strings.Index(x, ".") - 1
	if lastIndex < 0 {
		lastIndex = len(x) - 1
	}

	var buffer []byte
	var strBuffer bytes.Buffer

	j := 0
	for i := lastIndex; i >= 0; i-- {
		j++
		buffer = append(buffer, x[i])

		if j == 3 && i > 0 && !(i == 1 && x[0] == '-') {
			buffer = append(buffer, ',')
			j = 0
		}
	}

	for i := len(buffer) - 1; i >= 0; i-- {
		strBuffer.WriteByte(buffer[i])
	}
	result := strBuffer.String()

	if thousand != "," {
		result = strings.Replace(result, ",", thousand, -1)
	}

	extra := x[lastIndex+1:]
	if decimal != "." {
		extra = strings.Replace(extra, ".", decimal, 1)
	}

	return result + extra
}
