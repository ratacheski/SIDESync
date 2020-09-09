package service

import (
	"encoding/xml"
	"errors"
	log "github.com/sirupsen/logrus"
	"sidegosynchronizer/core/config"
	"sidegosynchronizer/core/resty"
	"sidegosynchronizer/data/dto"
	"sidegosynchronizer/data/model"
	"sidegosynchronizer/data/repository"
	"strconv"
	"sync"
)

//SincronizaMedidoresCCK Sincroniza os medidores do banco com os medidores cadastrados na base da cck
func SincronizaMedidoresCCK() (medidoresToReturn []model.Medidor, err error) {
	log.Println("Listando Medidores No WebService")
	medidoresCCK, err := listaMedidoresCCK()
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(medidoresCCK.Medidores))
	for i, medidor := range medidoresCCK.Medidores {
		go func(i int, medidor dto.Medidor) {
			defer wg.Done()
			cck, err := listaDatasMedicoesMedidorCCK(medidor)
			if err != nil {
				log.Error(`Erro ao buscar Detalhes do Medidor `+
					medidor.Nome+` no CCK WebService. `, err)
			}
			medidoresCCK.Medidores[i] = cck
		}(i, medidor)
	}
	wg.Wait()
	medidoresToReturn, err = insereMedidoresNoBanco(medidoresCCK)
	if err != nil {
		log.Error("Erro ao inserir medidores no Banco: ", err)
	}
	return
}

//listaMedidoresCCK Lista os medidores cadastrados na api da CCK
func listaMedidoresCCK() (medidores dto.Medidores, err error) {
	var urlListagem = config.LoadedConfigs.URLConexaoCCK + "?id_medidor=?"
	resp, err := resty.Client.R().Get(urlListagem)
	if err != nil {
		log.Error("Erro ao buscar Medidores CCK: ", err)
		return
	}
	body := resp.Body()
	medidores = dto.Medidores{}
	err = xml.Unmarshal(body, &medidores)
	if err != nil {
		return
	}
	return
}

//listaDatasMedicoesMedidorCCK lista as datas de primeira e última leitura do medidor na api da CCK
func listaDatasMedicoesMedidorCCK(medidor dto.Medidor) (medidorRetorno dto.Medidor, err error) {
	var urlListagem = config.LoadedConfigs.URLConexaoCCK +
		`?id_medidor=` + medidor.ID
	resp, err := resty.Client.R().Get(urlListagem)
	if err != nil {
		log.Error("Erro ao buscar Datas Medicoes Medidor CCK: ", err)
		return
	}
	body := resp.Body()
	medidorRetorno = dto.Medidor{}
	err = xml.Unmarshal(body, &medidorRetorno)
	if err != nil {
		medidorRetorno.NaoConsultar = true
		return medidorRetorno, errors.New("medidor sem arquivo de medição")
	}
	return
}

//insereMedidoresNoBanco serviço que trata os medidores e chama o repository responsável
//pela inserção dos medidores no banco
func insereMedidoresNoBanco(medidores dto.Medidores) (meds []model.Medidor, err error) {
	meds = make([]model.Medidor, 0)
	for _, medidor := range medidores.Medidores {
		if !medidor.NaoConsultar {
			idMedidor, _ := strconv.Atoi(medidor.ID)
			var med model.Medidor
			med.ID = idMedidor
			med.Denominacao = medidor.Nome
			med.DataPrimeiraLeitura.Time = medidor.PrimeiraLeitura.Time
			med.DataPrimeiraLeitura.Valid = true
			med.DataUltimaLeitura.Time = medidor.UltimaLeitura.Time
			med.DataUltimaLeitura.Valid = true
			meds = append(meds, med)
		}
	}
	err = repository.InsereMedidores(meds)
	return
}
