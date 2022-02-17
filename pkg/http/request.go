package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	JSON_CONTENT_TYPE string = "application/json"
)

func SendJson(host string, payload interface{}) {
	data, _ := json.Marshal(payload)
	buff := bytes.NewBuffer(data)
	_, err := http.Post(host, JSON_CONTENT_TYPE, buff)
	if err != nil {
		log.Trace().Err(err).Msg("could not send payload to " + host)
	}
}

func Get(uri string) (body []byte, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Trace().Err(err).Msg("could not get url " + uri)
		return nil, err
	}
	defer resp.Body.Close()
	body, ioerr := io.ReadAll(resp.Body)
	if ioerr != nil {
		log.Trace().Err(ioerr).Msg("could not read response body")
		return nil, ioerr
	}
	return body, nil
}
