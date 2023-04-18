// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package s3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/util"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	input    chan []envelope.Envelope
	uploader *manager.Uploader
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not load aws config")
		return err
	}
	client := s3.NewFromConfig(cfg)
	s.uploader = manager.NewUploader(client)
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

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, bucket string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	reader, writer := io.Pipe()
	gzipWriter := gzip.NewWriter(writer)
	go func() {
		for _, envelope := range envelopes {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(envelope)
			if err != nil {
				log.Error().Err(err).Msg("could not encode envelope")
			}
			envelopeReader := bufio.NewReader(&buf)
			envelopeReader.WriteTo(gzipWriter)
		}
		gzipWriter.Close()
		writer.Close()
	}()
	uuid := uuid.New()
	today := time.Now()
	fileName := "year=" + strconv.Itoa(today.Year()) + "/month=" + strconv.Itoa(int(today.Month())) + "/day=" + strconv.Itoa(today.Day()) + "/" + uuid.String() + "-events.json.gz"
	util.Pprint(fileName)
	input := s3.PutObjectInput{
		Body:   reader,
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}
	uploadOutput, err := s.uploader.Upload(ctx, &input)
	if err != nil {
		log.Error().Err(err).Msg("could not upload file")
	}
	log.Debug().Interface("uploadOutput", uploadOutput).Msg("uploaded file")
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
