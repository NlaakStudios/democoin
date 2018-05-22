package server

import (
	"errors"

	"github.com/gelembjuk/democoin/lib/net"
	"github.com/gelembjuk/democoin/lib/utils"
)

type NodeTransit struct {
	Blocks        map[string][][]byte
	MaxKnownHeigh int
	Logger        *utils.LoggerMan
}

func (t *NodeTransit) Init(l *utils.LoggerMan) error {
	t.Logger = l
	t.Blocks = make(map[string][][]byte)

	return nil
}
func (t *NodeTransit) AddBlocks(fromaddr net.NodeAddr, blocks [][]byte) error {
	key := fromaddr.NodeAddrToString()

	_, ok := t.Blocks[key]

	if !ok {
		t.Blocks[key] = blocks
	} else {
		t.Blocks[key] = append(t.Blocks[key], blocks...)
	}

	return nil
}

func (t *NodeTransit) CleanBlocks(fromaddr net.NodeAddr) {
	key := fromaddr.NodeAddrToString()

	if _, ok := t.Blocks[key]; ok {
		delete(t.Blocks, key)
	}
}

func (t *NodeTransit) GetBlocksCount(fromaddr net.NodeAddr) int {
	if _, ok := t.Blocks[fromaddr.NodeAddrToString()]; ok {
		return len(t.Blocks[fromaddr.NodeAddrToString()])
	}
	return 0
}

func (t *NodeTransit) ShiftNextBlock(fromaddr net.NodeAddr) ([]byte, error) {
	key := fromaddr.NodeAddrToString()

	if _, ok := t.Blocks[key]; ok {
		data := t.Blocks[key][0][:]
		t.Blocks[key] = t.Blocks[key][1:]

		if len(t.Blocks[key]) == 0 {
			delete(t.Blocks, key)
		}

		return data, nil
	}

	return nil, errors.New("The address is not in blocks transit")
}