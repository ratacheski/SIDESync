package model

type MedidorEnel struct {
	ID               string
	UC               string
	Nome             string
	NumMedidor       string
	Demanda          int
	DemandaPonta     int
	DemandaForaPonta int
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
