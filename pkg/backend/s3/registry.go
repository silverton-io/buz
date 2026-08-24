// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package s3

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
)

type RegistryBackend struct {
	bucket string
	path   string
	client *transfermanager.Client
}

func (b *RegistryBackend) Initialize(conf config.Backend) error {
	log.Debug().Msg("🟡 initializing s3 schema cache backend")
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not load aws config")
		return err
	}
	s3Client := s3.NewFromConfig(cfg)
	tmClient := transfermanager.New(s3Client)
	b.bucket, b.path, b.client = conf.Bucket, conf.Path, tmClient
	return nil
}

func (b *RegistryBackend) GetRemote(schema string) (contents []byte, err error) {
	ctx := context.Background()
	var schemaLocation string
	if b.path == "/" {
		schemaLocation = schema
	} else {
		schemaLocation = filepath.Join(b.path, schema)
	}
	log.Debug().Msg("🟡 getting file from s3 backend " + schemaLocation)
	output, err := b.client.GetObject(ctx, &transfermanager.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(schemaLocation),
	})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not get file from s3: " + schemaLocation)
		return nil, err
	}
	data, err := io.ReadAll(output.Body)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not read s3 object body: " + schemaLocation)
		return nil, err
	}
	return data, nil
}

func (b *RegistryBackend) Close() {
	log.Debug().Msg("🟡 closing s3 schema cache backend")
	// This is no-op
}
