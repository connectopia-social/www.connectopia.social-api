package domains

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func (d *Domains) AddMutex() {
	d.m = &sync.Mutex{}
}

type domain_checker_response struct {
	Resp_type string `json:"type"`
}

func parseResponse(response string) (bool, error) {
	var dcr domain_checker_response
	err := json.Unmarshal([]byte(response), &dcr)
	if err != nil {
		log.Printf("json unmarshal error: %s", err)
		return false, err
	}

	if dcr.Resp_type == "EchoResponse" {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *Domains) isCorrect(tID, _domain string) (bool, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://" + _domain + "/cgi-bin/system.cgi?action=EchoRequest")
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", tID, "https://"+_domain+"/cgi-bin/system.cgi", err)
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("traceID: %s, HTTP request: %s, HTTP response: %s", tID, "https://"+_domain+"/cgi-bin/system.cgi", resp.Status)
		return false, nil
	}

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", tID, "https://"+_domain+"/cgi-bin/system.cgi", err)
		return false, nil
	}

	parse_status, err := parseResponse(string(text))
	if err != nil {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", tID, "https://"+_domain+"/cgi-bin/system.cgi", err)
		return false, nil
	}
	if parse_status == false {
		log.Printf("traceID: %s, HTTP request: %s, error: %s", tID, "https://"+_domain+"/cgi-bin/system.cgi", "not Connectopia")
		return false, nil
	}

	return true, nil
}

func (d *Domains) AddDomain(tID, _domain string) (bool, error) {
	_isCorrect, err := d.isCorrect(tID, _domain)
	if err != nil {
		return false, err
	}

	if !_isCorrect {
		return false, nil
	}

	if d.domains == nil {
		d.domains = make(map[string]Domain)
	}

	d.m.Lock()
	d.domains[_domain] = Domain{Name: _domain, Expire: time.Now().Unix() + lifetime_secs}
	d.m.Unlock()

	return true, nil
}

func (d *Domains) UpdateDomainExpirationTimer(tID, _domain string) (bool, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.domains[_domain]; ok {
		d.domains[_domain] = Domain{Name: _domain, Expire: time.Now().Unix() + lifetime_secs}
		return true, nil
	} else {
		return false, nil
	}
}

func (d *Domains) GetDomains(tID string) ([]Domain, error) {
	d.m.Lock()
	defer d.m.Unlock()

	var domain_slice []Domain
	for _, v := range d.domains {
		domain_slice = append(domain_slice, v)
	}

	return domain_slice, nil
}

func (d Domains) IsExist(tID, _domain string) (bool, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.domains[_domain]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *Domains) DeleteDomain(tID, _domain string) (bool, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.domains[_domain]; ok {
		delete(d.domains, _domain)
		return true, nil
	} else {
		return false, nil
	}
}

func (d *Domains) ExpireAllDomains(tID string) (bool, error) {
	d.m.Lock()
	defer d.m.Unlock()

	curr_time := time.Now().Unix()
	for k, v := range d.domains {
		if curr_time > v.Expire {
			delete(d.domains, k)
		}
	}

	return true, nil
}
