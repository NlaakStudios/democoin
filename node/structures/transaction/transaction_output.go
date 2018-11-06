package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"strings"

	"github.com/NlaakStudios/democoin/lib/utils"
)

// TXOutput represents a transaction output.
// An output contains instructions for sending bitcoins. Value is the number of Satoshi
// (1 BTC = 100,000,000 Satoshi) that this output will be worth when claimed. ScriptPubKey
// is the second half of a script (discussed later). There can be more than one output, and
// they share the combined value of the inputs. Because each output from one transaction can
// only ever be referenced once by an input of a subsequent transaction, the entire combined
// input value needs to be sent in an output if you don't want to lose it. If the input is
// worth 50 BTC but you only want to send 25 BTC, Bitcoin will create two outputs worth 25 BTC:
// one to the destination, and one back to you (known as "change", though you send it to yourself).
// Any input bitcoins not redeemed in an output is considered a transaction fee; whoever generates
// the block can claim it by inserting it into the coinbase transaction of that block.
type TXOutput struct {
	Value      float64
	PubKeyHash []byte
}

// Simplified output format. To use externally
// It has all info in human readable format
// this can be used to display info about outputs wihout references to transaction object
type TXOutputIndependent struct {
	Value          float64
	DestPubKeyHash []byte
	SendPubKeyHash []byte
	TXID           []byte
	OIndex         int
	IsBase         bool
	BlockHash      []byte
}

type TXOutputIndependentList []TXOutputIndependent

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// Same as IsLockedWithKey but for simpler structure
func (out *TXOutputIndependent) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.DestPubKeyHash, pubKeyHash) == 0
}

// build independent transaction from normal output
func (out *TXOutputIndependent) LoadFromSimple(sout TXOutput, txid []byte, ind int, sender []byte, iscoinbase bool, blockHash []byte) {
	out.OIndex = ind
	out.DestPubKeyHash = sout.PubKeyHash
	out.SendPubKeyHash = sender
	out.Value = sout.Value
	out.TXID = txid
	out.IsBase = iscoinbase
	out.BlockHash = blockHash
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value float64, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// TXOutputs collects TXOutput
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}

func (output TXOutput) String() string {
	lines := []string{}

	lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
	lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))

	return strings.Join(lines, "\n")
}

func (output TXOutput) ToBytes() ([]byte, error) {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, output.Value)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.BigEndian, output.PubKeyHash)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (a TXOutputIndependentList) Len() int           { return len(a) }
func (a TXOutputIndependentList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TXOutputIndependentList) Less(i, j int) bool { return a[i].Value < a[j].Value }
