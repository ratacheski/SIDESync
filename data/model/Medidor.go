package model

import (
	"database/sql"
)

type Medidor struct {
	ID                     int
	Denominacao            string
	Tipo                   sql.NullInt64
	Latitude               sql.NullFloat64
	Longitude              sql.NullFloat64
	DataPrimeiraLeitura    sql.NullTime
	DataUltimaLeitura      sql.NullTime
	UltimaDataSincronizada sql.NullTime
	MedidorEnel            *MedidorEnel
}

func (medidor *Medidor) NewParams() {
	if medidor.MedidorEnel == nil {
		medidor.MedidorEnel = new(MedidorEnel)
	}
}

func (medidor *Medidor) GetParams() (parametros []interface{}) {
	medidor.NewParams()
	parametros = append(parametros,
		&medidor.ID,
		&medidor.Denominacao,
		&medidor.Tipo,
		&medidor.Latitude,
		&medidor.Longitude,
		&medidor.DataPrimeiraLeitura,
		&medidor.DataUltimaLeitura,
		&medidor.UltimaDataSincronizada,
		&medidor.MedidorEnel.ID)
	return
}
