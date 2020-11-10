package main

import (
	"fmt"
	"github.com/AceDarkknight/GoProxyCollector/collector"
	"github.com/AceDarkknight/GoProxyCollector/scheduler"
	"github.com/AceDarkknight/GoProxyCollector/server"
	"github.com/AceDarkknight/GoProxyCollector/storage"
	"github.com/AceDarkknight/GoProxyCollector/verifier"
	"time"

	"github.com/cihub/seelog"
)

type IpInfo struct {
	IP string `json:"ip"`
	Port int `json:"port"`
	Location string `json:"location"`
	Source string `json:"source"`
	Speed int `json:"speed"`
}
func main() {
	// Load log.
	scheduler.SetLogger("logConfig.xml")
	defer func() {
		seelog.Flush()
		r:=recover()
		if r!=nil{
			fmt.Println(r)
		}
	}()
	// Load database.
	database, err := storage.NewStorage()
	defer database.Close()
	if err != nil {
		seelog.Critical(err)
		panic(err)
	}

	seelog.Infof("database initialize finish.")
	// Start server
	go server.NewServer(database)
	// Verify storage every 5min.
	verifyTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range verifyTicker.C {
			verifier.VerifyAndDelete(database)
			seelog.Debug("verify database.")
		}
	}()

	configs := collector.NewCollectorConfig("collectorConfig.xml")
	scheduler.Run(configs, database)
}