package handlers

import (
	jsonresp "github.com/shammishailaj/scli/pkg/http/response/json"
	"github.com/shammishailaj/scli/pkg/schemas"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Ping struct {
	l      *log.Logger
	semVer *schemas.SemanticVersion
}

func NewPing(l *log.Logger, version *schemas.SemanticVersion) *Ping {
	return &Ping{
		l:      l,
		semVer: version,
	}
}

func (p *Ping) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	jsonresp.OK(w, r, p.semVer)
}
