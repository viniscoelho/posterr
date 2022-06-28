package gen

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

func (r PseudoRand) Generate(length int) string {
	s := r.newString(length, charSet)
	// ensures that a random string is generated
	for _, exists := r.values[s]; exists; s = r.newString(length, charSet) {
	}
	r.values[s] = struct{}{}
	return s
}

func (r PseudoRand) newString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.seed.Intn(len(charset))]
	}
	return string(b)
}
