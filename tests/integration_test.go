package ledger

import (
	"math/big"
	"testing"

	ethLedger "github.com/evmos/ethereum-ledger-go"
	"github.com/evmos/ethereum-ledger-go/accounts"
)

func TestLedgerSignTx(t *testing.T) {
	ledger, err := ethLedger.New()
	if err != nil {
		panic("Could not create new Ledger hub")
	}

	if len(ledger.Wallets()) == 0 {
		panic("Could not find any associated Ledgers")
	}

	wallet := ledger.Wallets()[0]
	err = wallet.Open("")
	if err != nil {
		panic("Could not open wallet")
	}

	defer wallet.Close()

	path := accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}
	account, err := wallet.Derive(path, true)
	if err != nil {
		panic("Could not derive account")
	}

	var addr = "0x3535353535353535353535353535353535353535"

	tx := ethLedger.CreateTx(
		3, big.NewInt(10), 10, addr, big.NewInt(10), make([]byte, 0),
	)

	bytes, err := wallet.SignTx(account, tx, big.NewInt(0))
	if err != nil {
		println(err.Error())
		panic("Could not sign data")
	}

	t.Logf("Signed bytes: %v\n", bytes)
}
