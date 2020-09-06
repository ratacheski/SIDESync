package repository

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sidegosynchronizer/core/storage"
	"sidegosynchronizer/data/dto"
	"sidegosynchronizer/data/model"
	"sidegosynchronizer/data/query"
	"strings"
)

func BuscaUltimaLeituraMedidorNoBanco(medidor model.Medidor) (medicao model.Medicao, err error) {
	row := storage.DB.QueryRow(query.BuscaUltimaLeituraMedidor, medidor.ID)
	err = row.Scan(medicao.GetParams()...)
	return
}

func InsereMedicoes(registros []dto.Leitura, medidor model.Medidor) (err error) {
	sql := query.InsertMedicoesPrefix
	var vals []interface{}
	var inserts []string
	j := 0
	for _, row := range registros {
		if row.Ativa != 0 || row.Reativa != 0 {
			inserts = append(inserts, fmt.Sprintf("(newid(), $%d, $%d, $%d, $%d)",
				j*4+1, j*4+2, j*4+3, j*4+4))
			vals = append(vals,
				row.DataHora.Time,
				row.Ativa,
				row.Reativa,
				medidor.ID)
			j++
		}
	}
	if len(inserts) > 0 {
		sql += strings.Join(inserts, ",")
		ret, err := storage.DB.Exec(sql, vals...)
		if err != nil {
			log.Error("Erro no exec: ", err)
			return err
		}
		affected, _ := ret.RowsAffected()
		log.Info(medidor.ID, " - ", medidor.Denominacao, " - ", affected, " Medições Inseridas")
	} else {
		log.Info(medidor.ID, " - ", medidor.Denominacao, " - Sem Novas Medições Relevantes")
	}
	return
}
