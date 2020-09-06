package model

import "time"

type Medidor struct {
	ID                  int
	Denominacao         string
	Tipo                int
	Latitude            float64
	Longitude           float64
	DataPrimeiraLeitura time.Time
	DataUltimaLeitura   time.Time
	MedidorEnel         *MedidorEnel
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
		&medidor.MedidorEnel.ID)
	return
}
