package ledger

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethereum-ledger-go/accounts"
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
	to string,
	gas uint64,
	gasPrice *big.Int,
	amount *big.Int,
	data []byte,
) (*coretypes.Transaction, error) {
	if !common.IsHexAddress(to) {
		return nil, fmt.Errorf("invalid 'to' address: %s", to)
	}

	toAddr := common.HexToAddress(to)

	return coretypes.NewTransaction(
		nonce,
		toAddr,
		amount,
		gas,
		gasPrice,
		data,
	), nil
}
