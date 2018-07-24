package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	//Disuse now
	GoodHealth  = checkResult{URL: "http://10.11.22.77:9093/health", Result: 0, ErrorInfo: ""}
	BadHealth   = checkResult{URL: "http://10.11.22.77:9093/health", Result: 1, ErrorInfo: ""}
	checkHealth = GoodHealth
)

type (
	checkResult struct {
		URL       string
		Result    int
		ErrorInfo string
	}
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(checkHealth)

	w.Write(body)

	log.Println(ackMsg)

}

func SetHealth(err error) {
	if err == nil {
		checkHealth = GoodHealth
	}

	checkHealth = BadHealth
	checkHealth.ErrorInfo = err.Error()
}
