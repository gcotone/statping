package handlers

import (
	"fmt"
	"github.com/hunterlong/statup/core"
	"net/http"
	"strings"
)

func PrometheusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Prometheus /metrics Request From IP: %v\n", r.RemoteAddr)
	metrics := []string{}
	system := fmt.Sprintf("statup_total_failures %v\n", core.CountFailures())
	system += fmt.Sprintf("statup_total_services %v", len(core.CoreApp.Services))
	metrics = append(metrics, system)

	for _, v := range core.CoreApp.Services {
		online := 1
		if !v.Online {
			online = 0
		}
		met := fmt.Sprintf("statup_service_failures{id=\"%v\" name=\"%v\"} %v\n", v.Id, v.Name, len(v.Failures))
		met += fmt.Sprintf("statup_service_latency{id=\"%v\" name=\"%v\"} %0.0f\n", v.Id, v.Name, (v.Latency * 100))
		met += fmt.Sprintf("statup_service_online{id=\"%v\" name=\"%v\"} %v\n", v.Id, v.Name, online)
		met += fmt.Sprintf("statup_service_status_code{id=\"%v\" name=\"%v\"} %v\n", v.Id, v.Name, v.LastStatusCode)
		met += fmt.Sprintf("statup_service_response_length{id=\"%v\" name=\"%v\"} %v", v.Id, v.Name, len([]byte(v.LastResponse)))
		metrics = append(metrics, met)
	}
	output := strings.Join(metrics, "\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}