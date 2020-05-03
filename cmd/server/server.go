package main

import (
	"encoding/json"
	"flag"
	"github.com/KHYehor/design-lab3/httptools"
	"github.com/KHYehor/design-lab3/signal"
	"net/http"
	"os"
	"strconv"
	"time"
)

var port = flag.Int("port", 8080, "server port")

const confResponseDelaySec = "CONF_RESPONSE_DELAY_SEC"
const confHealthFailure = "CONF_HEALTH_FAILURE"

var TRAFFIC = 0

func main() {
	h := new(http.ServeMux)

	h.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/plain")
		if failConfig := os.Getenv(confHealthFailure); failConfig == "true" {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte("FAILURE"))
		} else {
			rw.WriteHeader(http.StatusOK)
			h := strconv.Itoa(TRAFFIC)
			//log.Printf("Server's Traffic volume is (String): %s", h)
			//log.Printf("Server's Traffic volume is (INT): %d", TRAFFIC)
			_, _ = rw.Write([]byte(h))
		}
	})

	report := make(Report)

	h.HandleFunc("/api/v1/some-data", func(rw http.ResponseWriter, r *http.Request) {
		respDelayString := os.Getenv(confResponseDelaySec)
		if delaySec, parseErr := strconv.Atoi(respDelayString); parseErr == nil && delaySec > 0 && delaySec < 300 {
			time.Sleep(time.Duration(delaySec) * time.Second)
		}

		report.Process(r)

		rw.Header().Set("content-type", "application/json")
		rw.WriteHeader(http.StatusOK)
		var response = []string{
			"1", "2",
		}
		_ = json.NewEncoder(rw).Encode(response)
		// Due to this comment https://github.com/golang/go/issues/19644#issuecomment-288178212
		// I choose just to count only response traffic - if we ae talking about real project
		// I will chnge by algorithm but in aim of education I count only this
		for _, value := range response {
			TRAFFIC += len(value)
		}
	})
	h.Handle("/report", report)

	server := httptools.CreateServer(*port, h)
	server.Start()
	signal.WaitForTerminationSignal()
}
