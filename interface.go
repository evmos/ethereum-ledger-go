package ledger

import (
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
