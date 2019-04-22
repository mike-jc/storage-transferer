package modelsDataEncryptor

import (
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

type HybridCryptoOptions struct {
	SymmetricalMethod int
	Hash              hash.Hash
	Mgf1Hash          hash.Hash
}

func NewHybridCryptoOptionsDefault() HybridCryptoOptions {
	return HybridCryptoOptions{
		SymmetricalMethod: AesGcm256,
		Hash:              sha256.New(),
		Mgf1Hash:          sha256.New(),
	}
}

func NewHybridCryptoOptionsForDracoon() HybridCryptoOptions {
	return HybridCryptoOptions{
		SymmetricalMethod: AesGcm256,
		Hash:              sha256.New(),
		Mgf1Hash:          sha1.New(),
	}
}
