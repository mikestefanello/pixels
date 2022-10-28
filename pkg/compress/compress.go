package compress

import (
	"cloud.google.com/go/pubsub"
)

// This will be a package shared outside of this application!

const (
	MessageKeyCompressionEnabled = "ce"
)

type Compressor interface {
	Compress([]byte) ([]byte, error)
	CompressMessage(message *pubsub.Message) error
}

type Decompressor interface {
	Decompress([]byte) ([]byte, error)
	DecompressMessage(message *pubsub.Message) error
}
