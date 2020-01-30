package types

import (
	"time"

	"github.com/dapperlabs/flow-go/model/flow"
)

type Block struct {
	// specified
	View        uint64
	QC          *QuorumCertificate
	PayloadHash []byte

	// configed
	Height  uint64
	ChainID string

	BlockID flow.Identifier

	// autogenerated
	Timestamp time.Time
}

// BlockID returns the Merkle Root Hash of the Block, which is computed from (View, QC, PayloadHash)

// NewBlock creates an instance of Block
func NewBlock(view uint64, qc *QuorumCertificate, payloadHash []byte, height uint64, chainID string) *Block {

	t := time.Now()

	return &Block{
		View:        view,
		QC:          qc,
		PayloadHash: payloadHash,
		Height:      height,
		ChainID:     chainID,
		Timestamp:   t,
	}
}

// ToVote converts a UnsignedBlockProposal to a UnsignedVote
func (b Block) ToVote() *UnsignedVote {
	return NewUnsignedVote(b.View, b.BlockID)
}
