package handler

import (
	"fmt"
	"time"

	gethCommon "github.com/onflow/go-ethereum/common"

	"github.com/onflow/flow-go/fvm/evm/types"
	"github.com/onflow/flow-go/model/flow"
)

const (
	BlockHashListCapacity            = 16
	BlockStoreLatestBlockKey         = "LatestBlock"
	BlockStoreLatestBlockProposalKey = "LatestBlockProposal"
	BlockStoreBlockHashesKey         = "LatestBlockHashes"
)

type BlockStore struct {
	backend     types.Backend
	rootAddress flow.Address
}

var _ types.BlockStore = &BlockStore{}

// NewBlockStore constructs a new block store
func NewBlockStore(backend types.Backend, rootAddress flow.Address) *BlockStore {
	return &BlockStore{
		backend:     backend,
		rootAddress: rootAddress,
	}
}

// BlockProposal returns the block proposal to be updated by the handler
func (bs *BlockStore) BlockProposal() (*types.BlockProposal, error) {
	// first fetch it from the storage
	data, err := bs.backend.GetValue(bs.rootAddress[:], []byte(BlockStoreLatestBlockProposalKey))
	if err != nil {
		return nil, err
	}
	if len(data) != 0 {
		return types.NewBlockProposalFromBytes(data)
	}

	// if available construct a new one
	cadenceHeight, err := bs.backend.GetCurrentBlockHeight()
	if err != nil {
		return nil, err
	}

	cadenceBlock, found, err := bs.backend.GetBlockAtHeight(cadenceHeight)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("cadence block not found")
	}

	lastExecutedBlock, err := bs.LatestBlock()
	if err != nil {
		return nil, err
	}

	parentHash, err := lastExecutedBlock.Hash()
	if err != nil {
		return nil, err
	}

	// cadence block timestamp is unix nanoseconds but evm blocks
	// expect timestamps in unix seconds so we convert here
	timestamp := uint64(cadenceBlock.Timestamp / int64(time.Second))

	blockProposal := types.NewBlockProposal(
		parentHash,
		lastExecutedBlock.Height+1,
		timestamp,
		lastExecutedBlock.TotalSupply,
	)
	return blockProposal, nil
}

// UpdateBlockProposal updates the block proposal
func (bs *BlockStore) UpdateBlockProposal(bp *types.BlockProposal) error {
	blockProposalBytes, err := bp.ToBytes()
	if err != nil {
		return types.NewFatalError(err)
	}

	return bs.backend.SetValue(
		bs.rootAddress[:],
		[]byte(BlockStoreLatestBlockProposalKey),
		blockProposalBytes,
	)
}

func (bs *BlockStore) ResetBlockProposal() error {
	return bs.backend.SetValue(
		bs.rootAddress[:],
		[]byte(BlockStoreLatestBlockProposalKey),
		nil,
	)
}

// CommitBlockProposal commits the block proposal to the chain
func (bs *BlockStore) CommitBlockProposal(bp *types.BlockProposal) error {
	bp.PopulateRoots()

	blockBytes, err := bp.Block.ToBytes()
	if err != nil {
		return types.NewFatalError(err)
	}

	err = bs.backend.SetValue(bs.rootAddress[:], []byte(BlockStoreLatestBlockKey), blockBytes)
	if err != nil {
		return err
	}

	hash, err := bp.Block.Hash()
	if err != nil {
		return err
	}

	err = bs.updateBlockHashList(bp.Block.Height, hash)
	if err != nil {
		return err
	}

	err = bs.ResetBlockProposal()
	if err != nil {
		return err
	}

	return nil
}

// LatestBlock returns the latest executed block
func (bs *BlockStore) LatestBlock() (*types.Block, error) {
	data, err := bs.backend.GetValue(bs.rootAddress[:], []byte(BlockStoreLatestBlockKey))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return types.GenesisBlock, nil
	}
	return types.NewBlockFromBytes(data)
}

// BlockHash returns the block hash for the last x blocks
func (bs *BlockStore) BlockHash(height uint64) (gethCommon.Hash, error) {
	bhl, err := bs.getBlockHashList()
	if err != nil {
		return gethCommon.Hash{}, err
	}
	_, hash := bhl.BlockHashByHeight(height)
	return hash, nil
}

func (bs *BlockStore) getBlockHashList() (*types.BlockHashList, error) {
	data, err := bs.backend.GetValue(bs.rootAddress[:], []byte(BlockStoreBlockHashesKey))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		bhl := types.NewBlockHashList(BlockHashListCapacity)
		err = bhl.Push(types.GenesisBlock.Height, types.GenesisBlockHash)
		return bhl, err
	}
	return types.NewBlockHashListFromEncoded(data)
}

func (bs *BlockStore) updateBlockHashList(height uint64, hash gethCommon.Hash) error {
	bhl, err := bs.getBlockHashList()
	if err != nil {
		return err
	}
	err = bhl.Push(height, hash)
	if err != nil {
		return err
	}
	err = bs.backend.SetValue(
		bs.rootAddress[:],
		[]byte(BlockStoreBlockHashesKey),
		bhl.Encode())
	if err != nil {
		return err
	}
	return nil
}
