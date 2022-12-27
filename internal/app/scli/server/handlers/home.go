package handlers

import (
	"github.com/shammishailaj/scli/pkg/schemas"
	"github.com/shammishailaj/scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

type Home struct {
	l *log.Logger
}

func NewHome(l *log.Logger) *Home {
	return &Home{
		l: l,
	}
}

func (h *Home) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// https://stackoverflow.com/a/40382340
	hp := schemas.HomePage{
		AppName: strings.Title(strings.ToLower(viper.Get("app_name").(string))),
		AppDesc: viper.Get("app_desc").(string),
	}

	h.l.Infof("hp = %v", hp)
	u := utils.NewUtils(h.l)
	u.RenderTemplate(w, "", "home", hp)
}
