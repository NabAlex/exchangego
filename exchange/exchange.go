package exchange

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shmel1k/exchangego/context/errs"
)

type User struct {
	ID               uint32
	Name             string
	Password         string
	RegistrationDate time.Time
}

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func WriteOK(w http.ResponseWriter, data interface{}) {
	// FIXME(shmel1k): add easyjson or something like that
	r := Response{
		Status: http.StatusOK,
		Body:   data,
	}
	dt, err := json.Marshal(r)
	if err != nil {
		errs.WriteInternal(w)
		return
	}
	w.Write(dt)
}
