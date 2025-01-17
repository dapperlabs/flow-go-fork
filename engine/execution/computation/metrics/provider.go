package metrics

import (
	"sync"

	"github.com/rs/zerolog"
)

// provider is responsible for providing the metrics for the rpc endpoint
// it has a circular buffer of metrics for the last N finalized and executed blocks.
type provider struct {
	log zerolog.Logger

	mu sync.RWMutex

	bufferSize               uint
	bufferIndex              uint
	blockHeightAtBufferIndex uint64

	buffer [][]TransactionExecutionMetrics
}

func newProvider(
	log zerolog.Logger,
	bufferSize uint,
	blockHeightAtBufferIndex uint64,
) *provider {
	if bufferSize == 0 {
		panic("buffer size must be greater than zero")
	}

	return &provider{
		log:                      log,
		bufferSize:               bufferSize,
		blockHeightAtBufferIndex: blockHeightAtBufferIndex,
		bufferIndex:              0,
		buffer:                   make([][]TransactionExecutionMetrics, bufferSize),
	}
}

// Push buffers the metrics for the given height.
// The call should ensure height are called in strictly increasing order, otherwise
// metrics for the skipped height will not buffered.
func (p *provider) Push(
	height uint64,
	data []TransactionExecutionMetrics,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if height <= p.blockHeightAtBufferIndex {
		p.log.Warn().
			Uint64("height", height).
			Uint64("blockHeightAtBufferIndex", p.blockHeightAtBufferIndex).
			Msg("received metrics for a block that is older or equal than the most recent block")
		return
	}
	if height > p.blockHeightAtBufferIndex+1 {
		p.log.Warn().
			Uint64("height", height).
			Uint64("blockHeightAtBufferIndex", p.blockHeightAtBufferIndex).
			Msg("received metrics for a block that is not the next block")

		// Fill in the gap with nil
		for i := p.blockHeightAtBufferIndex; i < height-1; i++ {
			p.pushData(nil)
		}
	}

	p.pushData(data)
}

func (p *provider) pushData(data []TransactionExecutionMetrics) {
	p.bufferIndex = (p.bufferIndex + 1) % p.bufferSize
	p.blockHeightAtBufferIndex++
	p.buffer[p.bufferIndex] = data
}

func (p *provider) GetTransactionExecutionMetricsAfter(height uint64) (GetTransactionExecutionMetricsAfterResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	data := make(map[uint64][]TransactionExecutionMetrics)

	if height+1 > p.blockHeightAtBufferIndex {
		return data, nil
	}

	// start index is the lowest block height that is in the buffer
	// missing heights are handled below
	startHeight := uint64(0)
	// assign startHeight with the lowest buffered height
	if p.blockHeightAtBufferIndex > uint64(p.bufferSize) {
		startHeight = p.blockHeightAtBufferIndex - uint64(p.bufferSize)
	}

	// if the starting index is lower than the height we only need to return the data for
	// the blocks that are later than the given height
	if height+1 > startHeight {
		startHeight = height + 1
	}

	for h := startHeight; h <= p.blockHeightAtBufferIndex; h++ {
		// 0 <= diff; because of the bufferSize check above
		diff := uint(p.blockHeightAtBufferIndex - h)

		// 0 <= diff < bufferSize; because of the bufferSize check above
		// we are about to do a modulo operation with p.bufferSize on p.bufferIndex - diff, but diff could
		// be larger than p.bufferIndex, which would result in a negative intermediate value.
		// To avoid this, we add p.bufferSize to diff, which will guarantee that (p.bufferSize + p.bufferIndex - diff)
		// is always positive, but the modulo operation will still return the same index.
		intermediateIndex := p.bufferIndex + p.bufferSize - diff
		index := intermediateIndex % p.bufferSize

		d := p.buffer[index]
		if len(d) == 0 {
			continue
		}

		data[h] = p.buffer[index]
	}

	return data, nil
}
