package shortener

import (
	"sync/atomic"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const base = uint64(len(charset))

// Generator Stuct
type Generator struct {
	counter uint64
}

// NewGenerator initializes and returns a new instance of the Generator struct, Constructor for Generator.
func NewGenerator(start uint64) *Generator {
	return &Generator{
		counter: start,
	}
}

// NextID atomically increments the counter and returns the new value, ensuring thread safety when generating unique IDs.
func (g *Generator) NextID() uint64 {
	return atomic.AddUint64(&g.counter, 1)
}

// Encode converts uint64 ID into string representation using the charset, encoding the ID into a short code format.
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

// Generate creates a new unique short code by generating the next ID and encoding it, combining the functionality of NextID and Encode to produce a short code for URL shortening.
func (g *Generator) Generate() string {
	id := g.NextID()
	return g.Encode(id)
}
