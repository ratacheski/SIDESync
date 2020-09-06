package dto

import (
	"encoding/xml"
	"time"
)

type Medidores struct {
	XMLName   xml.Name  `xml:"medidores"`
	Medidores []Medidor `xml:"medidor"`
}

type Medidor struct {
	XMLName         xml.Name   `xml:"medidor"`
	ID              string     `xml:"id_medidor,attr"`
	Nome            string     `xml:"nm_medidor,attr"`
	PrimeiraLeitura customTime `xml:"datahora_pri"`
	UltimaLeitura   customTime `xml:"datahora_ult"`
	Leituras        []Leitura  `xml:"leitura"`
	NaoConsultar    bool
}

type Leitura struct {
	XMLName  xml.Name   `xml:"leitura"`
	DataHora customTime `xml:"datahora"`
	Ativa    float64    `xml:"ativa"`
	Reativa  float64    `xml:"reativa"`
}

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02 15:04:05" //date format
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = customTime{parse}
	return nil
}

type Telemetria struct {
	XMLName xml.Name `xml:"telemetria"`
	Medidor Medidor  `xml:"medidor"`
}
