// This is the network client for communication with nodes
// It is used by nodes to communicate with other nodes and by lite wallets
// to communicate with nodes
package nodeclient

import (
	"bytes"
	"encoding/binary"

	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/gelembjuk/democoin/lib"
	"github.com/gelembjuk/democoin/lib/transaction"
)

type NodeClient struct {
	DataDir     string
	NodeAddress lib.NodeAddr
	Address     string // wallet address
	Logger      *lib.LoggerMan
	NodeNet     *lib.NodeNetwork
}

type ComBlock struct {
	AddrFrom lib.NodeAddr
	Block    []byte
}

// this struct can be used for 2 commands. to get blocks starting from some block to down or to up
type ComGetBlocks struct {
	AddrFrom  lib.NodeAddr
	StartFrom []byte // has of block from which to start and go down or go up in case of Up command
}

// Response of GetBlock request
type ComGetFirstBlocksData struct {
	Blocks [][]byte // lowest block first
	// it is serialised BlockShort structure
	Height int
}

type ComGetData struct {
	AddrFrom lib.NodeAddr
	Type     string
	ID       []byte
}

// New Transaction command. Is used by lite wallets
type ComNewTransaction struct {
	Address string
	TX      *transaction.Transaction
}

// To Request new transaction by wallet.
// Wallet sends address where to send and amount to send
// and own pubkey. Server returns transaction but wihout signatures
type ComRequestTransaction struct {
	PubKey    []byte
	To        string
	Amount    float64
	Signature []byte // to confirm request is from owner of PubKey (TODO)
}

// Response on prepare transaction request. Returns transaction without signs
// and data to sign
type ComRequestTransactionData struct {
	TX         transaction.Transaction
	DataToSign [][]byte
}

// For request to get list of unspent transactions by wallet
type ComGetUnspentTransactions struct {
	Address   string
	LastBlock []byte
}

// Unspent Transaction record
type ComUnspentTransaction struct {
	TXID   []byte
	Vout   int
	Amount float64
	IsBase bool
	From   string
}

// Lit of unspent transactions returned on request
type ComUnspentTransactions struct {
	AddrFrom     lib.NodeAddr
	Transactions []ComUnspentTransaction
	LastBlock    []byte
}

// Request for history of transactions
type ComGetHistoryTransactions struct {
	Address string
}

// Record of transaction in list of history transactions
type ComHistoryTransaction struct {
	IOType bool // In (false) or Out (true)
	TXID   []byte
	Amount float64
	From   string
	To     string
}

// Request for inventory. It can be used to get blocks and transactions from other node
type ComInv struct {
	AddrFrom lib.NodeAddr
	Type     string
	Items    [][]byte
}

// Transaction to send to other node
type ComTx struct {
	AddFrom     lib.NodeAddr
	Transaction []byte // Transaction serialised
}

// Version mesage to other nodes
type ComVersion struct {
	Version    int
	BestHeight int
	AddrFrom   lib.NodeAddr
}

// Check if node address looks fine
func (c *NodeClient) CheckNodeAddress(address lib.NodeAddr) error {
	if address.Port < 1024 {
		return errors.New("Node Address Port has wrong value")
	}
	if address.Port > 65536 {
		return errors.New("Node Address Port has wrong value")
	}
	if address.Host == "" {
		return errors.New("Node Address Host has wrong value")
	}
	return nil
}

// Set currrent node address , to include itin requests to other nodes
func (c *NodeClient) SetNodeAddress(address lib.NodeAddr) {
	c.NodeAddress = address
}

// Send void commant to other node
// It is used by a node to send to itself only when we want to stop a node
// And unblock port listetining
func (c *NodeClient) SendVoid(address lib.NodeAddr) error {
	request := lib.CommandToBytes("viod")

	return c.SendData(address, request)
}

// Send list of nodes addresses to other node
func (c *NodeClient) SendAddrList(address lib.NodeAddr, addresses []lib.NodeAddr) error {
	request, err := c.BuildCommandData("addr", &addresses)

	if err != nil {
		return err
	}

	return c.SendData(address, request)
}

// Send block to other node
func (c *NodeClient) SendBlock(addr lib.NodeAddr, BlockSerialised []byte) error {
	data := ComBlock{c.NodeAddress, BlockSerialised}
	request, err := c.BuildCommandData("block", &data)

	if err != nil {
		return err
	}

	return c.SendData(addr, request)
}

// Send inventory. Blocks hashes or transactions IDs
func (c *NodeClient) SendInv(address lib.NodeAddr, kind string, items [][]byte) error {
	data := ComInv{c.NodeAddress, kind, items}

	request, err := c.BuildCommandData("inv", &data)

	if err != nil {
		return err
	}

	return c.SendData(address, request)
}

// Sedn request to get list of blocks on other node.
func (c *NodeClient) SendGetBlocks(address lib.NodeAddr, startfrom []byte) error {
	data := ComGetBlocks{c.NodeAddress, startfrom}

	request, err := c.BuildCommandData("getblocks", &data)

	if err != nil {
		return err
	}

	return c.SendData(address, request)
}

// Request for blocks but result must be upper from some starting block
func (c *NodeClient) SendGetBlocksUpper(address lib.NodeAddr, startfrom []byte) error {
	data := ComGetBlocks{c.NodeAddress, startfrom}

	request, err := c.BuildCommandData("getblocksup", &data)

	if err != nil {
		return err
	}

	return c.SendData(address, request)
}

// Request for list of first blocks in blockchain.
// This is used by new nodes
// TODO we can use SendGetBlocksUpper and empty hash. This will e same
func (c *NodeClient) SendGetFirstBlocks(address lib.NodeAddr) (*ComGetFirstBlocksData, error) {
	request, err := c.BuildCommandData("getfblocks", nil)

	if err != nil {
		return nil, err
	}
	datapayload := ComGetFirstBlocksData{}

	err = c.SendDataWaitResponse(address, request, &datapayload)

	if err != nil {
		return nil, err
	}

	return &datapayload, nil
}

// Request for a transaction or a block to get full info by ID or Hash
func (c *NodeClient) SendGetData(address lib.NodeAddr, kind string, id []byte) error {
	data := ComGetData{c.NodeAddress, kind, id}

	request, err := c.BuildCommandData("getdata", &data)

	if err != nil {
		return err
	}

	return c.SendData(address, request)
}

// Send Transaction to other node
func (c *NodeClient) SendTx(addr lib.NodeAddr, tnxserialised []byte) error {
	data := ComTx{c.NodeAddress, tnxserialised}
	request, err := c.BuildCommandData("tx", &data)

	if err != nil {
		return err
	}

	return c.SendData(addr, request)
}

// Send own version and blockchain state to other node
func (c *NodeClient) SendVersion(addr lib.NodeAddr, bestHeight int) error {
	data := ComVersion{lib.NodeVersion, bestHeight, c.NodeAddress}

	request, err := c.BuildCommandData("version", &data)

	if err != nil {
		return err
	}

	return c.SendData(addr, request)
}

// Request for history of transaction from a wallet
func (c *NodeClient) SendGetHistory(addr lib.NodeAddr, address string) ([]ComHistoryTransaction, error) {
	data := ComGetHistoryTransactions{address}

	request, err := c.BuildCommandData("gethistory", &data)

	if err != nil {
		return nil, err
	}

	datapayload := []ComHistoryTransaction{}

	err = c.SendDataWaitResponse(addr, request, &datapayload)

	if err != nil {
		return nil, err
	}

	return datapayload, nil
}

// Send new transaction from a wallet to a node
func (c *NodeClient) SendNewTransaction(addr lib.NodeAddr, from string, tx *transaction.Transaction) error {
	data := ComNewTransaction{}
	data.Address = from
	data.TX = tx

	request, err := c.BuildCommandData("txfull", &data)

	if err != nil {
		return nil
	}
	err = c.SendDataWaitResponse(addr, request, nil)

	if err != nil {
		return err
	}
	// no data are returned. only success or not
	return err
}

// Request to prepare new transaction by wallet.
// It returns a transaction without signature.
// Wallet has to sign it and then use SendNewTransaction to send completed transaction
func (c *NodeClient) SendRequestNewTransaction(addr lib.NodeAddr,
	PubKey []byte, to string, amount float64) (*transaction.Transaction, [][]byte, error) {

	data := ComRequestTransaction{}
	data.PubKey = PubKey
	data.To = to
	data.Amount = amount

	request, err := c.BuildCommandData("txrequest", &data)

	if err != nil {
		return nil, nil, err
	}

	datapayload := ComRequestTransactionData{}

	err = c.SendDataWaitResponse(addr, request, &datapayload)

	if err != nil {
		return nil, nil, err
	}

	return &datapayload.TX, datapayload.DataToSign, nil
}

// Request for list of unspent transactions outputs
// It can be used by wallet to see a state of balance
func (c *NodeClient) SendGetUnspent(addr lib.NodeAddr, address string, chaintip []byte) (ComUnspentTransactions, error) {
	data := ComGetUnspentTransactions{address, chaintip}

	request, err := c.BuildCommandData("getunspent", &data)

	datapayload := ComUnspentTransactions{}

	err = c.SendDataWaitResponse(addr, request, &datapayload)

	if err != nil {
		return ComUnspentTransactions{}, err
	}

	return datapayload, nil
}

// Builds a command data. It prepares a slice of bytes from given data
func (c *NodeClient) BuildCommandData(command string, data interface{}) ([]byte, error) {
	var payload []byte
	var err error

	if data != nil {
		payload, err = lib.GobEncode(data)

		if err != nil {
			return nil, err
		}
	} else {
		payload = []byte{}
	}

	payloadlength := uint32(len(payload))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, payloadlength) // convert int to []byte

	request := append(lib.CommandToBytes(command), bs...)
	request = append(request, payload...)

	return request, nil
}

// Sends prepared command to a node. This doesn't wait any response
func (c *NodeClient) SendData(addr lib.NodeAddr, data []byte) error {
	err := c.CheckNodeAddress(addr)

	if err != nil {
		return err
	}

	c.Logger.Trace.Println("Sending data to " + addr.NodeAddrToString())
	conn, err := net.Dial(lib.Protocol, addr.NodeAddrToString())

	if err != nil {
		c.Logger.Error.Println(err.Error())
		c.Logger.Trace.Println("Error: ", err.Error())

		// we can not connect.
		// we could remove this node from known
		// but this is not always good. we need somethign more smart here
		// TODO this needs analysis . if removing of a node is good idea
		//c.NodeNet.RemoveNodeFromKnown(addr)

		return errors.New(fmt.Sprintf("%s is not available", addr.NodeAddrToString()))
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))

	if err != nil {
		c.Logger.Error.Println(err.Error())
		c.Logger.Trace.Println("Error: ", err.Error())
		return err
	}
	return nil
}

// Send data to a node and wait for response
func (c *NodeClient) SendDataWaitResponse(addr lib.NodeAddr, data []byte, datapayload interface{}) error {

	err := c.CheckNodeAddress(addr)

	if err != nil {
		return err
	}

	c.Logger.Trace.Println("Sending data to " + addr.NodeAddrToString() + " and waiting response")

	// connect
	conn, err := net.Dial(lib.Protocol, addr.NodeAddrToString())

	if err != nil {
		c.Logger.Error.Println(err.Error())
		c.Logger.Trace.Println("Error: ", err.Error())

		// we can not connect.
		// we could remove this node from known
		// but this is not always good. we need somethign more smart here
		// TODO this needs analysis . if removing of a node is good idea
		//c.NodeNet.RemoveNodeFromKnown(addr)

		return errors.New(fmt.Sprintf("%s is not available", addr.NodeAddrToString()))
	}
	defer conn.Close()

	// send command bytes
	_, err = io.Copy(conn, bytes.NewReader(data))

	if err != nil {
		c.Logger.Error.Println(err.Error())
		c.Logger.Trace.Println("Error: ", err.Error())
		return err
	}
	// read response
	// read everything
	response, err := ioutil.ReadAll(conn)

	if err != nil {
		c.Logger.Error.Println(err.Error())
		c.Logger.Trace.Println("Response Read Error: ", err.Error())
		return err
	}

	if len(response) == 0 {
		err := errors.New("Received 0 bytes as a response. Expected at least 1 byte")
		c.Logger.Error.Println(err.Error())
		c.Logger.Trace.Println("Response Read Error: ", err.Error())
		return err
	}

	c.Logger.Trace.Printf("Received %d bytes as a response\n", len(response))

	// convert response for provided structure
	var buff bytes.Buffer
	buff.Write(response[1:])
	dec := gob.NewDecoder(&buff)

	if response[0] != 1 {
		// fail

		var payload string

		err := dec.Decode(&payload)

		if err != nil {
			return err
		}

		return errors.New(payload)
	}

	if datapayload != nil {
		err = dec.Decode(datapayload)

		if err != nil {
			return err
		}
	}

	return nil
}
