package repository

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sidegosynchronizer/core/storage"
	"sidegosynchronizer/data/model"
	"sidegosynchronizer/data/query"
	"strings"
)

//InsereMedidores Insere os medidores no banco
func InsereMedidores(registros []model.Medidor) (err error) {
	sql := query.UpsertMedidoresPrefix
	var vals []interface{}
	var inserts []string
	for i, row := range registros {
		inserts = append(inserts, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		vals = append(vals, row.ID,
			row.Denominacao,
			row.DataPrimeiraLeitura,
			row.DataUltimaLeitura)
	}
	sql += strings.Join(inserts, ",")
	sql = sql + query.UpsertMedidoresSufix
	_, err = storage.DB.Exec(sql, vals...)
	if err != nil {
		log.Error("Erro no exec: ", err)
		return
	}
	return
}
