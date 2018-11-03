package transactions

import (
	"crypto/ecdsa"

	"github.com/NlaakStudios/democoin/node/structures/transaction"

	"github.com/NlaakStudios/democoin/lib/wallet"
	"github.com/NlaakStudios/democoin/node/structures"
)

type UnApprovedTransactionCallbackInterface func(txhash, txstr string) error
type UnspentTransactionOutputCallbackInterface func(fromaddr string, value float64, txID []byte, output int, isbase bool) error

type TransactionsManagerInterface interface {
	GetAddressBalance(address string) (wallet.WalletBalance, error)
	GetUnapprovedCount() (int, error)
	GetUnspentCount() (int, error)
	GetUnapprovedTransactionsForNewBlock(number int) ([]*transaction.Transaction, error)
	GetIfExists(txid []byte) (*transaction.Transaction, error)
	GetIfUnapprovedExists(txid []byte) (*transaction.Transaction, error)

	VerifyTransaction(tx *transaction.Transaction, prevtxs []*transaction.Transaction, tip []byte) (bool, error)

	ForEachUnspentOutput(address string, callback UnspentTransactionOutputCallbackInterface) error
	ForEachUnapprovedTransaction(callback UnApprovedTransactionCallbackInterface) (int, error)

	// Create transaction methods
	CreateTransaction(PubKey []byte, privKey ecdsa.PrivateKey, to string, amount float64) (*transaction.Transaction, error)
	ReceivedNewTransaction(tx *transaction.Transaction) error
	ReceivedNewTransactionData(txBytes []byte, Signatures [][]byte) (*transaction.Transaction, error)
	PrepareNewTransaction(PubKey []byte, to string, amount float64) ([]byte, [][]byte, error)

	// new block was created in blockchain DB. It must not be on top of primary blockchain
	BlockAdded(block *structures.Block, ontopofchain bool) error
	// block was removed from blockchain DB from top
	BlockRemoved(block *structures.Block) error
	// block was not in primary chain and now is
	BlockAddedToPrimaryChain(block *structures.Block) error
	// block was in primary chain and now is not
	BlockRemovedFromPrimaryChain(block *structures.Block) error

	CancelTransaction(txID []byte) error
	ReindexData() (map[string]int, error)
	CleanUnapprovedCache() error
}
