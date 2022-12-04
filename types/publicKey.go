package types

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Since go-ethereum does not have a public key type, we create our own here

// Lengths of hashes and addresses in bytes.
const (
	// PublicKeyLength is the expected length of a public key
	PublicKeyLength = 65
)

// Public Key represents the 65 byte address of an Ethereum public key
type PublicKey [PublicKeyLength]byte

func (pk PublicKey) Bytes() []byte { return pk[:] }

func (pk PublicKey) Hex() string { return hexutil.Encode(pk[:]) }
