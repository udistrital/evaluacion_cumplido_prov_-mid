package models

import "time"

type InformacionPago struct {
	Id                   int
	TipoPagoId           TipoPago
	CumplidoProveedorId  CumplidoProveedor
	TipoDocumentoCobroId int
	TipoCuentaBancariaId int
	BancoId              int
	FechaInicial         time.Time
	FechaFinal           time.Time
	NumeroFactura        string
	ValorCumplido        float64
	NumeroCuenta         string
	Activo               bool
	FechaCreacion        time.Time
	FechaModificacion    time.Time
}

type TipoPago struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}
