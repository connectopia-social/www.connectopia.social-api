package api_handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ivankuchin/connectopia.org/internal/pkg/utils"
)

func GetDomains(w http.ResponseWriter, r *http.Request) {
	traceID := utils.GenerateTraceID()
	log.Printf("traceID: %s, HTTP request: %s", traceID, r.RequestURI)

	if _, err := domain_list.ExpireAllDomains(traceID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	domains_slice, err := domain_list.GetDomains(traceID)
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", traceID, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	domains_json, err := json.Marshal(domains_slice)
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", traceID, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(domains_json))
}
