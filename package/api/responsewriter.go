package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

func RawResponseWriter(w http.ResponseWriter, code int, payload []byte, contenttype string) {
	w.Header().Add("Content-Type", contenttype)
	w.WriteHeader(code)
	w.Write(payload)
}

func ErrorResponseWriter(w http.ResponseWriter, code int, err error) {
	jsonResponseWriter(w, code, Response{Status: "Error", Error: err.Error()})
}

func ResponseWriter(w http.ResponseWriter, code int, payload interface{}) {
	jsonResponseWriter(w, code, Response{Status: "OK", Data: payload})
}

func jsonResponseWriter(w http.ResponseWriter, code int, raw interface{}) {
	payload, e := json.Marshal(raw)
	if e != nil {
		ErrorResponseWriter(w, 500, errors.New("JSON Marshal error"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(payload)
	}
}
