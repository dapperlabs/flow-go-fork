package ingest_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

	"github.com/dapperlabs/flow-go/consensus/hotstuff/model"
	"github.com/dapperlabs/flow-go/engine"
	"github.com/dapperlabs/flow-go/engine/testutil"
	"github.com/dapperlabs/flow-go/engine/testutil/mock"
	"github.com/dapperlabs/flow-go/engine/verification"
	"github.com/dapperlabs/flow-go/engine/verification/test"
	chmodel "github.com/dapperlabs/flow-go/model/chunks"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/model/messages"
	network "github.com/dapperlabs/flow-go/network/mock"
	"github.com/dapperlabs/flow-go/network/stub"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

// testConcurrency evaluates behavior of verification node against:
// - ingest engine receives concurrent receipts from different sources
// - not all chunks of the receipts are assigned to the ingest engine
// - for each assigned chunk ingest engine emits a single result approval to verify engine only once
// (even in presence of duplication)
// - also the test stages to drop the first request on each collection to evaluate the retrial
// - also the test stages to drop the first request on each chunk data pack to evaluate the retrial
func TestConcurrency(t *testing.T) {
	var mu sync.Mutex
	testcases := []struct {
		erCount, // number of execution receipts
		senderCount, // number of (concurrent) senders for each execution receipt
		chunksNum int // number of chunks in each execution receipt
		lightIngest bool // indicates if light ingest engine should replace the original one
	}{
		{
			erCount:     1,
			senderCount: 1,
			chunksNum:   2,
			lightIngest: true,
		},
		{
			erCount:     1,
			senderCount: 5,
			chunksNum:   2,
			lightIngest: true,
		},
		{
			erCount:     5,
			senderCount: 1,
			chunksNum:   2,
			lightIngest: true,
		},
		{
			erCount:     5,
			senderCount: 5,
			chunksNum:   2,
			lightIngest: true,
		},
		{
			erCount:     1,
			senderCount: 1,
			chunksNum:   10, // choosing a higher number makes the test longer and longer timeout needed
			lightIngest: true,
		},
		{
			erCount:     2,
			senderCount: 5,
			chunksNum:   4,
			lightIngest: true,
		},
		{
			erCount:     1,
			senderCount: 1,
			chunksNum:   2,
			lightIngest: true,
		},
		{
			erCount:     1,
			senderCount: 5,
			chunksNum:   2,
			lightIngest: false,
		},
		{
			erCount:     5,
			senderCount: 1,
			chunksNum:   2,
			lightIngest: false,
		},
		{
			erCount:     5,
			senderCount: 5,
			chunksNum:   2,
			lightIngest: false,
		},
		{
			erCount:     1,
			senderCount: 1,
			chunksNum:   10, // choosing a higher number makes the test longer and longer timeout needed
			lightIngest: false,
		},
		{
			erCount:     2,
			senderCount: 5,
			chunksNum:   4,
			lightIngest: false,
		},
	}

	for _, tc := range testcases {

		t.Run(fmt.Sprintf("%d-ers/%d-senders/%d-chunks/%t-lightIngest",
			tc.erCount, tc.senderCount, tc.chunksNum, tc.lightIngest), func(t *testing.T) {
			mu.Lock()
			defer mu.Unlock()
			testConcurrency(t, tc.erCount, tc.senderCount, tc.chunksNum, tc.lightIngest)

		})
	}
}

func testConcurrency(t *testing.T, erCount, senderCount, chunksNum int, lightIngest bool) {
	log := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
	// to demarcate the logs
	log.Debug().
		Int("execution_receipt_count", erCount).
		Int("sender_count", senderCount).
		Int("chunks_num", chunksNum).
		Bool("light_ingest", lightIngest).
		Msg("TestConcurrency started")
	hub := stub.NewNetworkHub()

	// ingest engine parameters
	// parameters added based on following issue:
	requestInterval := uint(1000)
	failureThreshold := uint(2)

	// creates test id for each role
	colID := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
	conID := unittest.IdentityFixture(unittest.WithRole(flow.RoleConsensus))
	exeID := unittest.IdentityFixture(unittest.WithRole(flow.RoleExecution))
	verID := unittest.IdentityFixture(unittest.WithRole(flow.RoleVerification))

	identities := flow.IdentityList{colID, conID, exeID, verID}

	// new chunk assignment
	assignment := chmodel.NewAssignment()

	// create `erCount` ER fixtures that will be concurrently delivered
	ers := make([]verification.CompleteExecutionResult, 0)
	// list of assigned chunks to the verifier node
	vChunks := make([]*verification.VerifiableChunk, 0)
	// a counter to assign chunks every other one, so to check if
	// ingest only sends the assigned chunks to verifier

	for i := 0; i < erCount; i++ {
		er := test.LightExecutionResultFixture(chunksNum)
		ers = append(ers, er)
		// assigns all chunks to the verifier node
		for j, chunk := range er.Receipt.ExecutionResult.Chunks {
			assignment.Add(chunk, []flow.Identifier{verID.NodeID})

			var endState flow.StateCommitment
			// last chunk
			if int(chunk.Index) == len(er.Receipt.ExecutionResult.Chunks)-1 {
				endState = er.Receipt.ExecutionResult.FinalStateCommit
			} else {
				endState = er.Receipt.ExecutionResult.Chunks[j+1].StartState
			}

			vc := &verification.VerifiableChunk{
				ChunkIndex:    chunk.Index,
				EndState:      endState,
				Block:         er.Block,
				Receipt:       er.Receipt,
				Collection:    er.Collections[chunk.Index],
				ChunkDataPack: er.ChunkDataPacks[chunk.Index],
			}
			vChunks = append(vChunks, vc)
		}
	}

	// set up mock verifier engine that asserts each receipt is submitted
	// to the verifier exactly once.
	verifierEng, verifierEngWG := test.SetupMockVerifierEng(t, vChunks)
	assigner := test.NewMockAssigner(verID.NodeID)
	verNode := testutil.VerificationNode(t, hub, verID, identities, assigner, requestInterval, failureThreshold,
		lightIngest,
		testutil.WithVerifierEngine(verifierEng))

	// starts the ingest engine
	if lightIngest {
		<-verNode.LightIngestEngine.Ready()
	} else {
		<-verNode.IngestEngine.Ready()
	}

	collections := make([]*flow.Collection, 0)
	for _, completeER := range ers {
		collections = append(collections, completeER.Collections...)
	}

	colNode := testutil.GenericNode(t, hub, colID, identities)
	setupMockCollectionNode(t, colNode, verID.NodeID, collections)

	// mock the execution node with a generic node and mocked engine
	// to handle requests for chunk state
	exeNode := testutil.GenericNode(t, hub, exeID, identities)
	setupMockExeNode(t, exeNode, verID.NodeID, ers)

	// creates a network instance for verification node
	// and sets it in continuous delivery mode
	verNet, ok := hub.GetNetwork(verID.NodeID)
	assert.True(t, ok)
	verNet.StartConDev(requestInterval, true)

	// the wait group tracks goroutines for each ER sending it to VER
	var senderWG sync.WaitGroup
	senderWG.Add(erCount * senderCount)

	var blockStorageLock sync.Mutex

	for _, completeER := range ers {

		// spin up `senderCount` sender goroutines to mimic receiving
		// the same resource multiple times
		for i := 0; i < senderCount; i++ {
			go func(j int, id flow.Identifier, block *flow.Block, receipt *flow.ExecutionReceipt) {

				sendBlock := func() {
					// adds the block to the storage of the node
					// Note: this is done by the follower
					// this block should be done in a thread-safe way
					blockStorageLock.Lock()
					// we don't check for error as it definitely returns error when we
					// have duplicate blocks, however, this is not the concern for this test
					_ = verNode.Blocks.Store(block)
					blockStorageLock.Unlock()

					// casts block into a Hotstuff block for notifier
					hotstuffBlock := &model.Block{
						BlockID:     block.ID(),
						View:        block.Header.View,
						ProposerID:  block.Header.ProposerID,
						QC:          nil,
						PayloadHash: block.Header.PayloadHash,
						Timestamp:   block.Header.Timestamp,
					}
					// starts the ingest engine
					if lightIngest {
						verNode.LightIngestEngine.OnFinalizedBlock(hotstuffBlock)
					} else {
						verNode.IngestEngine.OnFinalizedBlock(hotstuffBlock)
					}

				}

				sendReceipt := func() {
					if lightIngest {
						err := verNode.LightIngestEngine.Process(exeID.NodeID, receipt)
						require.NoError(t, err)
					} else {
						err := verNode.IngestEngine.Process(exeID.NodeID, receipt)
						require.NoError(t, err)
					}
				}

				switch j % 2 {
				case 0:
					// block then receipt
					sendBlock()
					// allow another goroutine to run before sending receipt
					time.Sleep(time.Nanosecond)
					sendReceipt()
				case 1:
					// receipt then block
					sendReceipt()
					// allow another goroutine to run before sending block
					time.Sleep(time.Nanosecond)
					sendBlock()
				}

				senderWG.Done()
			}(i, completeER.Receipt.ExecutionResult.ID(), completeER.Block, completeER.Receipt)
		}
	}

	// wait for all ERs to be sent to VER
	unittest.RequireReturnsBefore(t, senderWG.Wait, time.Duration(senderCount*chunksNum*erCount*5)*time.Second)
	unittest.RequireReturnsBefore(t, verifierEngWG.Wait, time.Duration(senderCount*chunksNum*erCount*5)*time.Second)

	// stops ingest engine of verification node
	// Note: this should be done prior to any evaluation to make sure that
	// the checkTrackers method of Ingest engine is done working.
	// starts the ingest engine
	if lightIngest {
		<-verNode.LightIngestEngine.Done()
	} else {
		<-verNode.IngestEngine.Done()
	}

	// stops the network continuous delivery mode
	verNet.StopConDev()

	for _, c := range vChunks {
		if test.IsAssigned(c.ChunkIndex) {
			// assigned chunks should have their result to be added to ingested results mempool
			assert.True(t, verNode.IngestedResultIDs.Has(c.Receipt.ExecutionResult.ID()))
		}
	}

	exeNode.Done()
	colNode.Done()
	verNode.Done()

	// to demarcate the logs
	log.Debug().
		Int("execution_receipt_count", erCount).
		Int("sender_count", senderCount).
		Int("chunks_num", chunksNum).
		Bool("light_ingest", lightIngest).
		Msg("TestConcurrency finished")
}

// setupMockExeNode sets up a mocked execution node that responds to requests for
// chunk states. Any requests that don't correspond to an execution receipt in
// the input ers list result in the test failing.
// It also drops the first request for each chunk to evaluate retrials.
func setupMockExeNode(t *testing.T, node mock.GenericNode, verID flow.Identifier, ers []verification.CompleteExecutionResult) {
	eng := new(network.Engine)
	chunksConduit, err := node.Net.Register(engine.ChunkDataPackProvider, eng)
	assert.Nil(t, err)

	retriedChunks := make(map[flow.Identifier]struct{})

	eng.On("Process", verID, testifymock.Anything).
		Run(func(args testifymock.Arguments) {
			if req, ok := args[1].(*messages.ChunkDataPackRequest); ok {
				if _, ok := retriedChunks[req.ChunkID]; !ok {
					// this is the first request for this chunk
					// the request is dropped to evaluate retry functionality
					//retriedChunks[req.ChunkID] = struct{}{}
					//log.Debug().
					//	Hex("collection_id", logging.ID(req.ChunkID)).
					//	Msg("mock execution node drops first collection request for this collection")
					// TODO as it is switched to light node, retrial evaluation is disabled temporarily
					// return
				}

				for _, er := range ers {
					for _, chunk := range er.Receipt.ExecutionResult.Chunks {
						if chunk.ID() == req.ChunkID {
							res := &messages.ChunkDataPackResponse{
								Data:  *er.ChunkDataPacks[chunk.Index],
								Nonce: rand.Uint64(),
							}
							err := chunksConduit.Submit(res, verID)
							assert.Nil(t, err)
							return
						}
					}
				}
			}
			t.Logf("invalid chunk request (%T): %v ", args[1], args[1])
			t.Fail()
		}).
		Return(nil)

}

// setupMockCollectionNode sets up a mocked collection node that responds to requests for collections.
// Any requests that don't correspond to a collection ID in the input colls list result in the test failing.
// It also drops the first request for each collection to evaluate retrials.
func setupMockCollectionNode(t *testing.T, node mock.GenericNode, verID flow.Identifier, colls []*flow.Collection) {
	eng := new(network.Engine)
	chunksConduit, err := node.Net.Register(engine.CollectionProvider, eng)
	assert.Nil(t, err)

	retriedColl := make(map[flow.Identifier]struct{})

	eng.On("Process", verID, testifymock.Anything).
		Run(func(args testifymock.Arguments) {
			if req, ok := args[1].(*messages.CollectionRequest); ok {
				if _, ok := retriedColl[req.ID]; !ok {
					// this is the first request for this collection
					// the request is dropped to evaluate retry functionality
					retriedColl[req.ID] = struct{}{}
					//log.Debug().
					//	Hex("collection_id", logging.ID(req.ID)).
					//	Msg("mock collection node drops first collection request for this collection")
					// TODO as it is switched to light node, retrial evaluation is disabled temporarily
					// return
				}

				for _, coll := range colls {
					if coll.ID() == req.ID {
						res := &messages.CollectionResponse{
							Collection: *coll,
							Nonce:      rand.Uint64(),
						}
						err := chunksConduit.Submit(res, verID)
						assert.Nil(t, err)
						return
					}

				}
			}
			t.Logf("invalid collection request (%T): %v ", args[1], args[1])
			t.Fail()
		}).
		Return(nil)
}
