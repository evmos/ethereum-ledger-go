package ledger

import (
	"fmt"
	"log"
	"math/big"
	"testing"

	ethLedger "github.com/evmos/ethereum-ledger-go"
	"github.com/evmos/ethereum-ledger-go/accounts"
	"github.com/evmos/ethereum-ledger-go/common/math"
	"github.com/evmos/ethereum-ledger-go/types"
)

func initWallet(t *testing.T) (accounts.Wallet, accounts.Account) {
	t.Helper()
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

	path := accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}
	account, err := wallet.Derive(path, true)
	if err != nil {
		panic("Could not derive account")
	}

	return wallet, account
}

func TestLedgerSignTx(t *testing.T) {
	wallet, account := initWallet(t)
	defer wallet.Close()

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

func TestLedgerSignEIP712(t *testing.T) {
	wallet, account := initWallet(t)
	log.Printf("Account public key: %v\n", account.PublicKey)
	log.Printf("Account public key hex: %v\n", account.PublicKey.Hex())
	log.Printf("Account address: %v\n", account.Address)

	defer wallet.Close()

	const primaryType = "Mail"

	domain := types.TypedDataDomain{
		Name:              "Ether Mail",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(1),
		VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		Salt:              "",
	}

	domainTypes := types.Types{
		"EIP712Domain": {
			{
				Name: "name",
				Type: "string",
			},
			{
				Name: "version",
				Type: "string",
			},
			{
				Name: "chainId",
				Type: "uint256",
			},
			{
				Name: "verifyingContract",
				Type: "address",
			},
		},
		"Person": {
			{
				Name: "name",
				Type: "string",
			},
			{
				Name: "wallet",
				Type: "address",
			},
		},
		"Mail": {
			{
				Name: "from",
				Type: "Person",
			},
			{
				Name: "to",
				Type: "Person",
			},
			{
				Name: "contents",
				Type: "string",
			},
		},
	}

	messageStandard := map[string]interface{}{
		"from": map[string]interface{}{
			"name":   "Cow",
			"wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		"to": map[string]interface{}{
			"name":   "Bob",
			"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		"contents": "Hello, Bob!",
	}

	typedData := types.TypedData{
		Types:       domainTypes,
		PrimaryType: primaryType,
		Domain:      domain,
		Message:     messageStandard,
	}

	bytes, err := wallet.SignTypedData(account, typedData)
	if err != nil {
		panic(fmt.Sprintf("Could not sign with error: %v\n", err.Error()))
	}

	log.Printf("Signed bytes: %v\n", bytes)
}
