/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fileledger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hyperledger/udo/common/flogging"
	cl "github.com/hyperledger/udo/common/ledger"
	"github.com/hyperledger/udo/common/ledger/blockledger"
	genesisconfig "github.com/hyperledger/udo/common/tools/configtxgen/localconfig"
	cb "github.com/hyperledger/udo/protos/common"
	ab "github.com/hyperledger/udo/protos/orderer"
	"github.com/hyperledger/udo/protos/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var genesisBlock = cb.NewBlock(0, nil)

func init() {
	flogging.SetModuleLevel(pkgLogID, "DEBUG")
}

type testEnv struct {
	t        *testing.T
	location string
	flf      blockledger.Factory
}

func initialize(t *testing.T) (*testEnv, *FileLedger) {
	name, err := ioutil.TempDir("", "hyperledger_udo")
	assert.NoError(t, err, "Error creating temp dir: %s", err)

	flf := New(name).(*fileLedgerFactory)
	fl, err := flf.GetOrCreate(genesisconfig.TestChainID)
	assert.NoError(t, err, "Error GetOrCreate chain")

	fl.Append(genesisBlock)
	return &testEnv{location: name, t: t, flf: flf}, fl.(*FileLedger)
}

func (tev *testEnv) tearDown() {
	tev.shutDown()
	err := os.RemoveAll(tev.location)
	if err != nil {
		tev.t.Fatalf("Error tearing down env: %s", err)
	}
}

func (tev *testEnv) shutDown() {
	tev.flf.Close()
}

type mockBlockStore struct {
	blockchainInfo             *cb.BlockchainInfo
	resultsIterator            cl.ResultsIterator
	block                      *cb.Block
	envelope                   *cb.Envelope
	txValidationCode           peer.TxValidationCode
	defaultError               error
	getBlockchainInfoError     error
	retrieveBlockByNumberError error
}

func (mbs *mockBlockStore) AddBlock(block *cb.Block) error {
	return mbs.defaultError
}

func (mbs *mockBlockStore) GetBlockchainInfo() (*cb.BlockchainInfo, error) {
	return mbs.blockchainInfo, mbs.getBlockchainInfoError
}

func (mbs *mockBlockStore) RetrieveBlocks(startNum uint64) (cl.ResultsIterator, error) {
	return mbs.resultsIterator, mbs.defaultError
}

func (mbs *mockBlockStore) RetrieveBlockByHash(blockHash []byte) (*cb.Block, error) {
	return mbs.block, mbs.defaultError
}

func (mbs *mockBlockStore) RetrieveBlockByNumber(blockNum uint64) (*cb.Block, error) {
	return mbs.block, mbs.retrieveBlockByNumberError
}

func (mbs *mockBlockStore) RetrieveTxByID(txID string) (*cb.Envelope, error) {
	return mbs.envelope, mbs.defaultError
}

func (mbs *mockBlockStore) RetrieveTxByBlockNumTranNum(blockNum uint64, tranNum uint64) (*cb.Envelope, error) {
	return mbs.envelope, mbs.defaultError
}

func (mbs *mockBlockStore) RetrieveBlockByTxID(txID string) (*cb.Block, error) {
	return mbs.block, mbs.defaultError
}

func (mbs *mockBlockStore) RetrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	return mbs.txValidationCode, mbs.defaultError
}

func (*mockBlockStore) Shutdown() {
}

type mockBlockStoreIterator struct {
	mock.Mock
}

func (m *mockBlockStoreIterator) Next() (cl.QueryResult, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

func (m *mockBlockStoreIterator) Close() {
	m.Called()
}

func TestInitialization(t *testing.T) {
	tev, fl := initialize(t)
	defer tev.tearDown()

	assert.Equal(t, uint64(1), fl.Height(), "Block height should be 1")

	block := blockledger.GetBlock(fl, 0)
	assert.NotNil(t, block, "Error retrieving genesis block")
	assert.Equal(t, genesisBlock.Header.Hash(), block.Header.Hash(), "Block hashes did no match")
}

func TestReinitialization(t *testing.T) {
	tev, ledger1 := initialize(t)
	defer tev.tearDown()

	// create a block to add to the ledger
	b1 := blockledger.CreateNextBlock(ledger1, []*cb.Envelope{{Payload: []byte("My Data")}})

	// add the block to the ledger
	ledger1.Append(b1)

	fl, err := tev.flf.GetOrCreate(genesisconfig.TestChainID)
	ledger1, ok := fl.(*FileLedger)
	assert.NoError(t, err, "Expected to sucessfully get test chain")
	assert.Equal(t, 1, len(tev.flf.ChainIDs()), "Exptected not new chain to be created")
	assert.True(t, ok, "Exptected type assertion to succeed")

	// shut down the ledger provider
	tev.shutDown()

	// re-initialize the ledger provider (not the test ledger itself!)
	provider2 := New(tev.location)

	// assert expected ledgers exist
	chains := provider2.ChainIDs()
	assert.Equal(t, 1, len(chains), "Should have recovered the chain")

	// get the existing test chain ledger
	ledger2, err := provider2.GetOrCreate(chains[0])
	assert.NoError(t, err, "Unexpected error: %s", err)

	fl = ledger2.(*FileLedger)
	assert.Equal(t, uint64(2), fl.Height(), "Block height should be 2. Got %v", fl.Height())

	block := blockledger.GetBlock(fl, 1)
	assert.NotNil(t, block, "Error retrieving block 1")
	assert.Equal(t, b1.Header.Hash(), block.Header.Hash(), "Block hashes did no match")
}

func TestAddition(t *testing.T) {
	tev, fl := initialize(t)
	defer tev.tearDown()
	info, _ := fl.blockStore.GetBlockchainInfo()
	prevHash := info.CurrentBlockHash
	fl.Append(blockledger.CreateNextBlock(fl, []*cb.Envelope{{Payload: []byte("My Data")}}))
	assert.Equal(t, uint64(2), fl.Height(), "Block height should be 2")

	block := blockledger.GetBlock(fl, 1)
	assert.NotNil(t, block, "Error retrieving genesis block")
	assert.Equal(t, prevHash, block.Header.PreviousHash, "Block hashes did no match")
}

func TestRetrieval(t *testing.T) {
	tev, fl := initialize(t)
	defer tev.tearDown()
	fl.Append(blockledger.CreateNextBlock(fl, []*cb.Envelope{{Payload: []byte("My Data")}}))
	it, num := fl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Oldest{}})
	defer it.Close()
	assert.Zero(t, num, "Expected genesis block iterator, but got %d", num)

	block, status := it.Next()
	assert.Equal(t, cb.Status_SUCCESS, status, "Expected to successfully read the genesis block")
	assert.Zero(t, block.Header.Number, "Expected to successfully retrieve the genesis block")

	block, status = it.Next()
	assert.Equal(t, cb.Status_SUCCESS, status, "Expected to successfully read the second block")
	assert.Equal(
		t,
		uint64(1),
		block.Header.Number,
		"Expected to successfully retrieve the second block but got block number %d", block.Header.Number)
}

func TestBlockedRetrieval(t *testing.T) {
	tev, fl := initialize(t)
	defer tev.tearDown()
	it, num := fl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 1}}})
	defer it.Close()
	if num != 1 {
		t.Fatalf("Expected block iterator at 1, but got %d", num)
	}
	assert.Equal(t, uint64(1), num, "Expected block iterator at 1, but got %d", num)

	fl.Append(blockledger.CreateNextBlock(fl, []*cb.Envelope{{Payload: []byte("My Data")}}))

	block, status := it.Next()
	assert.Equal(t, cb.Status_SUCCESS, status, "Expected to successfully read the second block")
	assert.Equal(
		t,
		uint64(1),
		block.Header.Number,
		"Expected to successfully retrieve the second block but got block number %d", block.Header.Number)

	fl.Append(blockledger.CreateNextBlock(fl, []*cb.Envelope{{Payload: []byte("My Data")}}))

	block, status = it.Next()
	assert.Equal(t, cb.Status_SUCCESS, status, "Expected to successfully read the third block")
	assert.Equal(t, uint64(2), block.Header.Number, "Expected to successfully retrieve the third block")
}

func TestBlockstoreError(t *testing.T) {
	// Since this test only ensures failed GetBlockchainInfo
	// is properly handled. We don't bother creating fully
	// legit ledgers here (without genesis block).
	{
		fl := &FileLedger{
			blockStore: &mockBlockStore{
				blockchainInfo:         nil,
				getBlockchainInfoError: fmt.Errorf("Error getting blockchain info"),
			},
			signal: make(chan struct{}),
		}
		assert.Panics(
			t,
			func() {
				fl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Newest{}})
			},
			"Expected Iterator() to panic if blockstore operation fails")

		assert.Panics(
			t,
			func() { fl.Height() },
			"Expected Height() to panic if blockstore operation fails ")
	}

	{
		fl := &FileLedger{
			blockStore: &mockBlockStore{
				blockchainInfo:             &cb.BlockchainInfo{Height: uint64(1)},
				getBlockchainInfoError:     nil,
				retrieveBlockByNumberError: fmt.Errorf("Error retrieving block by number"),
			},
			signal: make(chan struct{}),
		}
		it, _ := fl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 42}}})
		defer it.Close()
		assert.IsType(
			t,
			&blockledger.NotFoundErrorIterator{},
			it,
			"Expected Not Found Error if seek number is greater than ledger height")
	}

	{
		resultsIterator := &mockBlockStoreIterator{}
		resultsIterator.On("Next").Return(nil, errors.New("a mocked error"))
		resultsIterator.On("Close").Return()
		fl := &FileLedger{
			blockStore: &mockBlockStore{
				blockchainInfo:             &cb.BlockchainInfo{Height: uint64(1)},
				getBlockchainInfoError:     nil,
				retrieveBlockByNumberError: fmt.Errorf("Error retrieving block by number"),
				resultsIterator:            resultsIterator,
			},
			signal: make(chan struct{}),
		}
		it, _ := fl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Oldest{}})
		defer it.Close()
		_, status := it.Next()
		assert.Equal(t, cb.Status_SERVICE_UNAVAILABLE, status, "Expected service unavailable error")
	}
}
