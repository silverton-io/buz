package cache

import "github.com/silverton-io/gosnowplow/pkg/config"

type S3SchemaCacheBackend struct {
	location string
	path     string
}

func (b *S3SchemaCacheBackend) Initialize(config config.SchemaCacheBackend) {

}

func (b *S3SchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	return nil, nil
}
