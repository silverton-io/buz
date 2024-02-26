// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package materializewebhook

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/backend/postgresdb"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	gormDb   *gorm.DB
	input    chan []envelope.Envelope
	shutdown chan int
	host     string
	apiKey   string
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	connParams := db.ConnectionParams{
		Host: conf.Hosts[0], //Only use the first configured host.
		Port: conf.Port,
		Db:   conf.Database,
		User: conf.User,
		Pass: conf.Password,
	}
	s.host = conf.Hosts[0]
	s.apiKey = conf.ApiKey
	connString := postgresdb.GenerateDsn(connParams)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open " + s.metadata.SinkType + " connection")
		return err
	}
	s.gormDb = gormDb
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	for _, tbl := range []string{s.metadata.DefaultOutput, s.metadata.DeadletterOutput} {
		ensureErr := db.EnsureMaterializeWebhookSource(s.gormDb, tbl, conf.WebhookSecret, conf.Cluster, &envelope.JsonbEnvelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	key := []byte(s.apiKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(""))
	signature := hex.EncodeToString(h.Sum(nil))

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("X-Signature", signature)

	urlGenerated := fmt.Sprintf("https://%s/public/%s", s.host, output)
	url, err := url.Parse(urlGenerated)
	if err != nil {
		log.Error().Err(err).Msg("🔴 " + output + " is not a valid url")
		return err
	}
	_, err = request.PostEnvelopes(*url, envelopes, nil)
	if err != nil {
		log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("🔴 could not dequeue payloads")
	}
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	db, _ := s.gormDb.DB()
	s.shutdown <- 1
	err := db.Close()
	return err
}
