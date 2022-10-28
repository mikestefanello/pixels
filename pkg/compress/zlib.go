package compress

import (
	"compress/zlib"
	"sync"

	"cloud.google.com/go/pubsub"
)

type zlibCompressor struct {
	pool sync.Pool
}

func NewZlibCompressor() Compressor {
	return &zlibCompressor{
		pool: sync.Pool{
			New: func() interface{} {
				return zlib.NewWriter(nil)
			},
		},
	}
}

func (z *zlibCompressor) Compress(src []byte) ([]byte, error) {
	// todo
	return src, nil
}

func (z *zlibCompressor) CompressMessage(msg *pubsub.Message) error {
	// todo
	msg.Attributes[MessageKeyCompressionEnabled] = "true"
	return nil
}

type zlibDecompressor struct {
	pool sync.Pool
}

func NewZlibDecompressor() Decompressor {
	return &zlibDecompressor{
		pool: sync.Pool{
			New: func() interface{} {
				return nil // todo
			},
		},
	}
}

func (z *zlibDecompressor) Decompress(src []byte) ([]byte, error) {
	// todo
	return src, nil
}

func (z *zlibDecompressor) DecompressMessage(msg *pubsub.Message) error {
	// todo
	if msg.Attributes[MessageKeyCompressionEnabled] == "true" {
		// todo
	}
	return nil
}
