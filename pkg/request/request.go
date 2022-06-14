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

func PostPayload(url url.URL, payload interface{}) (resp *http.Response, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("could not marshal payload")
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	resp, err = http.Post(url.String(), JSON_CONTENT_TYPE, buf)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		log.Debug().Interface("response", string(bodyBytes)).Msg("non-200 response")
		if err != nil {
			log.Error().Err(err).Msg("could not read response body")
			return nil, err
		}
	}
	if err != nil {
		log.Error().Err(err).Msg("could not post payload")
		return nil, err
	}
	return resp, nil
}

func PostEvent(url url.URL, payload event.SelfDescribingEvent) (resp *http.Response, err error) {
	resp, err = PostPayload(url, payload)
	return resp, err
}

func PostEnvelopes(url url.URL, envelopes []envelope.Envelope) (resp *http.Response, err error) {
	resp, err = PostPayload(url, envelopes)
	return resp, err
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
