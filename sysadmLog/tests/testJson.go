package main

import (
	"os"

	log "sysadm/sysadmLog"
)


func main() {
	
	logger := log.New()
	logger.Out = os.Stdout
	logger.ReportCaller = true
	jsonF := &log.JSONFormatter{}
	jsonF.DisableTimestamp = true
	logger.SetFormatter(jsonF)
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
}
