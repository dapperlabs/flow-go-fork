package testutils

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/onflow/atree"
	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"

	"github.com/onflow/flow-go/fvm/environment"
	"github.com/onflow/flow-go/fvm/evm/types"
	"github.com/onflow/flow-go/fvm/meter"
	"github.com/onflow/flow-go/model/flow"
)

var TestFlowEVMRootAddress = flow.BytesToAddress([]byte("FlowEVM"))
var TestComputationLimit = uint(math.MaxUint64 - 1)

func RunWithTestFlowEVMRootAddress(t testing.TB, backend types.Backend, f func(flow.Address)) {
	as := environment.NewAccountStatus()
	backend.SetValue(TestFlowEVMRootAddress[:], []byte(flow.AccountStatusKey), as.ToBytes())
	f(TestFlowEVMRootAddress)
}

func RunWithTestBackend(t testing.TB, f func(*TestBackend)) {
	tb := &TestBackend{
		testValueStore:   getSimpleValueStore(),
		testEventEmitter: getSimpleEventEmitter(),
		testMeter:        getSimpleMeter(),
	}
	f(tb)
}

func ConvertToCadence(data []byte) []cadence.Value {
	ret := make([]cadence.Value, len(data))
	for i, v := range data {
		ret[i] = cadence.UInt8(v)
	}
	return ret
}

func fullKey(owner, key []byte) string {
	return string(owner) + "~" + string(key)
}

func getSimpleValueStore() *testValueStore {
	data := make(map[string][]byte)
	allocator := make(map[string]uint64)
	return &testValueStore{
		getValue: func(owner, key []byte) ([]byte, error) {
			fk := fullKey(owner, key)
			return data[fk], nil
		},
		setValue: func(owner, key, value []byte) error {
			fk := fullKey(owner, key)
			data[fk] = value
			return nil
		},
		valueExists: func(owner, key []byte) (bool, error) {
			fk := fullKey(owner, key)
			data := data[fk]
			return len(data) > 0, nil

		},
		allocateStorageIndex: func(owner []byte) (atree.StorageIndex, error) {
			index := allocator[string(owner)]
			var data [8]byte
			allocator[string(owner)] = index + 1
			binary.BigEndian.PutUint64(data[:], index)
			return atree.StorageIndex(data), nil
		},
		totalStorageSize: func() int {
			size := 0
			for k, v := range data {
				size += len(k) + len(v)
			}
			for k := range allocator {
				size += len(k) + 8
			}
			return size
		},
	}
}

func getSimpleEventEmitter() *testEventEmitter {
	events := make(flow.EventsList, 0)
	return &testEventEmitter{
		emitFlowEvent: func(etype flow.EventType, payload []byte) error {
			events = append(events, flow.Event{Type: etype, Payload: payload})
			return nil
		},
		events: func() flow.EventsList {
			return events
		},
		reset: func() {
			events = make(flow.EventsList, 0)
		},
	}
}

func getSimpleMeter() *testMeter {
	computationLimit := TestComputationLimit
	compUsed := uint(0)
	return &testMeter{
		meterComputation: func(kind common.ComputationKind, intensity uint) error {
			compUsed += intensity
			if compUsed > computationLimit {
				return fmt.Errorf("computation limit has hit %d", computationLimit)
			}
			return nil
		},
		computationAvailable: func(kind common.ComputationKind, intensity uint) bool {
			return compUsed+intensity < computationLimit
		},
		computationUsed: func() (uint64, error) {
			return uint64(compUsed), nil
		},
	}
}

type TestBackend struct {
	*testValueStore
	*testMeter
	*testEventEmitter
}

func (tb *TestBackend) TotalStorageSize() int {
	if tb.totalStorageSize == nil {
		panic("method not set")
	}
	return tb.totalStorageSize()
}

func (tb *TestBackend) DropEvents() {
	if tb.reset == nil {
		panic("method not set")
	}
	tb.reset()
}

type testValueStore struct {
	getValue             func(owner, key []byte) ([]byte, error)
	setValue             func(owner, key, value []byte) error
	valueExists          func(owner, key []byte) (bool, error)
	allocateStorageIndex func(owner []byte) (atree.StorageIndex, error)
	totalStorageSize     func() int
}

var _ environment.ValueStore = &testValueStore{}

func (vs *testValueStore) GetValue(owner, key []byte) ([]byte, error) {
	if vs.getValue == nil {
		panic("method not set")
	}
	return vs.getValue(owner, key)
}

func (vs *testValueStore) SetValue(owner, key, value []byte) error {
	if vs.setValue == nil {
		panic("method not set")
	}
	return vs.setValue(owner, key, value)
}

func (vs *testValueStore) ValueExists(owner, key []byte) (bool, error) {
	if vs.valueExists == nil {
		panic("method not set")
	}
	return vs.valueExists(owner, key)
}

func (vs *testValueStore) AllocateStorageIndex(owner []byte) (atree.StorageIndex, error) {
	if vs.allocateStorageIndex == nil {
		panic("method not set")
	}
	return vs.allocateStorageIndex(owner)
}

func (vs *testValueStore) TotalStorageSize() int {
	if vs.totalStorageSize == nil {
		panic("method not set")
	}
	return vs.totalStorageSize()
}

type testMeter struct {
	meterComputation       func(common.ComputationKind, uint) error
	computationAvailable   func(common.ComputationKind, uint) bool
	computationUsed        func() (uint64, error)
	computationIntensities func() meter.MeteredComputationIntensities

	meterMemory func(usage common.MemoryUsage) error
	memoryUsed  func() (uint64, error)

	meterEmittedEvent      func(byteSize uint64) error
	totalEmittedEventBytes func() uint64

	interactionUsed func() (uint64, error)
}

var _ environment.Meter = &testMeter{}

func (m *testMeter) MeterComputation(
	kind common.ComputationKind,
	intensity uint,
) error {
	if m.meterComputation == nil {
		panic("method not set")
	}
	return m.meterComputation(kind, intensity)
}

func (m *testMeter) ComputationAvailable(
	kind common.ComputationKind,
	intensity uint,
) bool {
	if m.computationAvailable == nil {
		panic("method not set")
	}
	return m.computationAvailable(kind, intensity)
}

func (m *testMeter) ComputationIntensities() meter.MeteredComputationIntensities {
	if m.computationIntensities == nil {
		panic("method not set")
	}
	return m.computationIntensities()
}

func (m *testMeter) ComputationUsed() (uint64, error) {
	if m.computationUsed == nil {
		panic("method not set")
	}
	return m.computationUsed()
}

func (m *testMeter) MeterMemory(usage common.MemoryUsage) error {
	if m.meterMemory == nil {
		panic("method not set")
	}
	return m.meterMemory(usage)
}

func (m *testMeter) MemoryUsed() (uint64, error) {
	if m.memoryUsed == nil {
		panic("method not set")
	}
	return m.memoryUsed()
}

func (m *testMeter) InteractionUsed() (uint64, error) {
	if m.interactionUsed == nil {
		panic("method not set")
	}
	return m.interactionUsed()
}

func (m *testMeter) MeterEmittedEvent(byteSize uint64) error {
	if m.meterEmittedEvent == nil {
		panic("method not set")
	}
	return m.meterEmittedEvent(byteSize)
}

func (m *testMeter) TotalEmittedEventBytes() uint64 {
	if m.totalEmittedEventBytes == nil {
		panic("method not set")
	}
	return m.totalEmittedEventBytes()
}

type testEventEmitter struct {
	emitEvent              func(event cadence.Event) error
	emitFlowEvent          func(etype flow.EventType, payload []byte) error
	events                 func() flow.EventsList
	serviceEvents          func() flow.EventsList
	convertedServiceEvents func() flow.ServiceEventList
	reset                  func()
}

var _ environment.EventEmitter = &testEventEmitter{}

func (vs *testEventEmitter) EmitEvent(event cadence.Event) error {
	if vs.emitEvent == nil {
		panic("method not set")
	}
	return vs.emitEvent(event)
}

func (vs *testEventEmitter) EmitFlowEvent(etype flow.EventType, payload []byte) error {
	if vs.emitFlowEvent == nil {
		panic("method not set")
	}
	return vs.emitFlowEvent(etype, payload)
}

func (vs *testEventEmitter) Events() flow.EventsList {
	if vs.events == nil {
		panic("method not set")
	}
	return vs.events()
}

func (vs *testEventEmitter) ServiceEvents() flow.EventsList {
	if vs.serviceEvents == nil {
		panic("method not set")
	}
	return vs.serviceEvents()
}

func (vs *testEventEmitter) ConvertedServiceEvents() flow.ServiceEventList {
	if vs.convertedServiceEvents == nil {
		panic("method not set")
	}
	return vs.convertedServiceEvents()
}

func (vs *testEventEmitter) Reset() {
	if vs.reset == nil {
		panic("method not set")
	}
	vs.reset()
}
