package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sidegosynchronizer/core/config"
	"sidegosynchronizer/core/resty"
	"sidegosynchronizer/core/storage"
	"sidegosynchronizer/data/model"
	"sidegosynchronizer/data/service"
	"sidegosynchronizer/logger"
	"sync"
	"time"
)

// Entry point application.
func main() {
	startConfigs()
	for {
		log.Print("Sincronização Iniciada.")
		if err := sincroniza(); err != nil {
			break
		}
		log.Print("Sincronização Finalizada. Tempo de espera acionado.")
		time.Sleep(time.Duration(config.LoadedConfigs.IntervaloSincronizacao) * time.Second)
	}
	log.Error("Erro no sincronizador.")
}

func startConfigs() {
	var cfg config.Config
	if err := cfg.LoadFromJSON("config.json"); err != nil {
		log.Fatal(err)
	}
	logger.CarregaConfiguracoes(config.LoadedConfigs.LogConfig)
	defaultLogger := logger.ConfiguracaoPadrao()
	log.SetLevel(logger.Config.ObtemLogLevel())
	mw := io.MultiWriter(os.Stdout, defaultLogger)
	log.SetOutput(mw)
	resty.SetupResty()
}

func sincroniza() (err error) {
	err = storage.New(config.LoadedConfigs)
	if err != nil {
		log.Fatal(err)
		return
	}
	medidores, err := service.SincronizaMedidoresCCK()
	if err == nil {
		var wg sync.WaitGroup
		wg.Add(len(medidores))
		for i, medidor := range medidores {
			go func(i int, medidor model.Medidor) {
				defer wg.Done()
				service.SincronizaMedicoesCCK(medidor)
			}(i, medidor)
		}
		wg.Wait()
	}
	return nil
}
