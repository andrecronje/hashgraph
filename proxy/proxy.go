package proxy

import "github.com/andrecronje/hashgraph/hashgraph"

type AppProxy interface {
	SubmitCh() chan []byte
	CommitBlock(block hashgraph.Block) ([]byte, error)
}

type HashgraphProxy interface {
	CommitCh() chan hashgraph.Block
	SubmitTx(tx []byte) error
}
