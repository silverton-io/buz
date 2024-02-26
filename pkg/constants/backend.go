// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package constants

const (
	// Databases
	POSTGRES            string = "postgres"
	MYSQL               string = "mysql"
	MATERIALIZE         string = "materialize"
	MATERIALIZE_WEBHOOK string = "materializeWebhook"
	CLICKHOUSE          string = "clickhouse"
	MONGODB             string = "mongodb"
	ELASTICSEARCH       string = "elasticsearch"
	TIMESCALE           string = "timescale"
	BIGQUERY            string = "bigquery"
	// Streams and Queues
	PUBSUB           string = "pubsub"
	REDPANDA         string = "redpanda"
	KAFKA            string = "kafka"
	KINESIS          string = "kinesis"
	KINESIS_FIREHOSE string = "kinesis-firehose"
	NATS             string = "nats"
	NATS_JETSTREAM   string = "nats-jetstream"
	EVENTBRIDGE      string = "eventbridge"
	// Object Stores
	GCS   string = "gcs"
	S3    string = "s3"
	MINIO string = "minio"
	// System
	STDOUT    string = "stdout"
	BLACKHOLE string = "blackhole"
	FILE      string = "file"
	// Web
	HTTP  string = "http"
	HTTPS string = "https"
	// Third Party
	INDICATIVE string = "indicative"
	AMPLITUDE  string = "amplitude"
	PUBNUB     string = "pubnub"
	IGLU       string = "iglu"
	SPLUNK     string = "splunk"
)
