package http

import (
	"bytes"
	"encoding/json"
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
