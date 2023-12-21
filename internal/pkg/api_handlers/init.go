package api_handlers

import (
	"github.com/ivankuchin/connectopia.org/internal/pkg/domains"
)

var domain_list domains.Domains

func init() {
	domain_list.AddMutex()
}
