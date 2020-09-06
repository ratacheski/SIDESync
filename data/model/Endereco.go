package model

type Endereco struct {
	ID          string
	CEP         string
	Logradouro  string
	Numero      int
	Complemento string
	Bairro      string
	Cidade      string
	Estado      *string
}

func (endereco *Endereco) GetParams() (parametros []interface{}) {
	parametros = append(parametros,
		&endereco.ID,
		&endereco.CEP,
		&endereco.Logradouro,
		&endereco.Numero,
		&endereco.Complemento,
		&endereco.Bairro,
		&endereco.Cidade,
		&endereco.Estado)
	return
}
