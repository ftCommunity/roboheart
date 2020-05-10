package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func RequestLoader(r *http.Request, w http.ResponseWriter, v interface{}) bool {
	bodyb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponseWriter(w, 400, err)
		return false
	}
	if err := json.Unmarshal(bodyb, v); err != nil {
		ErrorResponseWriter(w, 400, err)
		return false
	}
	return true
}
