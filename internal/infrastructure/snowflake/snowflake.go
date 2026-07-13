package snowflake

import (
	"errors"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	epoch int64 = 1704067200000

	nodeBits     uint8 = 10
	sequenceBits uint8 = 12

	maxNode int64 = -1 ^ (-1 << nodeBits)
	maxSeq  int64 = -1 ^ (-1 << sequenceBits)

	timeShift uint8 = nodeBits + sequenceBits
	nodeShift uint8 = sequenceBits
)

type Node struct {
	state atomic.Uint64
	node  int64
}

func NewSnowflakeNode(nodeID int64) (*Node, error) {
	if nodeID < 0 || nodeID > maxNode {
		return nil, errors.New("snowflake: node id must be between 0 and 1023")
	}
	return &Node{node: nodeID}, nil
}

func (n *Node) Generate() int64 {
	for {
		// 1. Read current state
		currentState := n.state.Load()
		lastTime := int64(currentState >> sequenceBits)
		sequence := int64(currentState) & maxSeq

		now := time.Now().UnixMilli()

		if now < lastTime {
			panic("snowflake: clock moved backwards")
			// Clock went backwards. Yield and retry (or panic/return error).
			// runtime.Gosched()
			// continue
		}

		if now == lastTime {
			sequence = (sequence + 1) & maxSeq
			if sequence == 0 {
				// Millisecond exhausted, yield and try again
				runtime.Gosched()
				continue
			}
		} else {
			sequence = 0
		}

		// 2. Pack the new state
		newState := (uint64(now) << sequenceBits) | uint64(sequence)

		// 3. Attempt to swap. If successful, calculate and return the ID.
		// If it fails, another goroutine beat us to it, so the loop restarts.
		if n.state.CompareAndSwap(currentState, newState) {
			return ((now - epoch) << timeShift) | (n.node << nodeShift) | sequence
		}
	}
}

var defaultNode *Node

func Init(nodeID int) error {
	node, err := NewSnowflakeNode(int64(nodeID))
	if err != nil {
		return err
	}
	defaultNode = node
	return nil
}

func Generate() int64 {
	if defaultNode == nil {
		panic("snowflake: Generate called before Init")
	}
	return defaultNode.Generate()
}
