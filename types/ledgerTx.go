package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// LedgerTx is the transaction data of regular Ethereum transactions.
type LedgerTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *LedgerTx) copy() TxData {
	cpy := &LedgerTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are initialized below.
		Value:    new(big.Int),
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *LedgerTx) txType() byte           { return 0 }
func (tx *LedgerTx) chainID() *big.Int      { return deriveChainId(tx.V) }
func (tx *LedgerTx) accessList() AccessList { return nil }
func (tx *LedgerTx) data() []byte           { return tx.Data }
func (tx *LedgerTx) gas() uint64            { return tx.Gas }
func (tx *LedgerTx) gasPrice() *big.Int     { return tx.GasPrice }
func (tx *LedgerTx) gasTipCap() *big.Int    { return tx.GasPrice }
func (tx *LedgerTx) gasFeeCap() *big.Int    { return tx.GasPrice }
func (tx *LedgerTx) value() *big.Int        { return tx.Value }
func (tx *LedgerTx) nonce() uint64          { return tx.Nonce }
func (tx *LedgerTx) to() *common.Address    { return tx.To }

func (tx *LedgerTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *LedgerTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.V, tx.R, tx.S = v, r, s
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
