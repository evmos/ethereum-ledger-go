package ledger

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	ethLedger "github.com/evmos/ethereum-ledger-go"
	"github.com/evmos/ethereum-ledger-go/accounts"
)

// Test Mnemonic:
// glow spread dentist swamp people siren hint muscle first sausage castle metal
// cycle abandon accident logic again around mix dial knee organ episode usual
// (24 words)

func initWallet(t *testing.T, path accounts.DerivationPath) (accounts.Wallet, accounts.Account) {
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

	account, err := wallet.Derive(path, true)
	if err != nil {
		panic("Could not derive account")
	}

	return wallet, account
}

func createTypedDataPayload(message map[string]interface{}) apitypes.TypedData {
	const primaryType = "Mail"

	domain := apitypes.TypedDataDomain{
		Name:              "Ether Mail",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(1),
		VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		Salt:              "",
	}

	domainTypes := apitypes.Types{
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

	return apitypes.TypedData{
		Types:       domainTypes,
		PrimaryType: primaryType,
		Domain:      domain,
		Message:     message,
	}
}

// Test transaction is generated correctly using CreateTx
func TestSanityCreateTx(t *testing.T) {
	var addr = "0x3535353535353535353535353535353535353535"

	tx := ethLedger.CreateTx(
		3,               // Nonce
		big.NewInt(10),  // GasPrice
		10,              // Gas
		addr,            // To
		big.NewInt(10),  // Value
		make([]byte, 0), // Data
	)

	if tx.Nonce() != 3 {
		t.Errorf("Invalid nonce received")
	}

	if !reflect.DeepEqual(tx.GasPrice(), big.NewInt(10)) {
		t.Errorf("Invalid gas price received")
	}

	if tx.Gas() != 10 {
		t.Errorf("Invalid gas received")
	}

	addrBytes, err := hex.DecodeString(strings.TrimPrefix(addr, "0x"))
	if err != nil {
		t.Errorf("Could not convert address to bytes")
	}

	if !reflect.DeepEqual(tx.To()[:], addrBytes) {
		t.Errorf("Invalid 'to' received (bytes)")
	}

	if tx.To().Hex() != addr {
		t.Errorf("Invalid 'to' received (hex)")
	}

	if !reflect.DeepEqual(*tx.Value(), *big.NewInt(10)) {
		t.Errorf("Invalid value received")
	}

	if !reflect.DeepEqual(tx.Data(), make([]byte, 0)) {
		t.Errorf("Invalid data received")
	}
}

func TestInitWallet(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

	t.Logf("Account: %v\n", account.Address.Hex())
	if account.Address.Hex() != "0xbcf6368dF2C2999893064aDe8C4a4b1b6d3C077B" {
		t.Errorf("Invalid address for account")
	}

	t.Logf("Public Key: %v\n", account.PublicKey.Hex())
	if account.PublicKey.Hex() != "0x045f53cbc346997423fe843e2ee6d24fd7832211000a65975ba81d53c87ad1e5c863a5adb3cb919014903f13a68c9a4682b56ff5df3db888a2cbc3dc8fae1ec0fb" {
		t.Errorf("Invalid public key for account")
	}
}

func TestInvalidAccount(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

	account.Address = common.HexToAddress("0x3535353535353535353535353535353535353535")

	var sendAddr = "0x3636363636363636363636363636363636363636"
	tx := ethLedger.CreateTx(
		3, big.NewInt(10), 10, sendAddr, big.NewInt(10), make([]byte, 0),
	)
	_, err := wallet.SignTx(account, tx, big.NewInt(0))
	if err == nil {
		t.Errorf("Expected error on signing with invalid account")
	}
}

// Test deriving an account with path "m/44'/60'/0'/0/1"
func TestAlternateDerivation(t *testing.T) {
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/1")
	if err != nil {
		panic("Could not parse derivation path")
	}

	wallet, account := initWallet(t, path)
	defer wallet.Close()

	if account.PublicKey.Hex() != "0x044a5236e77ab81e094d7c6cfeac06d2e93fec455d01c7f80e22c592a89b44acebe99c2450425a184e5382362d5c52f5d996f12e73ccfb7694227f31b501e36ed7" {
		t.Errorf("Account public key did not match expected %v\n", account.PublicKey.Hex())
	}

	if account.Address.Hex() != "0x5737A62f68305D0D94564BED79f0d5347EdC060e" {
		t.Errorf("Account address did not match expected")
	}
}

func TestLedgerSignTx1(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

	var addr = "0x3535353535353535353535353535353535353535"

	tx := ethLedger.CreateTx(
		3, big.NewInt(10), 10, addr, big.NewInt(10), make([]byte, 0),
	)

	sigBytes, err := wallet.SignTx(account, tx, big.NewInt(0))
	if err != nil {
		println(err.Error())
		panic("Could not sign data")
	}

	sigHex := hex.EncodeToString(sigBytes)
	t.Logf("Signed bytes: %v\n", sigHex)

	// Test against signature generated using ethers.js
	if sigHex != "f85d030a0a9435353535353535353535353535353535353535350a801ca02e0b1b0ed24cd450488eb783e6c64ab0f1d681641970aef062434513731e829ca0721e7b6feedc989a8b114f3f622d5a525095b893b8ce81059e682f7333be3508" {
		t.Errorf("Invalid signature received")
	}
}

func TestLedgerSignTx2(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

	var addr = "0x4646464646464646464646464646464646464646"

	tx := ethLedger.CreateTx(
		8, big.NewInt(5), 50, addr, big.NewInt(70), []byte{4, 6, 8, 10},
	)

	sigBytes, err := wallet.SignTx(account, tx, big.NewInt(0))
	if err != nil {
		println(err.Error())
		panic("Could not sign data")
	}

	sigHex := hex.EncodeToString(sigBytes)
	t.Logf("Signed bytes: %v\n", sigHex)

	// Test against signature generated using ethers.js
	if sigHex != "f86108053294464646464646464646464646464646464646464646840406080a1ba0a2120857c6a2f9a2cabe59845b4e3925daf5d13394de52f87f2942f2ba4f9de3a031ecb1178393d2b6b4220eda7876f9a719498f4269f6444dfc5c270baec070cc" {
		t.Errorf("Invalid signature received")
	}
}

// Test signing a transaction with a different ChainID
func TestLedgerSignTx3(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

	var addr = "0x4646464646464646464646464646464646464646"

	tx := ethLedger.CreateTx(
		8, big.NewInt(5), 50, addr, big.NewInt(70), []byte{4, 6, 8, 10},
	)

	sigBytes, err := wallet.SignTx(account, tx, big.NewInt(1))
	if err != nil {
		println(err.Error())
		panic("Could not sign data")
	}

	sigHex := hex.EncodeToString(sigBytes)
	t.Logf("Signed bytes: %v\n", sigHex)

	// Test against signature generated using ethers.js
	if sigHex != "f86108053294464646464646464646464646464646464646464646840406080a26a0da90a513f9ecb1726bc0e77a88f3b1def3f468e41fd3d3f182703085e1b48feca017743f2d3cb4e4b5090e7bd70357a6c4d1b5bfd79531f47669b9a085b530ef37" {
		t.Errorf("Invalid signature received")
	}
}

// Test signing a transaction from an alternate (not default) derivation path
func TestLedgerSignTx4(t *testing.T) {
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/1")
	if err != nil {
		panic("Could not parse derivation path")
	}

	wallet, account := initWallet(t, path)
	defer wallet.Close()

	var sendAddr = "0x4646464646464646464646464646464646464646"
	tx := ethLedger.CreateTx(
		8, big.NewInt(5), 50, sendAddr, big.NewInt(70), []byte{4, 6, 8, 10},
	)

	sigBytes, err := wallet.SignTx(account, tx, big.NewInt(0))
	if err != nil {
		t.Errorf("Could not sign data")
	}

	sigHex := hex.EncodeToString(sigBytes)

	// Test against signature generated using ethers.js
	if sigHex != "f86108053294464646464646464646464646464646464646464646840406080a1ba06eaf33b72493638546d6294b39694c5e420fa55de74af6f1c87739da2582e231a016568d498cebfaa7ed80db99b7be76c31f5b78285cfa1794be3141dfae838f16" {
		t.Errorf("Invalid signature received: %v\n", sigHex)
	}
}

func TestLedgerSignTyped1(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

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

	typedData := createTypedDataPayload(messageStandard)

	sigBytes, err := wallet.SignTypedData(account, typedData)
	if err != nil {
		panic(fmt.Sprintf("Could not sign with error: %v\n", err.Error()))
	}

	sigHex := hex.EncodeToString(sigBytes)
	t.Logf("Signature: %v\n", sigHex)
	if sigHex != "fb35835539608d309ee5ee4b3dfbbc8cb4b591d7e8c9c745473848cbe1e13b037278226e2d6962b3b19145147314d9872ff853437e3ebd654d44aace09128acd1c" {
		t.Errorf("Invalid signature received")
	}
}

func TestLedgerSignTyped2(t *testing.T) {
	wallet, account := initWallet(t, accounts.DefaultBaseDerivationPath)
	defer wallet.Close()

	messageStandard := map[string]interface{}{
		"from": map[string]interface{}{
			"name":   "Charlie",
			"wallet": "0x1CfC9d8357cBE15E08Bb7084073B7E4ef790B625",
		},
		"to": map[string]interface{}{
			"name":   "Delta",
			"wallet": "0x53Fe71EDEFdF942dDE10834ed4d443A6df391F64",
		},
		"contents": "Message from Charlie to Delta!",
	}

	typedData := createTypedDataPayload(messageStandard)

	sigBytes, err := wallet.SignTypedData(account, typedData)
	if err != nil {
		panic(fmt.Sprintf("Could not sign with error: %v\n", err.Error()))
	}

	sigHex := hex.EncodeToString(sigBytes)
	t.Logf("Signature: %v\n", sigHex)
	if sigHex != "d929a56d69a98f3e491828fbd1555e66ddde17c8928a69704e710a9c34db1ab80314ffccf7014be6c8f819ca9c9603d59aad58cddaa1e6f43c7f66a6b9183c681c" {
		t.Errorf("Invalid signature received")
	}
}

// Test signing TypedData from a non-default derivation path
func TestLedgerSignTyped3(t *testing.T) {
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/1")
	if err != nil {
		panic("Could not parse derivation path")
	}

	wallet, account := initWallet(t, path)
	defer wallet.Close()

	messageStandard := map[string]interface{}{
		"from": map[string]interface{}{
			"name":   "Charlie",
			"wallet": "0x1CfC9d8357cBE15E08Bb7084073B7E4ef790B625",
		},
		"to": map[string]interface{}{
			"name":   "Delta",
			"wallet": "0x53Fe71EDEFdF942dDE10834ed4d443A6df391F64",
		},
		"contents": "Message from Charlie to Delta!",
	}

	typedData := createTypedDataPayload(messageStandard)

	sigBytes, err := wallet.SignTypedData(account, typedData)
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not sign with error: %v\n", err.Error()))
	}

	sigHex := hex.EncodeToString(sigBytes)
	t.Logf("Signature: %v\n", sigHex)
	if sigHex != "76984ce659f841975bdab7762ed9cb3c936791d1dcded3c0554147fca7accfdc543313669dcda04350990884e9e10c382fb20b722c123409a97c42ef6df617ca1c" {
		t.Errorf("Invalid signature received")
	}
}
