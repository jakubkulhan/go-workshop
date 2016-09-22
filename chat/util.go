package chat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func writeResponse(w http.ResponseWriter, r *http.Request, v interface{}, statusCode int) error {
	var buf []byte
	var err error

	if r.URL.Query().Get("pretty") == "true" {
		buf, err = json.MarshalIndent(v, "", "  ")
		buf = append(buf, '\n')
	} else {
		buf, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	for nWritten := 0; nWritten < len(buf); {
		n, err := w.Write(buf[nWritten:])
		if err != nil {
			return err
		}
		nWritten += n
	}

	return nil
}

func WriteOkResponse(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return writeResponse(w, r, v, 200)
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, errorCode int, err error) error {
	statusCode := errorCode / 1000

	return writeResponse(w, r, &ErrorResponse{
		OK:         false,
		StatusCode: statusCode,
		Code:       errorCode,
		Message:    err.Error(),
	}, statusCode)
}

func ReadRequest(r *http.Request, v interface{}) (err error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(buf, v)
}
