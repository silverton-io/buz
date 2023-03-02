// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/envelope"
)

const (
	JSON_CONTENT_TYPE string = "application/json"
)

func PostPayload(url url.URL, payload interface{}) (resp *http.Response, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not marshal payload")
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	resp, err = http.Post(url.String(), JSON_CONTENT_TYPE, buf)
	if resp == nil {
		return resp, nil
	}
	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not read response body")
			return nil, err
		}
	}
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not post payload")
		return nil, err
	}
	return resp, nil
}

func PostEvent(url url.URL, payload envelope.SelfDescribingEvent) (resp *http.Response, err error) {
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
