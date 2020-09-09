package repository

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"sidegosynchronizer/core/storage"
	"sidegosynchronizer/data/model"
	"sidegosynchronizer/data/query"
)

//InsereMedidores Insere os medidores no banco
func InsereMedidores(registros []model.Medidor) (err error) {

	tx, err := storage.DB.Begin()
	if err != nil {
		log.Error("Erro ao iniciar transaction: ", err)
		return
	}
	for _, row := range registros {
		_, err = tx.Exec(query.UpsertMedidores,
			&row.ID, &row.Denominacao,&row.DataPrimeiraLeitura,&row.DataUltimaLeitura)
		if err != nil {
			log.Error("Erro no exec: ", err)
			tx.Rollback()
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Error("Erro ao comittar. ", err)
	}
	return
}

func BuscaMedidorBanco(id int) (registro model.Medidor, err error) {
	row := storage.DB.QueryRow(query.BuscaMedidor, id)
	err = row.Scan(registro.GetParams()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return  registro, nil
		}
		log.Error("Erro no scan: ", err)
	}
	return
}
