package model

import "time"

type Medicao struct {
	ID              string
	DataMedicao     time.Time
	PotenciaAtiva   float64
	PotenciaReativa float64
	Medidor         *Medidor
}

func (medicao *Medicao) NewParams() {
	if medicao.Medidor == nil {
		medicao.Medidor = new(Medidor)
	}
}

func (medicao *Medicao) GetParams() (parametros []interface{}) {
	medicao.NewParams()
	parametros = append(parametros,
		&medicao.ID,
		&medicao.DataMedicao,
		&medicao.PotenciaAtiva,
		&medicao.PotenciaReativa,
		&medicao.Medidor.ID)
	return
}
