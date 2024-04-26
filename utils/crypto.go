package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type Cryptoer interface {
	GenerateFromPassword(password []byte) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type Crypto struct{}

func NewCrypto() Cryptoer {
	return Crypto{}
}

func (c Crypto) GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (c Crypto) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
