package api_handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ivankuchin/connectopia.org/internal/pkg/utils"
)

type register_request struct {
	Domain string `json:"domain"`
}

func RegisterDomain(w http.ResponseWriter, r *http.Request) {
	traceID := utils.GenerateTraceID()
	body, err := utils.GetClientRequestBody(traceID, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("traceID: %s, HTTP request: %s, HTTP request body: %s", traceID, r.RequestURI, string(body))

	var rr register_request
	err = json.Unmarshal(body, &rr)
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", traceID, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := domain_list.ExpireAllDomains(traceID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isExist, err := domain_list.IsExist(traceID, rr.Domain)
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", traceID, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch isExist {
	case true:
		// domain already exists
		_, err := domain_list.UpdateDomainExpirationTimer(traceID, rr.Domain)
		if err != nil {
			log.Printf("traceID: %s, HTTP request: %s, error: %s", traceID, r.RequestURI, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update domain expiration timer\n"))
		return
	case false:
		// first domain registration
		success, err := domain_list.AddDomain(traceID, rr.Domain)
		if err != nil {
			log.Printf("traceID: %s, HTTP request: %s, error: %s", traceID, r.RequestURI, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !success {
			log.Printf("traceID: %s, HTTP request: %s, error: can't add domain", traceID, r.RequestURI)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Domain registered successfully\n"))
	}
}
