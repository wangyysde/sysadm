package main

import (
	"os"
	"time"

	log "sysadm/sysadmLog"
)


func main() {
	
	logger := log.New()
	logger.Out = os.Stdout
	logger.ReportCaller = false
	textF := &log.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              true,
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,                                                                                                                                                                                                      
		FullTimestamp:             true,
		TimestampFormat:           time.RFC3339,
		DisableSorting:            true,
		DisableLevelTruncation:    true,
		PadLevelText:              false,
	}
	textF.PadLevelText = false
	//textF.DisableTimestamp = true
	logger.SetFormatter(textF)
	msg :=  log.Fields {
		"name": "MyName",
		"age": 12,
		"sex": true,
		"msg": "test my msg",
	}

	logger.WithFields(log.Fields(msg)).Info("aaaa")	

	msgnew :=  log.Fields {
		"name": "MyName",
		"age": 12,
		"sex": true,
		"msg": "Not message",
	}

	logger.WithFields(log.Fields(msgnew)).Info("")	
	logger.WithFields(log.Fields(msgnew))
	logger.Error()
	logger.Info("test message")
}
