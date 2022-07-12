package ledger

import (
	"math/big"
	"testing"
)

func TestMalformedAddressGen1(t *testing.T) {
	var addr = "0x35353535353535353535353535353535353535355" // Odd hex value
	defer func() { _ = recover() }()

	_ = CreateTx(
		3, big.NewInt(10), 10, addr, big.NewInt(10), make([]byte, 0),
	)

	t.Errorf("Expected panic due to malformed address")
}

func TestMalformedAddressGen2(t *testing.T) {
	var addr = "0x353535353535353535353535353535353535353535" // Hex len == 21
	defer func() { _ = recover() }()

	_ = CreateTx(
		3, big.NewInt(10), 10, addr, big.NewInt(10), make([]byte, 0),
	)

	t.Errorf("Expected panic due to malformed address")
}

func TestMalformedAddressGen3(t *testing.T) {
	var addr = "0x35353535353535353535353535353535353535" // Hex len == 19
	defer func() { _ = recover() }()

	_ = CreateTx(
		3, big.NewInt(10), 10, addr, big.NewInt(10), make([]byte, 0),
	)

	t.Errorf("Expected panic due to malformed address")
}
