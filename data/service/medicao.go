package service

import (
	"database/sql"
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"sidegosynchronizer/core/config"
	"sidegosynchronizer/core/resty"
	"sidegosynchronizer/core/utils"
	"sidegosynchronizer/data/dto"
	"sidegosynchronizer/data/model"
	"sidegosynchronizer/data/repository"
	"strconv"
	"time"
)

func SincronizaMedicoesCCK(medidor model.Medidor) {
	medidorBanco, err := repository.BuscaMedidorBanco(medidor.ID)
	if err != nil {
		log.Error("Erro ao buscar detalhes do medidor no banco: ",err)
		return
	}
	medicao, err := repository.BuscaUltimaLeituraMedidorNoBanco(medidorBanco)
	var leituras []dto.Leitura
	if err == sql.ErrNoRows {
		leituras, err = listaPrimeirasLeiturasCCK(medidorBanco)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	} else {
		leituras, err = listaProximasLeiturasCCK(medicao, medidorBanco)
		if err != nil {
			return
		}
	}
	if len(leituras) > 0 {
		err := repository.InsereMedicoes(leituras, medidorBanco)
		if err != nil {
			log.Error("Erro ao inserir medições do medidor ", medidorBanco.Denominacao, " no banco: ", err)
		}
	}
}

func listaProximasLeiturasCCK(medicao model.Medicao, medidor model.Medidor) (leituras []dto.Leitura, err error) {

	if medidor.DataUltimaLeitura.Time == medicao.DataMedicao {
		return
	}
	if medidor.UltimaDataSincronizada.Time.Add(14*time.Minute).After(medidor.DataUltimaLeitura.Time) {
		return
	}
	var urlListagem = config.LoadedConfigs.URLConexaoCCK + "?id_medidor=" +
		strconv.Itoa(medicao.Medidor.ID) + "&datahora_ini="
	if medidor.UltimaDataSincronizada.Time.After(medicao.DataMedicao) {
		urlListagem += utils.FormataDataHoraWebService(medidor.UltimaDataSincronizada.Time.Add(time.Minute*14))
	} else {
		urlListagem += utils.FormataDataHoraWebService(medicao.DataMedicao.Add(time.Minute*14))
	}
	resp, err := resty.Client.R().Get(urlListagem)
	if err != nil {
		log.Error("Erro ao buscar Próximas Leituras Medidor CCK: ", err)
		return
	}
	body := resp.Body()
	telemetria := dto.Telemetria{}
	err = xml.Unmarshal(body, &telemetria)
	if err != nil {
		log.Error("Erro no xml retornado: ", urlListagem, " ", err)
		return
	}
	return telemetria.Medidor.Leituras, nil
}

func listaPrimeirasLeiturasCCK(medidor model.Medidor) (leituras []dto.Leitura, err error) {
	if medidor.UltimaDataSincronizada.Time.Add(14*time.Minute).After(medidor.DataUltimaLeitura.Time) {
		return
	}
	var urlListagem = config.LoadedConfigs.URLConexaoCCK + "?id_medidor=" +
		strconv.Itoa(medidor.ID) + "&datahora_ini="
	if medidor.UltimaDataSincronizada.Time.After(medidor.DataPrimeiraLeitura.Time) {
		urlListagem += utils.FormataDataHoraWebService(medidor.UltimaDataSincronizada.Time.Add(time.Minute*14))
	} else {
		urlListagem += utils.FormataDataHoraWebService(medidor.DataPrimeiraLeitura.Time.Add(time.Minute*14))
	}
	resp, err := resty.Client.R().Get(urlListagem)
	if err != nil {
		log.Error("Erro ao buscar Primeiras Leituras Medidor CCK: ", err)
		return
	}
	body := resp.Body()
	telemetria := dto.Telemetria{}
	err = xml.Unmarshal(body, &telemetria)
	if err != nil {
		log.Error("Erro no xml retornado: ", urlListagem, " ", err)
		return
	}
	return telemetria.Medidor.Leituras, nil
}
