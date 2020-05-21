package testutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go/crypto/hash"
	"github.com/dapperlabs/flow-go/engine/execution/computation/virtualmachine"
	"github.com/dapperlabs/flow-go/model/flow"
)

func CreateContractDeploymentTransaction(contract string, authorizer flow.Address) flow.TransactionBody {
	encoded := hex.EncodeToString([]byte(contract))
	return flow.TransactionBody{
		Script: []byte(fmt.Sprintf(`transaction {
              prepare(signer: AuthAccount) {
                signer.setCode("%s".decodeHex())
              }
            }`, encoded)),
		Authorizers: []flow.Address{authorizer},
	}
}

func SignTransaction(tx *flow.TransactionBody, account flow.Address, privateKey flow.AccountPrivateKey, seqNum uint64) error {
	hasher, err := hash.NewHasher(privateKey.HashAlgo)
	if err != nil {
		return fmt.Errorf("cannot create hasher: %w", err)
	}

	err = tx.SetPayer(account).
		SetProposalKey(account, 0, seqNum).
		SignEnvelope(account, 0, privateKey.PrivateKey, hasher)

	if err != nil {
		return fmt.Errorf("cannot sign tx: %w", err)
	}
	return nil
}

// Generate a number of private keys
func GenerateAccountPrivateKeys(numberOfPrivateKeys int) ([]flow.AccountPrivateKey, error) {
	var privateKeys []flow.AccountPrivateKey
	for i := 0; i < numberOfPrivateKeys; i++ {
		seed := make([]byte, crypto.KeyGenSeedMinLenECDSAP256)
		_, err := rand.Read(seed)
		if err != nil {
			return nil, err
		}
		privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSAP256, seed)
		if err != nil {
			return nil, err
		}
		flowPrivateKey := flow.AccountPrivateKey{
			PrivateKey: privateKey,
			SignAlgo:   crypto.ECDSAP256,
			HashAlgo:   hash.SHA2_256,
		}
		privateKeys = append(privateKeys, flowPrivateKey)
	}
	return privateKeys, nil
}

// Create accounts on the ledger for the root account and for the private keys provided.
func BootstrappedLedger(ledger virtualmachine.Ledger, privateKeys []flow.AccountPrivateKey) (virtualmachine.Ledger, []flow.Address, error) {
	var accounts []flow.Address
	ledgerAccess := virtualmachine.LedgerDAL{Ledger: ledger}
	privateKeysIncludingRoot := []flow.AccountPrivateKey{flow.ServiceAccountPrivateKey}
	if len(privateKeys) > 0 {
		privateKeysIncludingRoot = append(privateKeysIncludingRoot, privateKeys...)
	}
	for _, account := range privateKeysIncludingRoot {
		accountPublicKey := account.PublicKey(virtualmachine.AccountKeyWeightThreshold)
		account, err := ledgerAccess.CreateAccountInLedger([]flow.AccountPublicKey{accountPublicKey})
		if err != nil {
			return nil, nil, err
		}
		accounts = append(accounts, account)
	}
	return ledger, accounts, nil
}

func SignTransactionByRoot(tx *flow.TransactionBody, seqNum uint64) error {
	return SignTransaction(tx, flow.RootAddress, flow.ServiceAccountPrivateKey, seqNum)
}

func RootBootstrappedLedger() (virtualmachine.Ledger, error) {
	ledger := make(virtualmachine.MapLedger)
	return ledger, BootstrapLedgerWithServiceAccount(ledger)
}

func BootstrapLedgerWithServiceAccount(ledger virtualmachine.Ledger) error {
	ledgerAccess := virtualmachine.LedgerDAL{Ledger: ledger}

	serviceAccountPublicKey := flow.ServiceAccountPrivateKey.PublicKey(virtualmachine.AccountKeyWeightThreshold)

	_, err := ledgerAccess.CreateAccountInLedger([]flow.AccountPublicKey{serviceAccountPublicKey})
	if err != nil {
		return err
	}

	return nil
}
