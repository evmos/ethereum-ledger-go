package ledger

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethereum-ledger-go/accounts"
	"github.com/evmos/ethereum-ledger-go/types"
	"github.com/evmos/ethereum-ledger-go/usbwallet"
)

type EthereumLedger struct {
	hub *usbwallet.Hub
}

func (el EthereumLedger) Wallets() []accounts.Wallet {
	return el.hub.Wallets()
}

func New() (*EthereumLedger, error) {
	l := &EthereumLedger{}
	hub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return &EthereumLedger{}, err
	}

	l.hub = hub
	return l, nil
}

func CreateTx(
	nonce uint64,
	gasPrice *big.Int,
	gas uint64,
	to string,
	value *big.Int,
	data []byte,
) *types.Transaction {
	addrBytes, err := hex.DecodeString(strings.TrimPrefix(to, "0x"))
	if err != nil {
		panic(fmt.Sprintf("Could not convert \"to\" field to bytes with error: %v\n", err.Error()))
	}
	if len(addrBytes) != 20 {
		panic(fmt.Sprintf("Improper size of \"to\" field, got %v, expected 20\n", len(addrBytes)))
	}

	var addrObj = &common.Address{}
	copy(addrObj[:], addrBytes)

	return types.NewLedgerTx(
		nonce,
		gasPrice,
		gas,
		addrObj,
		value,
		data,
	)
}
