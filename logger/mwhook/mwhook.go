package mwhook

import (
	"io"
	"log"
	"os"
	"sidegosynchronizer/logger"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// WriterMap é um mapa de writers para cada novo splitField
type WriterMap map[string]io.Writer

//SplitedOption é uma estrutura de opções para cada logger novo criado
type SplitedOption struct {
	LogDir         string
	MaxSize        int
	RotatePeriod   time.Duration
	LocalTime      bool
	RotateByPeriod bool
}

//MwHook é um hook para sirupsen/logrus que escreve logs em arquivos locais
type MwHook struct {
	writers          WriterMap
	lock             *sync.Mutex
	splitField       string
	splitedOption    SplitedOption
	defaultWriter    io.Writer
	hasDefaultWriter bool
	hasSplitField    bool
}

// SetDefaultWriter define o writer padrão
func (hook *MwHook) SetDefaultWriter(defaultWriter io.Writer) {
	hook.defaultWriter = defaultWriter
	hook.hasDefaultWriter = true
}

// SetSplitField é o campo usado para dividir o log através do valor deste campo
func (hook *MwHook) SetSplitField(splitField string, option SplitedOption) {
	if splitField != "" {
		hook.splitField = splitField
		hook.hasSplitField = true
		if option != (SplitedOption{}) {
			hook.splitedOption = option
		}
	}
}

// Fire escreve o arquivo de log para writer definido
func (hook *MwHook) Fire(entry *logrus.Entry) error {
	if hook.writers != nil || hook.hasDefaultWriter {
		return hook.ioWrite(entry)
	}
	return nil
}

// Escreve uma linha de log em um io.Writer.
func (hook *MwHook) ioWrite(entry *logrus.Entry) error {
	var (
		writer io.Writer
		msg    string
		err    error
	)

	hook.lock.Lock()
	defer hook.lock.Unlock()

	if hook.hasSplitField {
		if splitField, ok := entry.Data[hook.splitField]; ok {
			value := splitField.(string)
			if _, found := hook.writers[value]; !found {
				logPath := hook.splitedOption.LogDir + value + "/"
				err = CreateDirIfNotExist(logPath)
				if err != nil {
					log.Println("[mwhook] - Erro ao criar pasta de log: ", err)
					logPath = hook.splitedOption.LogDir
				}
				l := &logger.Logger{
					Filename:     logPath + value + ".log",
					MaxSize:      hook.splitedOption.MaxSize,      //MB
					RotatePeriod: hook.splitedOption.RotatePeriod, //Hours
					LocalTime:    hook.splitedOption.LocalTime,
				}
				if hook.splitedOption.RotateByPeriod {
					l.RotatePeriodically()
				}
				hook.writers[value] = l
			}
			writer = hook.writers[value]
		}
	}
	if writer == nil {
		if hook.hasDefaultWriter {
			writer = hook.defaultWriter
		} else {
			return nil
		}
	}

	msg, err = entry.String()
	if err != nil {
		log.Println("[mwhook] - Falha ao gerar mensagem["+entry.Message+"] de log:", err)
		return err
	}
	_, err = writer.Write([]byte(msg))
	return err
}

// Levels retornam níveis de log configurados
func (hook *MwHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func CreateDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0775)
		if err != nil {
			return err
		}
	}
	return err
}
