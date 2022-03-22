package config

type Manifold struct {
	BufferRecordThreshold int `json:"bufferRecordThreshold"`
	BufferByteThreshold   int `json:"bufferByteThreshold"`
	BufferTimeThreshold   int `json:"bufferTimeThreshold"`
}
