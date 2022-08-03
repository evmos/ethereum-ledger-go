# ethereum-ledger-go

This package is a general-purpose Ethereum Ledger library adapted from go-ethereum for applications written in Go. It is currently in Beta and experimental, so please use it at your own risk.

## Usage

### Initialize Wallet
```
import (
 	ethLedger "github.com/evmos/ethereum-ledger-go"
	"github.com/evmos/ethereum-ledger-go/accounts"
)

ledger, err := ethLedger.New()

wallet := ledger.Wallets()[0]
err = wallet.Open("")
path := accounts.DefaultBaseDerivationPath
account, err := wallet.Derive(path, true)
```

### Signing Transactions
```
import ethLedger "github.com/evmos/ethereum-ledger-go"

tx := ethLedger.CreateTx(
  3,               // Nonce
  big.NewInt(10),  // GasPrice
  10,              // Gas
  addr,            // To
  big.NewInt(10),  // Value
  make([]byte, 0), // Data
)

// Initialize Wallet
sigBytes, err := wallet.SignTx(
  account, // Wallet Account
  tx, // Tx
  big.NewInt(0) // Chain ID
)
```

### Signing Typed Data
```
import "github.com/ethereum/go-ethereum/signer/core/apitypes"

// Create typedData object that conforms to apitypes.TypedData
// (see tests/integration_test.go for a complete example)
sigBytes, err := wallet.SignTypedData(account, typedData)
```

