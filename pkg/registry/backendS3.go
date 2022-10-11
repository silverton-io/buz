// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package registry

import (
	"context"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
)

type S3SchemaCacheBackend struct {
	bucket     string
	path       string
	client     *s3.Client
	downloader *manager.Downloader
}

func (b *S3SchemaCacheBackend) Initialize(conf config.Backend) error {
	log.Debug().Msg("🟡 initializing s3 schema cache backend")
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not load aws config")
		return err
	}
	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)
	b.bucket, b.path, b.client, b.downloader = conf.Bucket, conf.Path, client, downloader
	return nil
}

func (b *S3SchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	ctx := context.Background()
	var schemaLocation string
	if b.path == "/" {
		schemaLocation = schema
	} else {
		schemaLocation = filepath.Join(b.path, schema)
	}
	buffer := manager.NewWriteAtBuffer([]byte{})
	log.Debug().Msg("🟡 getting file from s3 backend " + schemaLocation)
	_, err = b.downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(schemaLocation),
	})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not get file from s3: " + schemaLocation)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (b *S3SchemaCacheBackend) Close() {
	log.Debug().Msg("🟡 closing s3 schema cache backend")
	// This is no-op
}
