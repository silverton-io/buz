// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type bigqueryEnvelope envelope.Envelope

func getEnvelopeSchema() bigquery.Schema {
	return bigquery.Schema{
		{Name: "uuid", Type: bigquery.StringFieldType},
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "buz_timestamp", Type: bigquery.TimestampFieldType},
		{Name: "buz_version", Type: bigquery.StringFieldType},
		{Name: "buz_name", Type: bigquery.StringFieldType},
		{Name: "buz_env", Type: bigquery.StringFieldType},
		{Name: "protocol", Type: bigquery.StringFieldType},
		{Name: "schema", Type: bigquery.StringFieldType},
		{Name: "vendor", Type: bigquery.StringFieldType},
		{Name: "namespace", Type: bigquery.StringFieldType},
		{Name: "version", Type: bigquery.StringFieldType},
		{Name: "is_valid", Type: bigquery.BooleanFieldType},
		{Name: "validation_error", Type: bigquery.JSONFieldType},
		{Name: "contexts", Type: bigquery.JSONFieldType},
		{Name: "payload", Type: bigquery.JSONFieldType},
	}
}

func (e *bigqueryEnvelope) Save() (map[string]bigquery.Value, string, error) {
	vErr, _ := e.ValidationError.AsByte()
	contexts, _ := e.Contexts.AsByte()
	payload, _ := e.Payload.AsByte()

	return map[string]bigquery.Value{
		"uuid":             e.Uuid,
		"timestamp":        e.Timestamp,
		"buz_timestamp":    e.BuzTimestamp,
		"buz_version":      e.BuzVersion,
		"buz_name":         e.BuzName,
		"buz_env":          e.BuzEnv,
		"protocol":         e.Protocol,
		"schema":           e.Schema,
		"vendor":           e.Vendor,
		"namespace":        e.Namespace,
		"version":          e.Version,
		"is_valid":         e.IsValid,
		"validation_error": string(vErr),
		"contexts":         string(contexts),
		"payload":          string(payload),
	}, bigquery.NoDedupeID, nil
}

type Sink struct {
	metadata backendutils.SinkMetadata
	client   *bigquery.Client
	dataset  *bigquery.Dataset
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, conf.Project)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not instantiate a bigquery client to project " + conf.Project)
		return err
	}
	s.client = client
	s.dataset = client.Dataset(conf.Dataset)
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
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
	inserter := s.dataset.Table(output).Inserter()
	var items []*bigqueryEnvelope
	for _, e := range envelopes {
		bqEnvelope := bigqueryEnvelope(e)
		items = append(items, &bqEnvelope)
	}
	err := inserter.Put(ctx, items)
	if err != nil {
		log.Warn().Err(err).Msg("could not dequeue envelopes - ensuring table exists")
		tbl := s.dataset.Table(output)
		md, _ := tbl.Metadata(ctx)
		if md == nil {
			// The table doesn't exist and should be created
			schema := getEnvelopeSchema()
			tblMetadata := &bigquery.TableMetadata{
				Schema: schema,
			}
			log.Debug().Msg("creating " + output + " table")
			err := tbl.Create(ctx, tblMetadata)
			if err != nil {
				log.Error().Err(err).Msg("could not create " + output + " table")
				return err
			}
			if err != nil {
				log.Error().Err(err).Msg("could not dequeue envelopes")
				return err
			}
		}
	}
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	err := s.client.Close()
	return err
}
