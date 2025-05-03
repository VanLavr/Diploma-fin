package hasher

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	cost int
}

var Hshr = NewHasher(4)

func NewHasher(cost int) *Hasher {
	return &Hasher{cost: cost}
}

func (h *Hasher) Hash(data string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), h.cost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

func (h *Hasher) Validate(origin string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(origin), []byte(data))
	return err == nil
}
