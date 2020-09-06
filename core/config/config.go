package config

import (
	"encoding/json"
	"os"
)

var LoadedConfigs *Config

// Config contém as informações de conexãocom o banco de dados
type Config struct {
	Host                   string `json:"host"`
	Password               string `json:"password"`
	User                   string `json:"user"`
	DbName                 string `json:"dbName"`
	Port                   string `json:"port"`
	LogConfig              string `json:"logConfig"`
	URLConexaoCCK          string `json:"urlConexaoCCK"`
	IntervaloSincronizacao int    `json:"intervaloSincronizacao"`
}

// LoadFromJSON carrega as informações de conexão com o banco de dados de um arquivo do formato json.
func (config *Config) LoadFromJSON(fileName string) (err error) {
	var file *os.File
	var decoder *json.Decoder
	if file, err = os.Open(fileName); err != nil {
		return err
	}
	decoder = json.NewDecoder(file)
	if err = decoder.Decode(&LoadedConfigs); err != nil {
		return err
	}
	return nil
}
