// Package crypto ...
package crypto

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dapperlabs/flow-go/crypto/hash"
)

// revive:disable:var-naming

// revive:enable

// Signer interface
type signer interface {
	// generatePrKey generates a private key
	generatePrivateKey([]byte) (PrivateKey, error)
	// decodePrKey loads a private key from a byte array
	decodePrivateKey([]byte) (PrivateKey, error)
	// decodePubKey loads a public key from a byte array
	decodePublicKey([]byte) (PublicKey, error)
}

// commonSigner holds the common data for all signers
type commonSigner struct {
	algo            SigningAlgorithm
	prKeyLength     int
	pubKeyLength    int
	signatureLength int
}

// newNonRelicSigner initializes a signer that does not depend on the Relic library.
func newNonRelicSigner(algo SigningAlgorithm) (signer, error) {
	switch algo {
	case ECDSA_P256:
		return newEcdsaP256(), nil
	case ECDSA_SECp256k1:
		return newEcdsaSecp256k1(), nil
	default:
		return nil, fmt.Errorf("the signature scheme %s is not supported.", algo)
	}
}

// GeneratePrivateKey generates a private key of the algorithm using the entropy of the given seed
func GeneratePrivateKey(algo SigningAlgorithm, seed []byte) (PrivateKey, error) {
	signer, err := newSigner(algo)
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}
	return signer.generatePrivateKey(seed)
}

// DecodePrivateKey decodes an array of bytes into a private key of the given algorithm
func DecodePrivateKey(algo SigningAlgorithm, data []byte) (PrivateKey, error) {
	signer, err := newSigner(algo)
	if err != nil {
		return nil, fmt.Errorf("decode private key failed: %w", err)
	}
	return signer.decodePrivateKey(data)
}

// DecodePublicKey decodes an array of bytes into a public key of the given algorithm
func DecodePublicKey(algo SigningAlgorithm, data []byte) (PublicKey, error) {
	signer, err := newSigner(algo)
	if err != nil {
		return nil, fmt.Errorf("decode public key failed: %w", err)
	}
	return signer.decodePublicKey(data)
}

// Signature type tools

// Bytes returns a byte array of the signature data
func (s Signature) Bytes() []byte {
	return s[:]
}

// String returns a String representation of the signature data
func (s Signature) String() string {
	const zero = "00"
	var sb strings.Builder
	sb.WriteString("0x")
	for _, i := range s {
		hex := strconv.FormatUint(uint64(i), 16)
		sb.WriteString(zero[:2-len(hex)])
		sb.WriteString(hex)
	}
	return sb.String()
}

// Key Pair

// PrivateKey is an unspecified signature scheme private key
type PrivateKey interface {
	// Algorithm returns the signing algorithm related to the private key.
	Algorithm() SigningAlgorithm
	// KeySize return the key size in bytes.
	KeySize() int
	// Sign generates a signature using the provided hasher.
	Sign([]byte, hash.Hasher) (Signature, error)
	// PublicKey returns the public key.
	PublicKey() PublicKey
	// Encode returns a bytes representation of the private key
	Encode() ([]byte, error)
	// Equals returns true if the given PrivateKeys are equal. Keys are considered unequal if their algorithms are
	// unequal or if their encoded representations are unequal. If the encoding of either key fails, they are considered
	// unequal as well.
	Equals(PrivateKey) bool
}

// PublicKey is an unspecified signature scheme public key.
type PublicKey interface {
	// Algorithm returns the signing algorithm related to the public key.
	Algorithm() SigningAlgorithm
	// KeySize return the key size in bytes.
	KeySize() int
	// Verify verifies a signature of an input message using the provided hasher.
	Verify(Signature, []byte, hash.Hasher) (bool, error)
	// Encode returns a bytes representation of the public key.
	Encode() ([]byte, error)
	// Equals returns true if the given PublicKeys are equal. Keys are considered unequal if their algorithms are
	// unequal or if their encoded representations are unequal. If the encoding of either key fails, they are considered
	// unequal as well.
	Equals(PublicKey) bool
}
