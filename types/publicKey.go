package types

import (
	"encoding/hex"
)

// Lengths of hashes and addresses in bytes.
const (
	// PublicKeyLength is the expected length of a public key
	PublicKeyLength = 65
)

// Public Key represents the 40 byte address of an Ethereum public key
type PublicKey [PublicKeyLength]byte

func (pk PublicKey) Bytes() []byte { return pk[:] }

func (pk PublicKey) Hex() string {
	var buf [len(pk)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], pk[:])
	return string(buf[:])
}
