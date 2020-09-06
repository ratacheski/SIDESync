package logger

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//Config - Variavel global  com as informações do arquivo de configuração
var Config Configuracao

//Configuracao - Estrutura do arquivo de configuração do logger
type Configuracao struct {
	Directory      string        `json:"Diretorio"`
	Filename       string        `json:"NomeArquivo"`
	MaxSize        int           `json:"TamanhoMaximo"`
	RotatePeriod   time.Duration `json:"PeriodoRotacionar"`
	RotateByPeriod bool          `json:"RotacionarPorPeriodo"`
	LocalTime      bool          `json:"HorarioLocal"`
	LogLvl         string        `json:"NivelLog"`
	Formatter      formatter     `json:"Formato"`
}

func (cfg Configuracao) NomeArquivo() string {
	return Config.Directory + Config.Filename
}

type formatter struct {
	TimestampFormat  string `json:"DataHora"`
	ForceFormatting  bool   `json:"Formatar"`
	FullTimestamp    bool   `json:"DataHoraCompleto"`
	DisableSorting   bool   `json:"Ordenar"`
	EnableCallerFunc bool   `json:"RegistrarFuncao"`
}

func (cfg *Configuracao) ObtemLogLevel() log.Level {
	lvl, err := log.ParseLevel(strings.ToLower(cfg.LogLvl))
	if err != nil {
		log.Error("Erro ao obter o nivel de log: ", err)
		log.Info("Atualizando o nivel de log para Debug")
		lvl = log.DebugLevel
	}
	return lvl
}

//CarregaConfiguracoes - Inicia leitura do arquivo de configuração
func CarregaConfiguracoes(filename string) {
	var err error
	Config, err = readConfig(filename)
	if err != nil {
		log.Fatalln("Erro ao carregar as configurações de log: ", err)
	}
}

func ConfiguracaoPadrao() *Logger {
	log.SetFormatter(&TextFormatter{
		ForceFormatting:  Config.Formatter.ForceFormatting,
		FullTimestamp:    Config.Formatter.FullTimestamp,
		TimestampFormat:  Config.Formatter.TimestampFormat,
		DisableSorting:   Config.Formatter.DisableSorting,
		EnableCallerFunc: Config.Formatter.EnableCallerFunc,
	})
	// Habilita o log da linha e do nome do arquivo
	log.SetReportCaller(true)
	defaultLogger := &Logger{
		Filename:     Config.NomeArquivo(),
		MaxSize:      Config.MaxSize,      //MB
		RotatePeriod: Config.RotatePeriod, //Hours
		LocalTime:    Config.LocalTime,
	}
	//Rotaciona o arquivo de log
	if Config.RotateByPeriod {
		defaultLogger.RotatePeriodically()
	}
	return defaultLogger
}

//lê as informações do arquivo de configuração
func readConfig(fileName string) (Configuracao, error) {
	config := Configuracao{}
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return config, err
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		log.Println(err)
		return config, err
	}
	return config, err
}
