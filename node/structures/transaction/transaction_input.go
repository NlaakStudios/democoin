package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)
import "github.com/NlaakStudios/democoin/lib/utils"

// TXInput represents a transaction input
//
// An input is a reference to an output from a previous transaction. Multiple inputs are often listed
// in a transaction. All of the new transaction's input values (that is, the total coin value of the
// previous outputs referenced by the new transaction's inputs) are added up, and the total (less any
// transaction fee) is completely used by the outputs of the new transaction. Previous tx is a hash of
// a previous transaction. Index is the specific output in the referenced transaction. ScriptSig is the
// first half of a script (discussed in more detail later).
//
// The script contains two components, a signature and a public key. The public key must match the hash
// given in the script of the redeemed output. The public key is used to verify the redeemers signature,
// which is the second component. More precisely, the second component is an ECDSA signature over a hash
// of a simplified version of the transaction. It, combined with the public key, proves the transaction
// was created by the real owner of the address in question. Various flags define how the transaction is
// simplified and can be used to create different types of payment.
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte // this is the wallet who spends transaction
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash, _ := utils.HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (input TXInput) String() string {
	lines := []string{}

	lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
	lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
	lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
	lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))

	return strings.Join(lines, "\n")
}

func (input TXInput) ToBytes() ([]byte, error) {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, input.Txid)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.BigEndian, int32(input.Vout))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.BigEndian, input.Signature)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.BigEndian, input.PubKey)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
