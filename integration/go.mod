module github.com/onflow/flow-go/integration

go 1.16

require (
	github.com/dapperlabs/testingdock v0.4.4
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/docker/docker v1.4.2-0.20190513124817-8c8457b0f2f8
	github.com/docker/go-connections v0.4.0
	github.com/ethereum/go-ethereum v1.10.1 // indirect
	github.com/go-openapi/strfmt v0.20.1 // indirect
	github.com/go-test/deep v1.0.7 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/onflow/cadence v0.21.3-0.20220419065337-d5202c162010
	github.com/onflow/flow-core-contracts/lib/go/contracts v0.11.2-0.20220413172500-d89ca96e6db3
	github.com/onflow/flow-core-contracts/lib/go/templates v0.11.2-0.20220413172500-d89ca96e6db3
	github.com/onflow/flow-emulator v0.30.1-0.20220421153717-0a0abc4d580b
	github.com/onflow/flow-ft/lib/go/templates v0.2.0
	github.com/onflow/flow-go v0.25.8-0.20220420171205-833e2e8849b1 // replaced by version on-disk
	github.com/onflow/flow-go-sdk v0.24.1-0.20220421152843-9ce4d554036e
	github.com/onflow/flow-go/crypto v0.24.3
	github.com/onflow/flow/protobuf/go/flow v0.2.5
	github.com/plus3it/gorecurcopy v0.0.1
	github.com/rs/zerolog v1.26.1
	github.com/stretchr/testify v1.7.1-0.20210824115523-ab6dc3262822
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	google.golang.org/grpc v1.44.0
	gopkg.in/yaml.v2 v2.4.0
)

// temp fix for MacOS build. See comment https://github.com/ory/dockertest/issues/208#issuecomment-686820414
//replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6

replace github.com/onflow/flow-go => ../
