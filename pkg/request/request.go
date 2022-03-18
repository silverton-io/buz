package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/event"
)

const (
	JSON_CONTENT_TYPE string = "application/json"
)

func PostEvent(url url.URL, payload event.SelfDescribingEvent) error {
	data, err := json.Marshal(payload)
	buff := bytes.NewBuffer(data)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not marshal payload")
		return err
	}
	_, err = http.Post(url.String(), JSON_CONTENT_TYPE, buff)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not post payload")
		return err
	}
	return nil
}

func PostEnvelopes(url url.URL, envelopes []envelope.Envelope) error {
	data, err := json.Marshal(envelopes)
	buff := bytes.NewBuffer(data)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not marshal envelopes")
		return err
	}
	_, err = http.Post(url.String(), JSON_CONTENT_TYPE, buff)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not post envelopes")
		return err
	}
	return nil
}

func Get(url url.URL) (body []byte, err error) {
	resp, err := http.Get(url.String())
	if err != nil {
		log.Trace().Err(err).Msg("could not get url " + url.String())
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
