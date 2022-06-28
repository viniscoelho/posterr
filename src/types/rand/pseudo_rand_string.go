package rand

import (
	"math/rand"
	"time"
)

const charSet = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type PseudoRand struct {
	seed   *rand.Rand
	values map[string]struct{}
}

func NewPseudoRandomString() PseudoRand {
	return PseudoRand{
		seed:   rand.New(rand.NewSource(time.Now().UnixNano())),
		values: make(map[string]struct{}),
	}
}

// GenerateUnique returns a non previously generated pseudo random string
func (r PseudoRand) GenerateUnique(length int) string {
	s := r.newString(length)
	// Ensures that a random string is generated
	for _, exists := r.values[s]; exists; s = r.newString(length) {
	}
	r.values[s] = struct{}{}
	return s
}

// GenerateAny returns a pseudo random string
func (r PseudoRand) GenerateAny(length int) string {
	return r.newString(length)
}

func (r PseudoRand) newString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[r.seed.Intn(len(charSet))]
	}
	return string(b)
}
