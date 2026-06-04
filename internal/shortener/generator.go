package shortener

import (
	"sync/atomic"
)

// Base62 alphabet used to turn numeric IDs into short codes (a-z, A-Z, 0-9).
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const base = uint64(len(charset))

// Generator hands out unique IDs using an atomic counter — no random collisions.
type Generator struct {
	counter uint64
}

func NewGenerator(start uint64) *Generator {
	return &Generator{counter: start}
}

func (g *Generator) NextID() uint64 {
	return atomic.AddUint64(&g.counter, 1)
}

func (g *Generator) Encode(id uint64) string {
	if id == 0 {
		return string(charset[0])
	}

	var result []byte
	for id > 0 {
		remainder := id % base
		result = append([]byte{charset[remainder]}, result...)
		id = id / base
	}
	return string(result)
}

func (g *Generator) Generate() string {
	return g.Encode(g.NextID())
}
