package storage

import (
	"database/sql"
	"fmt"
	"sidegosynchronizer/core/config"
	"strings"
)

var DB *sql.DB

// New cria uma instância do objeto de acesso e manipulação da base de dados.
func New(config *config.Config) (err error) {
	if DB, err = sql.Open("postgres", geraStringConexao(config)); err != nil {
		return
	}
	DB.SetMaxIdleConns(-1)
	if err = DB.Ping(); err != nil {
		return
	}
	return
}

// geraStringConexao gera a string de conexão com o banco de dados de acordo com o 'objeto' Config responável pela chamada do método
func geraStringConexao(config *config.Config) string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName)
	return psqlInfo
}

// AdicionaAlias - Adiciona um alias nos campos passados (id_usuario > u.id_usuario)
func AdicionaAlias(campos string, alias string) string {
	arrCampos := strings.Split(campos, ",")
	var resCampos []string
	for _, campo := range arrCampos {
		resCampos = append(resCampos, alias+"."+strings.TrimSpace(campo))
	}
	return strings.Join(resCampos, ",")
}
