package model

import "database/sql"

type MedidorEnel struct {
	ID               sql.NullString
	UC               sql.NullString
	Nome             sql.NullString
	NumMedidor       sql.NullString
	Demanda          sql.NullInt64
	DemandaPonta     sql.NullInt64
	DemandaForaPonta sql.NullInt64
	Endereco         *Endereco
}

func (medidorEnel *MedidorEnel) NewParams() {
	if medidorEnel.Endereco == nil {
		medidorEnel.Endereco = new(Endereco)
	}
}

func (medidorEnel *MedidorEnel) GetParams() (parametro []interface{}) {
	medidorEnel.NewParams()
	parametro = append(parametro,
		&medidorEnel.ID,
		&medidorEnel.UC,
		&medidorEnel.Nome,
		&medidorEnel.NumMedidor,
		&medidorEnel.Demanda,
		&medidorEnel.DemandaPonta,
		&medidorEnel.DemandaForaPonta,
		&medidorEnel.Endereco.ID)
	return
}
