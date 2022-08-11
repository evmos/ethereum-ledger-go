# ethereum-ledger-go

This package is a general-purpose Ethereum Ledger library adapted from [go-ethereum](https://github.com/ethereum/go-ethereum) for applications written in Go.

## Usage

### Initialize Wallet
```
import (
 	ethLedger "github.com/evmos/ethereum-ledger-go"
	"github.com/evmos/ethereum-ledger-go/accounts"
)

ledger, err := ethLedger.New()

wallet := ledger.Wallets()[0]     // Use first USB device detected
err = wallet.Open("")

path := accounts.DefaultBaseDerivationPath      // m/44'/60'/0'/0/0
account, err := wallet.Derive(path, true)       // Boolean indicates whether the account should be cached on the wallet
```

### Sign Transactions
```
import ethLedger "github.com/evmos/ethereum-ledger-go"
import "math/big"

tx := ethLedger.CreateTx(
  3,                   // Nonce
  big.NewInt(10),      // GasPrice
  10,                  // Gas
  addr,                // To
  big.NewInt(10),      // Value
  make([]byte, 0),     // Data
)

// Initialize Wallet
sigBytes, err := wallet.SignTx(
  account,              // Wallet Account
  tx,                   // Tx
  big.NewInt(0)         // Chain ID
)
```

### Sign Typed Data
```
import "github.com/ethereum/go-ethereum/signer/core/apitypes"

// First, create a typedData object that conforms to apitypes.TypedData
// (see tests/integration_test.go for a complete example)
sigBytes, err := wallet.SignTypedData(
  account,          // Wallet Account
  typedData         // EIP-712 conformant Typed Data
)
```

## Notes
- This library currently does not support [personal signing](https://eips.ethereum.org/EIPS/eip-191)
